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
		return nil, ErrNoReplies
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

// Header returns the value of the header.
func (m *Message) Header(name string) string {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	return C.GoString(C.notmuch_message_get_header(m.cptr, cname))
}

// Tags returns the tags for the current message, returning a *Tags which can
// be used to iterate over all tags using `Tags.Next(Tag)`
func (m *Message) Tags() *Tags {
	ts := &Tags{
		cptr:    C.notmuch_message_get_tags(m.cptr),
		message: m,
		thread:  m.thread,
	}
	return ts
}

// AddTag adds a tag to the message.
func (m *Message) AddTag(tag string) error {
	ctag := C.CString(tag)
	defer C.free(unsafe.Pointer(ctag))
	return statusErr(C.notmuch_message_add_tag(m.cptr, ctag))
}

// RemoveTag removes a tag from the message.
func (m *Message) RemoveTag(tag string) error {
	ctag := C.CString(tag)
	defer C.free(unsafe.Pointer(ctag))
	return statusErr(C.notmuch_message_remove_tag(m.cptr, ctag))
}

// RemoveAllTags removes all tags from the message.
func (m *Message) RemoveAllTags() error {
	return statusErr(C.notmuch_message_remove_all_tags(m.cptr))
}

// Atomic allows a transactional change of tags to the message.
func (m *Message) Atomic(callback func(*Message)) error {
	if err := statusErr(C.notmuch_message_freeze(m.cptr)); err != nil {
		return err
	}
	callback(m)
	return statusErr(C.notmuch_message_thaw(m.cptr))
}


// Destroy a Message object.
//
// It can be useful to call this method in the case of a single
// Query object with many messages in the result, (such as iterating
// over the entire database). Otherwise, it's fine to never call this
// function and there will still be no memory leaks. (The memory from
// the messages get reclaimed when the containing query is destroyed.)
func (m *Message) Close() error {
	C.notmuch_message_destroy(m.cptr)
	return nil
}
