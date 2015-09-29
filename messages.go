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

// Messages represents notmuch messages.
type Messages struct {
	cptr   *C.notmuch_messages_t
	thread *Thread
}

// Next retrieves the next message from the result set. Next returns true if a message
// was successfully retrieved.
func (ms *Messages) Next(m *Message) bool {
	if !ms.valid() {
		return false
	}
	*m = *ms.get()
	C.notmuch_messages_move_to_next(ms.cptr)
	return true
}

// Tags return a list of tags from all messages.
//
// WARNING: You can no longer iterate over messages after calling this
// function, because the iterator will point at the end of the list.  We do not
// have a function to reset the iterator yet and the only way how you can
// iterate over the list again is to recreate the message list.
func (ms *Messages) Tags() *Tags {
	ts := &Tags{
		cptr:   C.notmuch_messages_collect_tags(ms.cptr),
		thread: ms.thread,
	}
	// TODO(kalbasit): notmuch_messages_collect_tags can return NULL on error
	// but there's not explanation on what kind of error can occur. We should handle
	// it as OOM for now but we eventually have to narrow it down.
	checkOOM(unsafe.Pointer(ts.cptr))
	runtime.SetFinalizer(ts, func(ts *Tags) {
		C.notmuch_tags_destroy(ts.cptr)
	})
	return ts
}

func (ms *Messages) get() *Message {
	cmessage := C.notmuch_messages_get(ms.cptr)
	checkOOM(unsafe.Pointer(cmessage))
	message := &Message{
		cptr:     cmessage,
		messages: ms,
	}
	runtime.SetFinalizer(message, func(m *Message) {
		C.notmuch_message_destroy(m.cptr)
	})
	return message
}

func (ms *Messages) valid() bool {
	cbool := C.notmuch_messages_valid(ms.cptr)
	return int(cbool) != 0
}
