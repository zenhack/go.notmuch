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

const (
	SORT_OLDEST_FIRST = C.NOTMUCH_SORT_OLDEST_FIRST
	SORT_NEWEST_FIRST = C.NOTMUCH_SORT_NEWEST_FIRST
	SORT_MESSAGE_ID   = C.NOTMUCH_SORT_MESSAGE_ID
	SORT_UNSORTED     = C.NOTMUCH_SORT_UNSORTED
)

// SortMode represents the sort behaviour of a query.
// One of SORT_{OLDEST_FIRST,NEWEST_FIRST,MESSAGE_ID,UNSORTED}
type SortMode C.notmuch_sort_t

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
	err := statusErr(C.notmuch_query_search_threads(q.toC(), &cthreads))
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
	var ccount C.uint
	C.notmuch_query_count_threads(q.toC(), &ccount)
	return int(ccount)
}

// CountMessages returns the number of messages for the current query.
func (q *Query) CountMessages() int {
	var cCount C.uint
	C.notmuch_query_count_messages(q.toC(), &cCount)
	return int(cCount)
}

// SetSortScheme is used to set the sort scheme on a query.
func (q *Query) SetSortScheme(mode SortMode) {
	cmode := C.notmuch_sort_t(mode)
	C.notmuch_query_set_sort(q.toC(), cmode)
}
