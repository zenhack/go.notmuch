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

// Get fetches the currently selected message.
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
