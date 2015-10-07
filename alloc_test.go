package notmuch

// Copyright © 2015 The go.notmuch Authors. Authors can be found in the AUTHORS file.
// Licensed under the GPLv3 or later.
// See COPYING at the root of the repository for details.

// This file contains tests that exercise the resource creation/reclamation model.

import (
	"testing"
)

// Basic test of the case where the user explicitly closes a parent object
// before the child object. This is inevitable if Close is used at all; the
// GC will close things child-first.
func TestOutOfOrderClose(t *testing.T) {
	db, err := Open(dbPath, DBReadOnly)
	if err != nil {
		t.Fatalf("Open(%q): unexpected error: %s", dbPath, err)
	}

	query := db.NewQuery("subject:\"Introducing myself\"")
	threads, err := query.Threads()
	if err != nil {
		t.Fatalf("error getting the threads: %s", err)
	}

	db.Close()
	threads.Close()
}
