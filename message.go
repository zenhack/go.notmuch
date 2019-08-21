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
		cptr:   unsafe.Pointer(cmsgs),
		parent: (*cStruct)(m),
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
		cptr:   unsafe.Pointer(ctags),
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

// MaildirFlagsToTags adds/removes tags according to maildir flags in the message filename(s).
// This function examines the filenames of 'message' for maildir
// flags, and adds or removes tags on 'message' as follows when these
// flags are present:
//
//      Flag    Action if present
//      ----    -----------------
//      'D'     Adds the "draft" tag to the message
//      'F'     Adds the "flagged" tag to the message
//      'P'     Adds the "passed" tag to the message
//      'R'     Adds the "replied" tag to the message
//      'S'     Removes the "unread" tag from the message
//
// For each flag that is not present, the opposite action (add/remove)
// is performed for the corresponding tags.
//
// Flags are identified as trailing components of the filename after a
// sequence of ":2,".
//
// If there are multiple filenames associated with this message, the
// flag is considered present if it appears in one or more
// filenames. (That is, the flags from the multiple filenames are
// combined with the logical OR operator.)
//
// A client can ensure that notmuch database tags remain synchronized
// with maildir flags by calling this function after each call to
// DB.AddMessage. See also Message.TagsToMaildirFlags for synchronizing
// tag changes back to maildir flags.
func (m *Message) MaildirFlagsToTags() error {
	return statusErr(C.notmuch_message_maildir_flags_to_tags(m.toC()))
}

// TagsToMaildirFlags renames message filename(s) to encode tags as maildir flags.
// Specifically, for each filename corresponding to this message:
//
// If the filename is not in a maildir directory, do nothing. (A
// maildir directory is determined as a directory named "new" or
// "cur".) Similarly, if the filename has invalid maildir info,
// (repeated or outof-ASCII-order flag characters after ":2,"), then
// do nothing.
//
// If the filename is in a maildir directory, rename the file so that
// its filename ends with the sequence ":2," followed by zero or more
// of the following single-character flags (in ASCII order):
//
//   * flag 'D' if the message has the "draft" tag
//   * flag 'F' if the message has the "flagged" tag
//   * flag 'P' if the message has the "passed" tag
//   * flag 'R' if the message has the "replied" tag
//   * flag 'S' if the message does not have the "unread" tag
//
// Any existing flags unmentioned in the list above will be preserved
// in the renaming.
//
// Also, if this filename is in a directory named "new", rename it to
// be within the neighboring directory named "cur".
//
// A client can ensure that maildir filename flags remain synchronized
// with notmuch database tags by calling this function after changing
// tags. See also Message.MaildirFlagsToTags for synchronizing maildir flag
// changes back to tags.
func (m *Message) TagsToMaildirFlags() error {
	return statusErr(C.notmuch_message_tags_to_maildir_flags(m.toC()))
}
