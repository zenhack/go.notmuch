package notmuch

// Copyright © 2015 The go.bindings Authors. Authors can be found in the AUTHORS file.
// Licensed under the GPLv3 or later.
// See COPYING at the root of the repository for details.

import "testing"

func TestSearchThreads(t *testing.T) {
	db, err := Open(dbPath, DBReadOnly)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	threads, err := db.NewQuery("").Threads()
	if err != nil {
		t.Fatalf("error getting the threads: %s", err)
	}

	var count int
	thread := &Thread{}
	for threads.Next(thread) {
		count++
	}

	if want, got := 24, count; want != got {
		t.Errorf("db.NewQuery(%q).Threads(): want %d got %d", "", want, got)
	}
}

func TestGetNoResult(t *testing.T) {
	db, err := Open(dbPath, DBReadOnly)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	threads, err := db.NewQuery("subject:notfoundnotfound").Threads()
	if err != nil {
		t.Fatalf("error getting the threads: %s", err)
	}

	var count int
	thread := &Thread{}
	for threads.Next(thread) {
		count++
	}

	if want, got := 0, count; want != got {
		t.Errorf("db.NewQuery(%q).Threads(): want %d got %d", "", want, got)
	}
}

func TestQueryCountMessages(t *testing.T) {
	db, err := Open(dbPath, DBReadOnly)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	q := db.NewQuery("subject:\"Introducing myself\"")
	if want, got := 3, q.CountMessages(); want != got {
		t.Errorf("q.Count(): want %d got %d", want, got)
	}
}

func TestQueryCountThreads(t *testing.T) {
	db, err := Open(dbPath, DBReadOnly)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	q := db.NewQuery("subject:\"Introducing myself\"")
	if want, got := 1, q.CountThreads(); want != got {
		t.Errorf("q.Count(): want %d got %d", want, got)
	}
}
