package notmuch

// #cgo LDFLAGS: -lnotmuch
// #include <stdlib.h>
// #include <notmuch.h>
import "C"

import "unsafe"

const (
	// DBReadOnly is the mode for opening the database in read only.
	DBReadOnly = C.NOTMUCH_DATABASE_MODE_READ_ONLY

	// DBReadWrite is the mode for opening the database in read write.
	DBReadWrite = C.NOTMUCH_DATABASE_MODE_READ_WRITE
)

type (
	// DBMode is the mode of the database opening, DBReadOnly or DBReadWrite
	DBMode C.notmuch_database_mode_t

	// DB represents a notmuch database.
	DB struct {
		cptr unsafe.Pointer
	}

	// Query represents a notmuch query.
	Query struct {
		cptr *C.notmuch_query_t
		db   *DB
	}

	// Threads represents notmuch threads.
	Threads struct {
		cptr  *C.notmuch_threads_t
		query *Query
	}

	// Thread represents a notmuch thread.
	Thread struct {
		id      string
		cptr    *C.notmuch_thread_t
		threads *Threads
	}

	// Messages represents notmuch messages.
	Messages struct {
		cptr   *C.notmuch_messages_t
		thread *Thread
	}

	// Message represents a notmuch message.
	Message struct {
		cptr     *C.notmuch_message_t
		messages *Messages
		thread   *Thread
	}

	status C.notmuch_status_t
)

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

func (db *DB) toC() *C.notmuch_database_t {
	return (*C.notmuch_database_t)(db.cptr)
}

func (ts *Threads) toC() *C.notmuch_threads_t {
	return (*C.notmuch_threads_t)(ts.cptr)
}

func (t *Thread) toC() *C.notmuch_thread_t {
	return (*C.notmuch_thread_t)(t.cptr)
}

func (m *Message) toC() *C.notmuch_message_t {
	return (*C.notmuch_message_t)(m.cptr)
}

// Open opens the database at the location path using mode. Caller is responsible
// for closing the database when done.
func Open(path string, mode DBMode) (*DB, error) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	cmode := C.notmuch_database_mode_t(mode)

	db := &DB{}
	cdb := (**C.notmuch_database_t)(unsafe.Pointer(&db.cptr))

	cerr := C.notmuch_database_open(cpath, cmode, cdb)

	err := statusErr(cerr)
	return db, err
}

// Close closes the database.
func (db *DB) Close() error {
	cdb := (*C.notmuch_database_t)(db.cptr)
	cerr := C.notmuch_database_close(cdb)
	err := statusErr(cerr)
	return err
}
