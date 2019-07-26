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

// SortMode represents the sort behaviour of a query.
// One of SORT_{OLDEST_FIRST,NEWEST_FIRST,MESSAGE_ID,UNSORTED}
type SortMode C.notmuch_sort_t

// ExcludeMode represents the exclude behaviour of a query.
// One of EXCLUDE_{ALL,FLAG,TRUE,FALSE}.
type ExcludeMode C.notmuch_exclude_t

var (
	SORT_OLDEST_FIRST SortMode = C.NOTMUCH_SORT_OLDEST_FIRST
	SORT_NEWEST_FIRST SortMode = C.NOTMUCH_SORT_NEWEST_FIRST
	SORT_MESSAGE_ID   SortMode = C.NOTMUCH_SORT_MESSAGE_ID
	SORT_UNSORTED     SortMode = C.NOTMUCH_SORT_UNSORTED
)

var (
	// Query.Messages and Query.Threads will return all matching
	// messages/threads regardless of exclude status.
	// The exclude flag will be set for any excluded message that is
	// returned by Query.SearchMessages, and the thread counts
	// for threads returned by Query.Threads will be the
	// number of non-excluded messages/matches.
	EXCLUDE_FLAG ExcludeMode = C.NOTMUCH_EXCLUDE_FLAG
	// Query.Messages and Query.Threads will return all matching
	// messages/threads regardless of exclude status.
	// The exclude status is completely ignored.
	EXCLUDE_FALSE ExcludeMode = C.NOTMUCH_EXCLUDE_FALSE
	// Query.Messages will omit excluded messages from the results,
	// and Query.Threads will omit threads that match only in excluded messages.
	// Query.Threads will include all messages in threads that
	// match in at least one non-excluded message.
	EXCLUDE_TRUE ExcludeMode = C.NOTMUCH_EXCLUDE_TRUE
	// Query.Messages will omit excluded messages from the results,
	// and Query.Threads will omit threads that match only in excluded messages.
	// Query.Threads will omit excluded messages from all threads.
	EXCLUDE_ALL ExcludeMode = C.NOTMUCH_EXCLUDE_ALL
)

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

// Messages returns the messages matching the query.
func (q *Query) Messages() (*Messages, error) {
	var cmsgs *C.notmuch_messages_t
	err := statusErr(C.notmuch_query_search_messages(q.toC(), &cmsgs))
	if err != nil {
		return nil, err
	}
	msgs := &Messages{
		cptr:   unsafe.Pointer(cmsgs),
		parent: (*cStruct)(q),
	}
	setGcClose(msgs)
	return msgs, nil
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

// SetExcludeScheme is used to set the exclude scheme on a query.
func (q *Query) SetExcludeScheme(mode ExcludeMode) {
	cmode := C.notmuch_exclude_t(mode)
	C.notmuch_query_set_omit_excluded(q.toC(), cmode)
}

// AddTagExclude adds a tag that will be excluded from the query results by default.
// Note that this function returns ErrIgnored if the provided tag appears in the query
func (q *Query) AddTagExclude(tag string) error {
	ctag := C.CString(tag)
	defer C.free(unsafe.Pointer(ctag))
	return statusErr(C.notmuch_query_add_tag_exclude(q.toC(), ctag))
}
