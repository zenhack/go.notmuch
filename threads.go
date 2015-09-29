package notmuch

// Copyright © 2015 The go.notmuch Authors. Authors can be found in the AUTHORS file.
// Licensed under the GPLv3 or later.
// See COPYING at the root of the repository for details.

// #cgo LDFLAGS: -lnotmuch
// #include <stdlib.h>
// #include <notmuch.h>
import "C"

import (
	"runtime"
	"unsafe"
)

// Threads represents notmuch threads.
type Threads struct {
	cptr  *C.notmuch_threads_t
	query *Query
}

// Next retrieves the next thread from the result set. Next returns true if a thread
// was successfully retrieved.
func (ts *Threads) Next(t *Thread) bool {
	if !ts.valid() {
		return false
	}
	*t = *ts.get()
	C.notmuch_threads_move_to_next(ts.cptr)
	return true
}

func (ts *Threads) get() *Thread {
	cthread := C.notmuch_threads_get(ts.cptr)
	checkOOM(unsafe.Pointer(cthread))
	thread := &Thread{
		cptr:    cthread,
		threads: ts,
	}
	runtime.SetFinalizer(thread, func(t *Thread) {
		C.notmuch_thread_destroy(t.cptr)
	})
	return thread
}

func (ts *Threads) valid() bool {
	cbool := C.notmuch_threads_valid(ts.cptr)
	return int(cbool) != 0
}
