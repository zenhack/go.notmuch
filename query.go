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

// Query represents a notmuch query.
type Query struct {
	cptr *C.notmuch_query_t
	db   *DB
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
	runtime.SetFinalizer(threads, func(t *Threads) {
		C.notmuch_threads_destroy(t.cptr)
	})
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
