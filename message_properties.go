package notmuch

// Copyright © 2015 The go.notmuch Authors. Authors can be found in the AUTHORS file.
// Licensed under the GPLv3 or later.
// See COPYING at the root of the repository for details.

// #cgo LDFLAGS: -lnotmuch
// #include <stdlib.h>
// #include <notmuch.h>
import "C"

// MessageProperties represent a notmuch properties type.
type MessageProperties cStruct

func (props *MessageProperties) toC() *C.notmuch_message_properties_t {
	return (*C.notmuch_message_properties_t)(props.cptr)
}

func (props *MessageProperties) Close() error {
	return (*cStruct)(props).doClose(func() error {
		C.notmuch_message_properties_destroy(props.toC())
		return nil
	})
}

// Next retrieves the next prop from the result set. Next returns true if a prop
// was successfully retrieved.
func (props *MessageProperties) Next(p **MessageProperty) bool {
	if !props.valid() {
		return false
	}
	*p = props.get()
	C.notmuch_message_properties_move_to_next(props.toC())
	return true
}

// Return a slice of strings containing each element of props.
func (props *MessageProperties) slice() []string {
	var prop *MessageProperty
	ret := []string{}
	for props.Next(&prop) {
		ret = append(ret, prop.Value)
	}
	return ret
}

func (props *MessageProperties) get() *MessageProperty {
	ckey := C.notmuch_message_properties_key(props.toC())
	cvalue := C.notmuch_message_properties_value(props.toC())

	prop := &MessageProperty{
		Key:        C.GoString(ckey),
		Value:      C.GoString(cvalue),
		properties: props,
	}
	return prop
}

func (props *MessageProperties) valid() bool {
	cbool := C.notmuch_message_properties_valid(props.toC())
	return int(cbool) != 0
}
