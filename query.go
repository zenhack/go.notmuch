package notmuch

// Copyright © 2015 The go.notmuch Authors. Authors can be found in the AUTHORS file.
// Licensed under the GPLv3 or later.
// See COPYING at the root of the repository for details.

// #cgo LDFLAGS: -lnotmuch
// #include <stdlib.h>
// #include <notmuch.h>
import "C"
import "unsafe"

// Query represents a notmuch query.
type Query struct {
	qs   string
	cptr *C.notmuch_query_t
	db   *DB
}

// String returns the query as a string, implements fmt.Stringer.
func (q *Query) String() string {
	return q.qs
}

// Threads returns the threads matching the query.
func (q *Query) Threads() (*Threads, error) {
	threads := &Threads{query: q}
	cthreads := (**C.notmuch_threads_t)(unsafe.Pointer(&threads.cptr))
	cerr := C.notmuch_query_search_threads_st(q.cptr, cthreads)
	err := statusErr(cerr)
	if err != nil {
		return nil, err
	}
	return threads, nil
}

// CountThreads returns the number of messages for the current query.
func (q *Query) CountThreads() int {
	return int(C.notmuch_query_count_threads(q.cptr))
}

// CountMessages returns the number of messages for the current query.
func (q *Query) CountMessages() int {
	return int(C.notmuch_query_count_messages(q.cptr))
}

// Destroy a Query along with any associated resources.
//
// This will in turn destroy any Threads and Messages objects generated
// by this query, (and in turn any Thread and Message objects generated
// from those results, etc.), if such objects haven't already been
// destroyed.
func (q *Query) Close() error {
	C.notmuch_query_destroy(q.cptr)
	return nil
}
