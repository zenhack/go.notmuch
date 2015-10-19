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
type Query cStruct

func (q *Query) Close() error {
	return (*cStruct)(q).doClose(func() error {
		C.notmuch_query_destroy(q.toC())
		return nil
	})
}

func (q *Query) toC() *C.notmuch_query_t {
	return (*C.notmuch_query_t)(q.cptr)
}

// String returns the query as a string, implements fmt.Stringer.
func (q *Query) String() string {
	return C.GoString(C.notmuch_query_get_query_string(q.toC()))
}

// Threads returns the threads matching the query.
func (q *Query) Threads() (*Threads, error) {
	var cthreads *C.notmuch_threads_t
	err := statusErr(C.notmuch_query_search_threads_st(q.toC(), &cthreads))
	if err != nil {
		return nil, err
	}
	threads := &Threads{
		cptr:   unsafe.Pointer(cthreads),
		parent: (*cStruct)(q),
	}
	setGcClose(threads)
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
