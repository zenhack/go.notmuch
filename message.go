package notmuch

// Copyright © 2015 The go.notmuch Authors. Authors can be found in the AUTHORS file.
// Licensed under the GPLv3 or later.
// See COPYING at the root of the repository for details.

// #cgo LDFLAGS: -lnotmuch
// #include <stdlib.h>
// #include <notmuch.h>
import "C"
import (
	"time"
	"unsafe"
)

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
		return nil, ErrNoRepliesOrPointerNotFromThread
	}
	return &Messages{
		cptr:   cmsgs,
		thread: m.thread,
	}, nil
}

// Filename returns the absolute path of the email message.
//
// Note: If this message corresponds to multiple files in the mail store, (that
// is, multiple files contain identical message IDs), this function will
// arbitrarily return a single one of those filenames. See Filenames for
// returning the complete list of filenames.
func (m *Message) Filename() string {
	return C.GoString(C.notmuch_message_get_filename(m.cptr))
}

// Filenames returns *Filenames an iterator to get the message's filenames.
// Each filename in the iterator is an absolute filename.
func (m *Message) Filenames() *Filenames {
	return &Filenames{
		cptr:    C.notmuch_message_get_filenames(m.cptr),
		message: m,
	}
}

// Date returns the date of the message.
func (m *Message) Date() time.Time {
	ctime := C.notmuch_message_get_date(m.cptr)
	return time.Unix(int64(ctime), 0)
}
