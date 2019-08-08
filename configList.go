package notmuch

// Copyright © 2015 The go.notmuch Authors. Authors can be found in the AUTHORS file.
// Licensed under the GPLv3 or later.
// See COPYING at the root of the repository for details.

// #cgo LDFLAGS: -lnotmuch
// #include <stdlib.h>
// #include <notmuch.h>
import "C"

// Messages represents a notmuch config list
type ConfigList cStruct

func (cl *ConfigList) Close() error {
	return (*cStruct)(cl).doClose(func() error {
		C.notmuch_config_list_destroy(cl.toC())
		return nil
	})
}

func (cl *ConfigList) toC() *C.notmuch_config_list_t {
	return (*C.notmuch_config_list_t)(cl.cptr)
}

// Next retrieves the nex config pair from the ConfigList.
// Neither key, nor value may be nil, or this function will panic.
// Next returns true if a pair was successfully retrieved.
func (cl *ConfigList) Next(key, value *string) bool {
	if !cl.valid() {
		return false
	}
	*key = cl.key()
	*value = cl.value()
	C.notmuch_config_list_move_to_next(cl.toC())
	return true
}

func (cl *ConfigList) valid() bool {
	cbool := C.notmuch_config_list_valid(cl.toC())
	return int(cbool) != 0
}

func (cl *ConfigList) key() string {
	return C.GoString(C.notmuch_config_list_key(cl.toC()))
}

func (cl *ConfigList) value() string {
	return C.GoString(C.notmuch_config_list_value(cl.toC()))
}
