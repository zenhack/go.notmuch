package notmuch

// Copyright © 2015 The go.notmuch Authors. Authors can be found in the AUTHORS file.
// Licensed under the GPLv3 or later.
// See COPYING at the root of the repository for details.

import "testing"

func TestGetThreadID(t *testing.T) {
	db, err := Open(dbPath, DBReadOnly)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	threads, err := db.NewQuery("Essai accentué").Threads()
	if err != nil {
		t.Fatalf("error getting the threads: %s", err)
	}
	thread := &Thread{}
	if !threads.Next(thread) {
		t.Fatalf("threads.Next(thread): unable to fetch the first and only thread")
	}
	if want, got := "0000000000000014", thread.GetID(); want != got {
		t.Errorf("db.NewQuery(%q).Threads().Get().GetID(): want %s got %s", want, want, got)
	}
}

func TestGetSubjectUTF8(t *testing.T) {
	db, err := Open(dbPath, DBReadOnly)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	threads, err := db.NewQuery("Essai accentué").Threads()
	if err != nil {
		t.Fatalf("error getting the threads: %s", err)
	}
	thread := &Thread{}
	if !threads.Next(thread) {
		t.Fatalf("threads.Next(thread): unable to fetch the first and only thread")
	}
	if want, got := "Essai accentué", thread.GetSubject(); want != got {
		t.Errorf("db.NewQuery(%q).Threads().Get().GetSubject(): want %s got %s", want, want, got)
	}
}
