package bypass

import "unsafe"

// The below need to match runtime/chan.go for this to work.

// hchan is the internal runtime representation of a go chan.
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

// mutex is the internal runtime mutex.
type mutex struct {
	// Futex-based impl treats it as uint32 key,
	// while sema-based impl as M* waitm.
	// Used to be a union, but unions break precise GC.
	key uintptr
}

// Internal runtime functions

//go:linkname lock runtime.lock
func lock(l *mutex)

//go:linkname unlock runtime.unlock
func unlock(l *mutex)
