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

//ConfigOption is used as a container for key / value pairs of the ConfigList
type ConfigOption struct {
	Key, Value string
}

func (cl *ConfigList) Close() error {
	return (*cStruct)(cl).doClose(func() error {
		C.notmuch_config_list_destroy(cl.toC())
		return nil
	})
}

func (cl *ConfigList) toC() *C.notmuch_config_list_t {
	return (*C.notmuch_config_list_t)(cl.cptr)
}

// Next retrieves the nex ConfigOption from the ConfigList.
// Next returns true if a ConfigOption was successfully retrieved.
func (cl *ConfigList) Next(opt **ConfigOption) bool {
	if !cl.valid() {
		return false
	}
	*opt = cl.getOption()
	C.notmuch_config_list_move_to_next(cl.toC())
	return true
}

func (cl *ConfigList) valid() bool {
	cbool := C.notmuch_config_list_valid(cl.toC())
	return int(cbool) != 0
}

func (cl *ConfigList) getOption() *ConfigOption {
	var opt ConfigOption
	opt.Key = cl.key()
	opt.Value = cl.value()
	return &opt
}

func (cl *ConfigList) key() string {
	return C.GoString(C.notmuch_config_list_key(cl.toC()))
}

func (cl *ConfigList) value() string {
	return C.GoString(C.notmuch_config_list_value(cl.toC()))
}
