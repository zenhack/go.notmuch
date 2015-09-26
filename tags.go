package notmuch

// Copyright © 2015 The go.notmuch Authors. Authors can be found in the AUTHORS file.
// Licensed under the GPLv3 or later.
// See COPYING at the root of the repository for details.

// #cgo LDFLAGS: -lnotmuch
// #include <stdlib.h>
// #include <notmuch.h>
import "C"

// Tags represent a notmuch tags type.
type Tags struct {
	cptr   *C.notmuch_tags_t
	thread *Thread
}

// Next retrieves the next tag from the result set. Next returns true if a tag
// was successfully retrieved.
func (ts *Tags) Next(t *Tag) bool {
	if !ts.valid() {
		return false
	}
	*t = *ts.get()
	C.notmuch_tags_move_to_next(ts.cptr)
	return true
}

// Get fetches the currently selected tag.
func (ts *Tags) get() *Tag {
	ctag := C.notmuch_tags_get(ts.cptr)
	tag := &Tag{
		Value: C.GoString(ctag),
		tags:  ts,
	}
	return tag
}

func (ts *Tags) valid() bool {
	cbool := C.notmuch_tags_valid(ts.cptr)
	return int(cbool) != 0
}
