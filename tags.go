package notmuch

// Copyright © 2015 The go.notmuch Authors. Authors can be found in the AUTHORS file.
// Licensed under the GPLv3 or later.
// See COPYING at the root of the repository for details.

// #cgo LDFLAGS: -lnotmuch
// #include <stdlib.h>
// #include <notmuch.h>
import "C"

// Tags represent a notmuch tags type.
type Tags cStruct

func (ts *Tags) toC() *C.notmuch_tags_t {
	return (*C.notmuch_tags_t)(ts.cptr)
}

func (ts *Tags) Close() error {
	return (*cStruct)(ts).doClose(func() error {
		C.notmuch_tags_destroy(ts.toC())
		return nil
	})
}

// Next retrieves the next tag from the result set. Next returns true if a tag
// was successfully retrieved.
func (ts *Tags) Next(t *Tag) bool {
	if !ts.valid() {
		return false
	}
	*t = *ts.get()
	C.notmuch_tags_move_to_next(ts.toC())
	return true
}

// Return a slice of strings containing each element of ts.
func (ts *Tags) slice() []string {
	tag := &Tag{}
	ret := []string{}
	for ts.Next(tag) {
		ret = append(ret, tag.Value)
	}
	return ret
}

func (ts *Tags) get() *Tag {
	ctag := C.notmuch_tags_get(ts.toC())
	tag := &Tag{
		Value: C.GoString(ctag),
		tags:  ts,
	}
	return tag
}

func (ts *Tags) valid() bool {
	cbool := C.notmuch_tags_valid(ts.toC())
	return int(cbool) != 0
}
