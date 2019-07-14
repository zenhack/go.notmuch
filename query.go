package notmuch

// Copyright © 2015 The go.notmuch Authors. Authors can be found in the AUTHORS file.
// Licensed under the GPLv3 or later.
// See COPYING at the root of the repository for details.

// #cgo LDFLAGS: -lnotmuch
// #include <stdlib.h>
// #include <notmuch.h>
import "C"

import "unsafe"

const (
	// notmuch_query_search_messages and notmuch_query_search_threads
	// will return all matching messages/threads regardless of exclude status.
	// The exclude flag will be set for any excluded message that is
	// returned by notmuch_query_search_messages, and the thread counts
	// for threads returned by notmuch_query_search_threads will be the
	// number of non-excluded messages/matches.
	QUERY_EXCLUDE_FLAG = C.NOTMUCH_EXCLUDE_FLAG
	// notmuch_query_search_messages and notmuch_query_search_threads
	// will return all matching messages/threads regardless of exclude status.
	// The exclude status is completely ignored.
	QUERY_EXCLUDE_FALSE = C.NOTMUCH_EXCLUDE_FALSE
	// notmuch_query_search_messages will omit excluded
	// messages from the results, and notmuch_query_search_threads will omit
	// threads that match only in excluded messages.
	// notmuch_query_search_threads will include all messages in threads that
	// match in at least one non-excluded message.
	QUERY_EXCLUDE_TRUE = C.NOTMUCH_EXCLUDE_TRUE
	// notmuch_query_search_messages will omit excluded
	// messages from the results, and notmuch_query_search_threads will omit
	// threads that match only in excluded messages.
	// notmuch_query_search_threads will omit excluded messages from all threads.
	QUERY_EXCLUDE_ALL = C.NOTMUCH_EXCLUDE_ALL
)

// Query represents a notmuch query.
type Query cStruct

// ExcludeMode represents the exclude behaviour of a query.
// One of QUERY_EXCLUDE_{ALL,FLAG,TRUE,FALSE}.
type ExcludeMode C.notmuch_exclude_t

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

// SetExcludeScheme is used to set the exclude scheme on a query.
func (q *Query) SetExcludeScheme(mode ExcludeMode) {
	cmode := C.notmuch_exclude_t(mode)
	C.notmuch_query_set_omit_excluded(q.toC(), cmode)
}
