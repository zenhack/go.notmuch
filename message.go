package notmuch

// Copyright © 2015 The go.notmuch Authors. Authors can be found in the AUTHORS file.
// Licensed under the GPLv3 or later.
// See COPYING at the root of the repository for details.

// #cgo LDFLAGS: -lnotmuch
// #include <stdlib.h>
// #include <notmuch.h>
import "C"

// Message represents a notmuch message.
type Message struct {
	cptr     *C.notmuch_message_t
	messages *Messages
	thread   *Thread
}

// GetID returns the message ID.
func (m *Message) GetID() string {
	return C.GoString(C.notmuch_message_get_message_id(m.cptr))
}
