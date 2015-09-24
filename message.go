package notmuch

// #cgo LDFLAGS: -lnotmuch
// #include <stdlib.h>
// #include <notmuch.h>
import "C"

// Message represents a notmuch message.
type Message struct {
	cptr     *C.notmuch_message_t
	messages *Messages
	thread   *Thread
}

func (m *Message) toC() *C.notmuch_message_t {
	return (*C.notmuch_message_t)(m.cptr)
}
