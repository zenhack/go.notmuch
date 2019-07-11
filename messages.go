package notmuch

// Copyright © 2015 The go.notmuch Authors. Authors can be found in the AUTHORS file.
// Licensed under the GPLv3 or later.
// See COPYING at the root of the repository for details.

// #cgo LDFLAGS: -lnotmuch
// #include <stdlib.h>
// #include <notmuch.h>
import "C"
import (
	"unsafe"
)

// Messages represents notmuch messages.
type Messages cStruct

func (ms *Messages) Close() error {
	return (*cStruct)(ms).doClose(func() error {
		C.notmuch_messages_destroy(ms.toC())
		return nil
	})
}

func (ms *Messages) toC() *C.notmuch_messages_t {
	return (*C.notmuch_messages_t)(ms.cptr)
}

// Next retrieves the next message from the result set. Next returns true if a message
// was successfully retrieved.
func (ms *Messages) Next(m **Message) bool {
	if !ms.valid() {
		return false
	}
	*m = ms.get()
	C.notmuch_messages_move_to_next(ms.toC())
	return true
}

// Tags return a list of tags from all messages.
//
// WARNING: You can no longer iterate over messages after calling this
// function, because the iterator will point at the end of the list.  We do not
// have a function to reset the iterator yet and the only way how you can
// iterate over the list again is to recreate the message list.
func (ms *Messages) Tags() *Tags {
	ctags := C.notmuch_messages_collect_tags(ms.toC())
	// TODO(kalbasit): notmuch_messages_collect_tags can return NULL on error
	// but there's not explanation on what kind of error can occur. We should handle
	// it as OOM for now but we eventually have to narrow it down.
	checkOOM(unsafe.Pointer(ctags))
	tags := &Tags{
		cptr:   unsafe.Pointer(ctags),
		parent: (*cStruct)(ms),
	}
	setGcClose(tags)
	return tags
}

func (ms *Messages) get() *Message {
	cmessage := C.notmuch_messages_get(ms.toC())
	checkOOM(unsafe.Pointer(cmessage))
	message := &Message{
		cptr:   unsafe.Pointer(cmessage),
		parent: (*cStruct)(ms),
	}
	setGcClose(message)
	return message
}

func (ms *Messages) valid() bool {
	cbool := C.notmuch_messages_valid(ms.toC())
	return int(cbool) != 0
}
