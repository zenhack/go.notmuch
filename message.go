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
type Message cStruct

func (m *Message) toC() *C.notmuch_message_t {
	return (*C.notmuch_message_t)(m.cptr)
}

func (m *Message) Close() error {
	return (*cStruct)(m).doClose(func() error {
		C.notmuch_message_destroy(m.toC())
		return nil
	})
}

// ID returns the message ID.
func (m *Message) ID() string {
	return C.GoString(C.notmuch_message_get_message_id(m.toC()))
}

// ThreadID returns the ID of the thread to which this message belongs to.
func (m *Message) ThreadID() string {
	return C.GoString(C.notmuch_message_get_thread_id(m.toC()))
}

// Replies returns the replies of a message.
func (m *Message) Replies() (*Messages, error) {
	cmsgs := C.notmuch_message_get_replies(m.toC())
	if unsafe.Pointer(cmsgs) == nil {
		return nil, ErrNoRepliesOrPointerNotFromThread
	}
	// We point the messages object directly at our thread, rather than having
	// the gc reference go through this message:
	msgs := &Messages{
		cptr: unsafe.Pointer(cmsgs),
		parent: m.parent,
	}
	setGcClose(msgs)
	return msgs, nil
}

// Filename returns the absolute path of the email message.
//
// Note: If this message corresponds to multiple files in the mail store, (that
// is, multiple files contain identical message IDs), this function will
// arbitrarily return a single one of those filenames. See Filenames for
// returning the complete list of filenames.
func (m *Message) Filename() string {
	return C.GoString(C.notmuch_message_get_filename(m.toC()))
}

// Filenames returns *Filenames an iterator to get the message's filenames.
// Each filename in the iterator is an absolute filename.
func (m *Message) Filenames() *Filenames {
	return &Filenames{
		cptr:    C.notmuch_message_get_filenames(m.toC()),
		message: m,
	}
}

// Date returns the date of the message.
func (m *Message) Date() time.Time {
	ctime := C.notmuch_message_get_date(m.toC())
	return time.Unix(int64(ctime), 0)
}

// Header returns the value of the header.
func (m *Message) Header(name string) string {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	return C.GoString(C.notmuch_message_get_header(m.toC(), cname))
}

// Tags returns the tags for the current message, returning a *Tags which can
// be used to iterate over all tags using `Tags.Next(Tag)`
func (m *Message) Tags() *Tags {
	ctags := C.notmuch_message_get_tags(m.toC())
	tags := &Tags{
		cptr: unsafe.Pointer(ctags),
		parent: (*cStruct)(m),
	}
	setGcClose(tags)
	return tags
}

// AddTag adds a tag to the message.
func (m *Message) AddTag(tag string) error {
	ctag := C.CString(tag)
	defer C.free(unsafe.Pointer(ctag))
	return statusErr(C.notmuch_message_add_tag(m.toC(), ctag))
}

// RemoveTag removes a tag from the message.
func (m *Message) RemoveTag(tag string) error {
	ctag := C.CString(tag)
	defer C.free(unsafe.Pointer(ctag))
	return statusErr(C.notmuch_message_remove_tag(m.toC(), ctag))
}

// RemoveAllTags removes all tags from the message.
func (m *Message) RemoveAllTags() error {
	return statusErr(C.notmuch_message_remove_all_tags(m.toC()))
}

// Atomic allows a transactional change of tags to the message.
func (m *Message) Atomic(callback func(*Message)) error {
	if err := statusErr(C.notmuch_message_freeze(m.toC())); err != nil {
		return err
	}
	callback(m)
	return statusErr(C.notmuch_message_thaw(m.toC()))
}
