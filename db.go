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
)

func (db *DB) toC() *C.notmuch_database_t {
	return (*C.notmuch_database_t)(db.cptr)
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

// Version returns the database version.
func (db *DB) Version() int {
	return int(C.notmuch_database_get_version(db.toC()))
}
