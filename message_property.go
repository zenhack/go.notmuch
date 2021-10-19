package notmuch

// Copyright © 2015 The go.notmuch Authors. Authors can be found in the AUTHORS file.
// Licensed under the GPLv3 or later.
// See COPYING at the root of the repository for details.

// #cgo LDFLAGS: -lnotmuch
// #include <stdlib.h>
// #include <notmuch.h>
import "C"

// MessageProperty represents a property in the database.
type MessageProperty struct {
	Key        string
	Value      string
	properties *MessageProperties
}

func (p *MessageProperty) String() string {
	return p.Key + "=" + p.Value
}
