package notmuch

// Copyright © 2015 The go.bindings Authors. Authors can be found in the AUTHORS file.
// Licensed under the GPLv3 or later.
// See COPYING at the root of the repository for details.

// #cgo LDFLAGS: -lnotmuch
// #include <stdlib.h>
// #include <notmuch.h>
import "C"

type (
	// Thread represents a notmuch thread.
	Thread struct {
		id      string
		cptr    *C.notmuch_thread_t
		threads *Threads
	}

	// Messages represents notmuch messages.
	Messages struct {
		cptr   *C.notmuch_messages_t
		thread *Thread
	}
)

// GetSubject returns the subject of a thread.
func (t *Thread) GetSubject() string {
	cstr := C.notmuch_thread_get_subject(t.toC())
	str := C.GoString(cstr)
	return str
}

// GetID returns the ID of the thread.
func (t *Thread) GetID() string {
	return t.id
}

func (t *Thread) toC() *C.notmuch_thread_t {
	return (*C.notmuch_thread_t)(t.cptr)
}
