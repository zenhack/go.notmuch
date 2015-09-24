package notmuch

import (
	"path/filepath"
	"testing"
)

var dbPath string

func init() {
	var err error
	dbPath, err = filepath.Abs("fixtures/database-v1")
	if err != nil {
		panic(err)
	}
}

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
