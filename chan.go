package bypass

import (
	"fmt"
	"reflect"
	"unsafe"

	_ "unsafe"
)

type Chan struct {
	hchan *hchan
	v     reflect.Value
}

// WrapChan wraps a standard channel in a bypass.Chan so that it can be
// modified in strange ways.
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

// Elems returns a slice with all of the elements in the channel.
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

// Lock locks the channel so nothing can write to it until it is unlocked.
func (c *Chan) Lock() {
	lock(&c.hchan.lock)
}

// Unlock unlocks the channel.
func (c *Chan) Unlock() {
	unlock(&c.hchan.lock)
}
