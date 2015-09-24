package notmuch

// #cgo LDFLAGS: -lnotmuch
// #include <stdlib.h>
// #include <notmuch.h>
import "C"

// GetSubject returns the subject of a thread.
func (t *Thread) GetSubject() string {
	cstr := C.notmuch_thread_get_subject(t.toC())
	str := C.GoString(cstr)
	return str
}

// GetID returns the ID of the thread.
func (t *Thread) GetID() string {
	return t.id
}
