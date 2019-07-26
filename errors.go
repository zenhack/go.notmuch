package notmuch

// Copyright © 2015 The go.notmuch Authors. Authors can be found in the AUTHORS file.
// Licensed under the GPLv3 or later.
// See COPYING at the root of the repository for details.

// #cgo LDFLAGS: -lnotmuch
// #include <stdlib.h>
// #include <notmuch.h>
import "C"
import "errors"

import "unsafe"

type status C.notmuch_status_t

var (
	// ErrOutOfMemory is returned when an Out of memory occured.
	ErrOutOfMemory = statusErr(C.NOTMUCH_STATUS_OUT_OF_MEMORY)

	// ErrReadOnlyDB is returned when an attempt was made to write to a
	// database opened in read-only mode.
	ErrReadOnlyDB = statusErr(C.NOTMUCH_STATUS_READ_ONLY_DATABASE)

	// ErrXapianException is returned when a xapian exception occured.
	ErrXapianException = statusErr(C.NOTMUCH_STATUS_XAPIAN_EXCEPTION)

	// ErrFileError is returned when an error occurred trying to read or write to
	// a file (this could be file not found, permission denied, etc.)
	ErrFileError = statusErr(C.NOTMUCH_STATUS_FILE_ERROR)

	// ErrFileNotEmail is returned when a file was presented that doesn't appear
	// to be an email message.
	ErrFileNotEmail = statusErr(C.NOTMUCH_STATUS_FILE_NOT_EMAIL)

	// ErrDuplicateMessageID is returned when a file contains a message ID that
	// is identical to a message already in the database.
	ErrDuplicateMessageID = statusErr(C.NOTMUCH_STATUS_DUPLICATE_MESSAGE_ID)

	// ErrNullPointer is returned when the user erroneously passed a NULL pointer
	// to a notmuch function.
	ErrNullPointer = statusErr(C.NOTMUCH_STATUS_NULL_POINTER)

	// ErrTagTooLong is returned when a tag value is too long (exceeds TagMax)
	ErrTagTooLong = statusErr(C.NOTMUCH_STATUS_TAG_TOO_LONG)

	// ErrUnbalancedFreezeThaw is returned when Message.Thaw() was called more
	// times than Message.Freeze().
	ErrUnbalancedFreezeThaw = statusErr(C.NOTMUCH_STATUS_UNBALANCED_FREEZE_THAW)

	// ErrUnbalancedAtomic DB.EndAtomic() has been called more times than DB.BeginAtomic()
	ErrUnbalancedAtomic = statusErr(C.NOTMUCH_STATUS_UNBALANCED_ATOMIC)

	// ErrUnsupportedOperation is returned when the operation is not supported.
	ErrUnsupportedOperation = statusErr(C.NOTMUCH_STATUS_UNSUPPORTED_OPERATION)

	// ErrUpgradeRequired is returned when the database requires an upgrade.
	ErrUpgradeRequired = statusErr(C.NOTMUCH_STATUS_UPGRADE_REQUIRED)

	// ErrIgnored is returned if the operation was ignored
	ErrIgnored = statusErr(C.NOTMUCH_STATUS_IGNORED)

	// ErrPathError is returned when there is a problem with the proposed path,
	// e.g. a relative path passed to a function expecting an absolute path.
	// TODO(kalbasit): this is currently on master. uncomment when released.
	// ErrPathError = statusErr(C.NOTMUCH_STATUS_PATH_ERROR)

	// ErrNotFound is returned when Find* did not find the thread/message by id or filename.
	ErrNotFound = errors.New("not found")

	// ErrUnknownError is returned when notmuch returns NULL indicating an error.
	ErrUnknownError = errors.New("unknown error occured")

	// ErrNoRepliesOrPointerNotFromThread is returned if a message has no replies or if the message's C
	// pointer did not come from a thread.
	ErrNoRepliesOrPointerNotFromThread = errors.New("message has no replies or message's pointer not from a thread")
)

// Notmuch returns NULL in several instances on out of memory errors. The
// expected go behavior is to panic. This function checks that if argument is nil
// and if so, panics with an out-of-memory message.
func checkOOM(ptr unsafe.Pointer) {
	if ptr == nil {
		panic(ErrOutOfMemory)
	}
}

// Convert a notmuch status to an error. This is almost a simple cast, but
// we need to return nil if it's a success, rather than NOTMUCH_STATUS_SUCCESS.
func statusErr(s C.notmuch_status_t) error {
	if s != C.NOTMUCH_STATUS_SUCCESS {
		return status(s)
	}
	return nil
}

func (s status) Error() string {
	cstr := C.notmuch_status_to_string(C.notmuch_status_t(s))
	return C.GoString(cstr)
}
