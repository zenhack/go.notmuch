package notmuch

// Copyright © 2015 The go.notmuch Authors. Authors can be found in the AUTHORS file.
// Licensed under the GPLv3 or later.
// See COPYING at the root of the repository for details.

// #cgo LDFLAGS: -lnotmuch
// #include <stdlib.h>
// #include <notmuch.h>
import "C"

// Tag represents a tag in the database.
type Tag struct {
	Value string
	tags  *Tags
}

func (t *Tag) String() string {
	return t.Value
}
