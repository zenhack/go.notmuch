package notmuch

// Copyright © 2015 The go.bindings Authors. Authors can be found in the AUTHORS file.
// Licensed under the GPLv3 or later.
// See COPYING at the root of the repository for details.

// #cgo LDFLAGS: -lnotmuch
// #include <stdlib.h>
// #include <notmuch.h>
import "C"

import "unsafe"

type status C.notmuch_status_t

// Notmuch returns NULL in several instances on out of memory errors. The
// expected go behavior is to panic. This function checks that if argument is nil
// and if so, panics with an out-of-memory message.
func checkOOM(ptr unsafe.Pointer) {
	if ptr == nil {
		panic("Notmuch reported an out of memory error!")
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
