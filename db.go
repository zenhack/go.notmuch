package notmuch

// Copyright © 2015 The go.bindings Authors. Authors can be found in the AUTHORS file.
// Licensed under the GPLv3 or later.
// See COPYING at the root of the repository for details.

// #cgo LDFLAGS: -lnotmuch
// #include <stdlib.h>
// #include <notmuch.h>
import "C"

import (
	"runtime"
	"unsafe"
)

const (
	// DBReadOnly is the mode for opening the database in read only.
	DBReadOnly = C.NOTMUCH_DATABASE_MODE_READ_ONLY

	// DBReadWrite is the mode for opening the database in read write.
	DBReadWrite = C.NOTMUCH_DATABASE_MODE_READ_WRITE

	// TagMax is the maximum number of allowed tags.
	TagMax = C.NOTMUCH_TAG_MAX
)

type (
	// DBMode is the mode of the database opening, DBReadOnly or DBReadWrite
	DBMode C.notmuch_database_mode_t

	// DB represents a notmuch database.
	DB struct {
		cptr unsafe.Pointer
	}
)

// Create creates a new, empty notmuch database located at 'path'.
func Create(path string) (*DB, error) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	db := &DB{}
	cdb := (**C.notmuch_database_t)(unsafe.Pointer(&db.cptr))
	cerr := C.notmuch_database_create(cpath, cdb)
	runtime.SetFinalizer(db, func(db *DB) {
		C.notmuch_database_destroy(db.toC())
	})
	return db, statusErr(cerr)
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
	runtime.SetFinalizer(db, func(db *DB) {
		C.notmuch_database_destroy(db.toC())
	})
	return db, statusErr(cerr)
}

// NewQuery creates a new query from a string following xapian format.
func (db *DB) NewQuery(queryString string) *Query {
	cstr := C.CString(queryString)
	defer C.free(unsafe.Pointer(cstr))
	cquery := C.notmuch_query_create(db.toC(), cstr)
	query := &Query{
		cptr: cquery,
		db:   db,
	}
	runtime.SetFinalizer(query, func(q *Query) {
		C.notmuch_query_destroy(q.toC())
	})
	return query
}

// Close closes the database.
func (db *DB) Close() error {
	cerr := C.notmuch_database_close(db.toC())
	err := statusErr(cerr)
	return err
}

// Version returns the database version.
func (db *DB) Version() int {
	return int(C.notmuch_database_get_version(db.toC()))
}

// LastStatus retrieves last status string for the notmuch database.
func (db *DB) LastStatus() string {
	return C.GoString(C.notmuch_database_status_string(db.toC()))
}

// Path returns the database path of the database.
func (db *DB) Path() string {
	return C.GoString(C.notmuch_database_get_path(db.toC()))
}

// NeedsUpgrade returns true if the database can be upgraded. This will always
// return false if the database was opened with DBReadOnly.
//
// If this function returns true then the caller may call
// Upgrade() to upgrade the database.
func (db *DB) NeedsUpgrade() bool {
	cbool := C.notmuch_database_needs_upgrade(db.toC())
	return int(cbool) != 0
}

func (db *DB) toC() *C.notmuch_database_t {
	return (*C.notmuch_database_t)(db.cptr)
}
