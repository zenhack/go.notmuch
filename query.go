package notmuch

// #cgo LDFLAGS: -lnotmuch
// #include <stdlib.h>
// #include <notmuch.h>
import "C"

import (
	"runtime"
	"unsafe"
)

type (
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
)

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

// Threads returns the threads matching the query.
func (q *Query) Threads() (*Threads, error) {
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

// CountThreads returns the number of messages for the current query.
func (q *Query) CountThreads() int {
	return int(C.notmuch_query_count_threads(q.toC()))
}

// CountMessages returns the number of messages for the current query.
func (q *Query) CountMessages() int {
	return int(C.notmuch_query_count_messages(q.toC()))
}

// Next retrieves the next thread from the result set. Next returns true if a thread
// was successfully retrieved.
func (ts *Threads) Next(t *Thread) bool {
	if !ts.valid() {
		return false
	}
	*t = *ts.get()
	C.notmuch_threads_move_to_next(ts.toC())
	return true
}

func (q *Query) toC() *C.notmuch_query_t {
	return (*C.notmuch_query_t)(q.cptr)
}

func (ts *Threads) toC() *C.notmuch_threads_t {
	return (*C.notmuch_threads_t)(ts.cptr)
}

// Get fetches the currently selected thread.
func (ts *Threads) get() *Thread {
	cthread := C.notmuch_threads_get(ts.toC())
	checkOOM(unsafe.Pointer(cthread))
	thread := &Thread{
		id:      C.GoString(C.notmuch_thread_get_thread_id(cthread)),
		cptr:    cthread,
		threads: ts,
	}
	runtime.SetFinalizer(thread, func(t *Thread) {
		C.notmuch_thread_destroy(t.toC())
	})
	return thread
}

func (ts *Threads) valid() bool {
	cbool := C.notmuch_threads_valid(ts.toC())
	return int(cbool) != 0
}
