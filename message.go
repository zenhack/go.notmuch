package notmuch

// Copyright © 2015 The go.notmuch Authors. Authors can be found in the AUTHORS file.
// Licensed under the GPLv3 or later.
// See COPYING at the root of the repository for details.

// #cgo LDFLAGS: -lnotmuch
// #include <stdlib.h>
// #include <notmuch.h>
import "C"
import "unsafe"

// Message represents a notmuch message.
type Message struct {
	cptr     *C.notmuch_message_t
	messages *Messages
	thread   *Thread
}

// ID returns the message ID.
func (m *Message) ID() string {
	return C.GoString(C.notmuch_message_get_message_id(m.cptr))
}

// ThreadID returns the ID of the thread to which this message belongs to.
func (m *Message) ThreadID() string {
	return C.GoString(C.notmuch_message_get_thread_id(m.cptr))
}

// Replies returns the replies of a message.
func (m *Message) Replies() (*Messages, error) {
	cmsgs := C.notmuch_message_get_replies(m.cptr)
	if unsafe.Pointer(cmsgs) == nil {
		return nil, ErrNoReplies
	}
	return &Messages{
		cptr:   cmsgs,
		thread: m.thread,
	}, nil
}
