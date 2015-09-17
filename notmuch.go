// Package notmuch provides a Go language binding to the notmuch mail library.
//
// The design is similar enough to the underlying C library that familiarity with
// one will inform the other. There are some differences, however:
//
// * The Go binding arranges for the garbage collector to deal with objects
//   allocated by notmuch correctly. You should close the database manually,i
//   but everything else will be garbage collected when it becomes unreachable,
//   and not before. Objects hold references to their parent objects to make this
//   go smoothly.
// * If notmuch returns NULL because of an out of memory error, Go will panic, as
//   it does with other out of memory errors.
// * Some of the names have been shortened or made more idiomatic. The documentation
//   indends to make it obvious when this is the case.
// * Functions which create a child object from a parent object are methods on the
//   parent object, rather than stand-alone functions.
// * Functions which in C return a status code and pass back a value via a pointer
//   argument now return a (value, error) pair.

package notmuch
// Copyright © 2015 Ian Denhardt <ian@zenhack.net>
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
	DB_RDONLY = C.NOTMUCH_DATABASE_MODE_READ_ONLY
	DB_RDWR = C.NOTMUCH_DATABASE_MODE_READ_WRITE
)

type status C.notmuch_status_t
type DBMode C.notmuch_database_mode_t

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
	if s == C.NOTMUCH_STATUS_SUCCESS {
		return nil
	} else {
		return status(s)
	}
}

func (s status) Error() string {
	cstr := C.notmuch_status_to_string(C.notmuch_status_t(s))
	return C.GoString(cstr)
}

type DB struct {
	cptr unsafe.Pointer
}

type Query struct {
	cptr *C.notmuch_query_t
	db *DB
}

type Threads struct {
	cptr *C.notmuch_threads_t
	query *Query
}

type Thread struct {
	cptr *C.notmuch_thread_t
	threads *Threads
}

func (db *DB) toC() *C.notmuch_database_t {
	return (*C.notmuch_database_t)(db.cptr)
}

func (q *Query) toC() *C.notmuch_query_t {
	return (*C.notmuch_query_t)(q.cptr)
}

func (t *Threads) toC() *C.notmuch_threads_t {
	return (*C.notmuch_threads_t)(t.cptr)
}

func (t *Thread) toC() *C.notmuch_thread_t {
	return (*C.notmuch_thread_t)(t.cptr)
}

func Open(path string, mode DBMode) (*DB, error) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	cmode := C.notmuch_database_mode_t(mode)

	db  := &DB{}
	cdb := (**C.notmuch_database_t)(unsafe.Pointer(&db.cptr))

	cerr := C.notmuch_database_open(cpath, cmode, cdb)

	err := statusErr(cerr)
	return db, err
}

func (db *DB) Close() error {
	cdb := (*C.notmuch_database_t)(db.cptr)
	cerr := C.notmuch_database_close(cdb)
	err := statusErr(cerr)
	return err
}

func (db *DB) QueryCreate(queryString string) *Query {
	cstr := C.CString(queryString)
	defer C.free(unsafe.Pointer(cstr))
	cquery := C.notmuch_query_create(db.toC(), cstr)
	query := &Query{
		cptr: cquery,
		db: db,
	}
	runtime.SetFinalizer(query, func(q *Query) {
		C.notmuch_query_destroy(q.toC())
	})
	return query
}

func (q *Query) SearchThreads() (*Threads, error) {
	threads := &Threads{query: q}
	cthreads := (**C.notmuch_threads_t)(unsafe.Pointer(&threads.cptr))
	cerr := C.notmuch_query_search_threads_st(q.toC(), cthreads)
	err := statusErr(cerr)
	if err != nil {
		return nil, err
	}
	runtime.SetFinalizer(threads, func(t *Threads) {
		C.notmuch_threads_destroy(t.toC())
	})
	return threads, nil
}

func (t *Threads) Valid() bool {
	cbool := C.notmuch_threads_valid(t.toC())
	return int(cbool) != 0
}


func (t *Threads) Get() *Thread {
	cthread := C.notmuch_threads_get(t.toC())
	// NOTE: we don't distinguish between OOM and calling Get when
	// !t.Valid(). As such, it's an error for the user to call Get
	// without first calling Valid.
	checkOOM(unsafe.Pointer(cthread))
	thread := &Thread{
		cptr: cthread,
		threads: t,
	}
	runtime.SetFinalizer(thread, func(t *Thread) {
		C.notmuch_thread_destroy(t.toC())
	})
	return thread
}

func (t *Threads) MoveToNext() {
	C.notmuch_threads_move_to_next(t.toC())
}

func (t *Thread) GetSubject() string {
	cstr := C.notmuch_thread_get_subject(t.toC())
	str := C.GoString(cstr)
	return str
}
