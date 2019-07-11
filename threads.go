package notmuch

// Copyright © 2015 The go.notmuch Authors. Authors can be found in the AUTHORS file.
// Licensed under the GPLv3 or later.
// See COPYING at the root of the repository for details.

// #cgo LDFLAGS: -lnotmuch
// #include <stdlib.h>
// #include <notmuch.h>
import "C"

import "unsafe"

// Threads represents notmuch threads.
type Threads cStruct

func (ts *Threads) toC() *C.notmuch_threads_t {
	return (*C.notmuch_threads_t)(ts.cptr)
}

func (ts *Threads) Close() error {
	return (*cStruct)(ts).doClose(func() error {
		C.notmuch_threads_destroy(ts.toC())
		return nil
	})
}

// Next retrieves the next thread from the result set. Next returns true if a thread
// was successfully retrieved.
func (ts *Threads) Next(t **Thread) bool {
	if !ts.valid() {
		return false
	}
	*t = ts.get()
	C.notmuch_threads_move_to_next(ts.toC())
	return true
}

func (ts *Threads) get() *Thread {
	cthread := C.notmuch_threads_get(ts.toC())
	checkOOM(unsafe.Pointer(cthread))
	thread := &Thread{
		cptr:   unsafe.Pointer(cthread),
		parent: (*cStruct)(ts),
	}
	setGcClose(thread)
	return thread
}

func (ts *Threads) valid() bool {
	cbool := C.notmuch_threads_valid(ts.toC())
	return int(cbool) != 0
}
