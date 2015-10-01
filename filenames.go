package notmuch

// Copyright © 2015 The go.notmuch Authors. Authors can be found in the AUTHORS file.
// Licensed under the GPLv3 or later.
// See COPYING at the root of the repository for details.

// #cgo LDFLAGS: -lnotmuch
// #include <stdlib.h>
// #include <notmuch.h>
import "C"

// Filenames is an iterator to get the message's filenames.
type Filenames struct {
	cptr    *C.notmuch_filenames_t
	message *Message
}

// Next retrieves the next filename from the iterator. Next returns true if a
// filename was successfully retrieved.
func (fs *Filenames) Next(f *string) bool {
	if !fs.valid() {
		return false
	}
	*f = fs.get()
	C.notmuch_filenames_move_to_next(fs.cptr)
	return true
}

func (fs *Filenames) get() string {
	return C.GoString(C.notmuch_filenames_get(fs.cptr))
}

func (fs *Filenames) valid() bool {
	cbool := C.notmuch_filenames_valid(fs.cptr)
	return int(cbool) != 0
}
