package bypass

import (
	"fmt"
	"reflect"
	"unsafe"

	_ "unsafe"
)

const locked uintptr = 1

type hchan struct {
	qcount   uint           // total data in the queue
	dataqsiz uint           // size of the circular queue
	buf      unsafe.Pointer // points to an array of dataqsiz elements
	elemsize uint16
	closed   uint32
	elemtype uintptr // *_type // element type
	sendx    uint    // send index
	recvx    uint    // receive index
	recvq    waitq   // list of recv waiters
	sendq    waitq   // list of send waiters
	lock     mutex
}

type waitq struct {
	first uintptr // *sudog
	last  uintptr // *sudog
}

type mutex struct {
	// Futex-based impl treats it as uint32 key,
	// while sema-based impl as M* waitm.
	// Used to be a union, but unions break precise GC.
	key uintptr
}

//go:linkname lock runtime.lock
func lock(l *mutex)

//go:linkname unlock runtime.unlock
func unlock(l *mutex)

type Chan struct {
	hchan *hchan
	v     reflect.Value
}

func WrapChan(c interface{}) *Chan {
	v := reflect.ValueOf(c)
	vt := v.Type()
	if vt.Kind() != reflect.Chan {
		panic("can only wrap chans")
	}
	size := uint16(vt.Elem().Size())
	upv := unsafe.Pointer(uintptr(unsafe.Pointer(&v)) + offsetPtr)
	ch := (*hchan)(*(*unsafe.Pointer)(upv))
	if ch.elemsize != size {
		panic(fmt.Sprintf("reflect elem size does not match chan size %d != %d", size, ch.elemsize))
	}
	return &Chan{
		hchan: ch,
		v:     v,
	}
}

func (c *Chan) Elems() interface{} {
	count := c.v.Len()
	elemType := c.v.Type().Elem()
	arr := reflect.MakeSlice(reflect.SliceOf(elemType), 0, count)
	c.Lock()
	defer c.Unlock()
	elemSize := uintptr(c.hchan.elemsize)

	for i := 0; i < count; i++ {
		idx := uint(i) + c.hchan.recvx
		if idx >= c.hchan.dataqsiz {
			idx -= c.hchan.dataqsiz
		}
		ptr := unsafe.Pointer(uintptr(c.hchan.buf) + elemSize*uintptr(idx))
		arr = reflect.Append(arr, reflect.NewAt(elemType, ptr).Elem())
	}
	return arr.Interface()
}

func (c *Chan) Lock() {
	lock(&c.hchan.lock)
}

func (c *Chan) Unlock() {
	unlock(&c.hchan.lock)
}
