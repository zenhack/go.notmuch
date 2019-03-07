package notmuch

// Copyright © 2015 The go.notmuch Authors. Authors can be found in the AUTHORS file.
// Licensed under the GPLv3 or later.
// See COPYING at the root of the repository for details.

// #cgo LDFLAGS: -lnotmuch
// #include <stdlib.h>
// #include <notmuch.h>
import "C"

import (
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
	DB cStruct
)

func (db *DB) toC() *C.notmuch_database_t {
	return (*C.notmuch_database_t)(db.cptr)
}

// Close closes the database.
func (db *DB) Close() error {
	return (*cStruct)(db).doClose(func() error {
		return statusErr(C.notmuch_database_destroy(db.toC()))
	})
}

// Create creates a new, empty notmuch database located at 'path'.
func Create(path string) (*DB, error) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	var cdb *C.notmuch_database_t
	err := statusErr(C.notmuch_database_create(cpath, &cdb))
	if err != nil {
		return nil, err
	}
	db := &DB{cptr: unsafe.Pointer(cdb)}
	setGcClose(db)
	return db, nil
}

// Open opens the database at the location path using mode. Caller is responsible
// for closing the database when done.
func Open(path string, mode DBMode) (*DB, error) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	cmode := C.notmuch_database_mode_t(mode)
	var cdb *C.notmuch_database_t
	cdbptr := (**C.notmuch_database_t)(&cdb)
	err := statusErr(C.notmuch_database_open(cpath, cmode, cdbptr))
	if err != nil {
		return nil, err
	}
	db := &DB{cptr: unsafe.Pointer(cdb)}
	setGcClose(db)
	return db, nil
}

// Compact compacts a notmuch database, backing up the original database to the
// given path. The database will be opened with DBReadWrite to ensure no writes
// are made.
func Compact(path, backup string) error {
	cpath := C.CString(path)
	cbackup := C.CString(backup)
	defer func() {
		C.free(unsafe.Pointer(cpath))
		C.free(unsafe.Pointer(cbackup))
	}()

	return statusErr(C.notmuch_database_compact(cpath, cbackup, nil, nil))
}

// Atomic opens an atomic transaction in the database and calls the callback.
func (db *DB) Atomic(callback func(*DB)) error {
	if err := statusErr(C.notmuch_database_begin_atomic(db.toC())); err != nil {
		return err
	}
	callback(db)
	return statusErr(C.notmuch_database_end_atomic(db.toC()))
}

// NewQuery creates a new query from a string following xapian format.
func (db *DB) NewQuery(queryString string) *Query {
	cstr := C.CString(queryString)
	defer C.free(unsafe.Pointer(cstr))
	cquery := C.notmuch_query_create(db.toC(), cstr)
	query := &Query{
		cptr:   unsafe.Pointer(cquery),
		parent: (*cStruct)(db),
	}
	setGcClose(query)
	return query
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

// Upgrade upgrades the current database to the latest supported version. The
// database must be opened with DBReadWrite.
func (db *DB) Upgrade() error {
	return statusErr(C.notmuch_database_upgrade(db.toC(), nil, nil))
}

// AddMessage adds a new message to the current database or associate an
// additional filename with an existing message.
func (db *DB) AddMessage(filename string) (*Message, error) {
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))

	var cmsg *C.notmuch_message_t
	if err := statusErr(C.notmuch_database_index_file(db.toC(), cfilename, nil, &cmsg)); err != nil {
		return nil, err
	}
	msg := &Message{
		cptr:   unsafe.Pointer(cmsg),
		parent: (*cStruct)(db),
	}
	setGcClose(msg)
	return msg, nil
}

// RemoveMessage remove a message filename from the current database. If the
// message has no more filenames, remove the message.
func (db *DB) RemoveMessage(filename string) error {
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))

	return statusErr(C.notmuch_database_remove_message(db.toC(), cfilename))
}

// FindMessage finds a message with the given message_id.
func (db *DB) FindMessage(id string) (*Message, error) {
	cid := C.CString(id)
	defer C.free(unsafe.Pointer(cid))
	var cmsg *C.notmuch_message_t
	if err := statusErr(C.notmuch_database_find_message(db.toC(), cid, &cmsg)); err != nil {
		return nil, err
	}
	if cmsg == nil {
		return nil, ErrNotFound
	}
	msg := &Message{
		cptr:   unsafe.Pointer(cmsg),
		parent: (*cStruct)(db),
	}
	setGcClose(msg)
	return msg, nil
}

// FindMessageByFilename finds a message with the given filename.
func (db *DB) FindMessageByFilename(filename string) (*Message, error) {
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))

	var cmsg *C.notmuch_message_t
	if err := statusErr(C.notmuch_database_find_message_by_filename(db.toC(), cfilename, &cmsg)); err != nil {
		return nil, err
	}
	if cmsg == nil {
		return nil, ErrNotFound
	}
	msg := &Message{
		cptr:   unsafe.Pointer(cmsg),
		parent: (*cStruct)(db),
	}
	setGcClose(msg)
	return msg, nil
}

// Tags returns the list of all tags in the database.
func (db *DB) Tags() (*Tags, error) {
	ctags := C.notmuch_database_get_all_tags(db.toC())
	if ctags == nil {
		return nil, ErrUnknownError
	}
	tags := &Tags{
		cptr:   unsafe.Pointer(ctags),
		parent: (*cStruct)(db),
	}
	setGcClose(tags)
	return tags, nil
}
