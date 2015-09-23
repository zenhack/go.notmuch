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
	for {
		threads.Get()
		count++

		if err := threads.MoveToNext(); err != nil {
			break
		}
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
	for {
		if _, err := threads.Get(); err != nil {
			break
		}
		count++
		threads.MoveToNext()
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
	thread, err := threads.Get()
	if err != nil {
		t.Errorf("threads.Get(): got error: %s", err)
	}
	if want, got := "Essai accentué", thread.GetSubject(); want != got {
		t.Errorf("db.NewQuery(%q).Threads().Get().GetSubject(): want %s got %s", want, want, got)
	}
}
