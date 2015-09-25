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
		cptr *C.notmuch_database_t
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
		C.notmuch_database_destroy(db.cptr)
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
		C.notmuch_database_destroy(db.cptr)
	})
	return db, statusErr(cerr)
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
	cerr := C.notmuch_database_begin_atomic(db.cptr)
	if err := statusErr(cerr); err != nil {
		return err
	}
	callback(db)
	return statusErr(C.notmuch_database_end_atomic(db.cptr))
}

// NewQuery creates a new query from a string following xapian format.
func (db *DB) NewQuery(queryString string) *Query {
	cstr := C.CString(queryString)
	defer C.free(unsafe.Pointer(cstr))
	cquery := C.notmuch_query_create(db.cptr, cstr)
	query := &Query{
		cptr: cquery,
		db:   db,
	}
	runtime.SetFinalizer(query, func(q *Query) {
		C.notmuch_query_destroy(q.cptr)
	})
	return query
}

// Close closes the database.
func (db *DB) Close() error {
	cerr := C.notmuch_database_close(db.cptr)
	err := statusErr(cerr)
	return err
}

// Version returns the database version.
func (db *DB) Version() int {
	return int(C.notmuch_database_get_version(db.cptr))
}

// LastStatus retrieves last status string for the notmuch database.
func (db *DB) LastStatus() string {
	return C.GoString(C.notmuch_database_status_string(db.cptr))
}

// Path returns the database path of the database.
func (db *DB) Path() string {
	return C.GoString(C.notmuch_database_get_path(db.cptr))
}

// NeedsUpgrade returns true if the database can be upgraded. This will always
// return false if the database was opened with DBReadOnly.
//
// If this function returns true then the caller may call
// Upgrade() to upgrade the database.
func (db *DB) NeedsUpgrade() bool {
	cbool := C.notmuch_database_needs_upgrade(db.cptr)
	return int(cbool) != 0
}

// Upgrade upgrades the current database to the latest supported version. The
// database must be opened with DBReadWrite.
func (db *DB) Upgrade() error {
	return statusErr(C.notmuch_database_upgrade(db.cptr, nil, nil))
}

// AddMessage adds a new message to the current database or associate an
// additional filename with an existing message.
func (db *DB) AddMessage(filename string) (*Message, error) {
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))

	msg := &Message{}
	cmsg := (**C.notmuch_message_t)(unsafe.Pointer(&msg.cptr))
	if err := statusErr(C.notmuch_database_add_message(db.cptr, cfilename, cmsg)); err != nil {
		return nil, err
	}
	runtime.SetFinalizer(msg, func(m *Message) {
		C.notmuch_message_destroy(m.cptr)
	})
	return msg, nil
}

// RemoveMessage remove a message filename from the current database. If the
// message has no more filenames, remove the message.
func (db *DB) RemoveMessage(filename string) error {
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))

	return statusErr(C.notmuch_database_remove_message(db.cptr, cfilename))
}

// FindMessage finds a message with the given message_id.
func (db *DB) FindMessage(id string) (*Message, error) {
	cid := C.CString(id)
	defer C.free(unsafe.Pointer(cid))

	msg := &Message{}
	cmsg := (**C.notmuch_message_t)(unsafe.Pointer(&msg.cptr))
	if err := statusErr(C.notmuch_database_find_message(db.cptr, cid, cmsg)); err != nil {
		return nil, err
	}
	if msg.cptr == nil {
		return nil, ErrNotFound
	}
	runtime.SetFinalizer(msg, func(m *Message) {
		C.notmuch_message_destroy(m.cptr)
	})
	return msg, nil
}

// FindMessageByFilename finds a message with the given filename.
func (db *DB) FindMessageByFilename(filename string) (*Message, error) {
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))

	msg := &Message{}
	cmsg := (**C.notmuch_message_t)(unsafe.Pointer(&msg.cptr))
	if err := statusErr(C.notmuch_database_find_message_by_filename(db.cptr, cfilename, cmsg)); err != nil {
		return nil, err
	}
	if msg.cptr == nil {
		return nil, ErrNotFound
	}
	runtime.SetFinalizer(msg, func(m *Message) {
		C.notmuch_message_destroy(m.cptr)
	})
	return msg, nil
}

// Tags returns the list of all tags in the database.
func (db *DB) Tags() (*Tags, error) {
	ctags := C.notmuch_database_get_all_tags(db.cptr)
	if ctags == nil {
		return nil, ErrUnknownError
	}
	tags := &Tags{
		cptr: ctags,
	}
	runtime.SetFinalizer(tags, func(t *Tags) {
		C.notmuch_tags_destroy(tags.cptr)
	})

	return tags, nil
}
