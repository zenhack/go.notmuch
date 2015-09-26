package notmuch

// Copyright © 2015 The go.notmuch Authors. Authors can be found in the AUTHORS file.
// Licensed under the GPLv3 or later.
// See COPYING at the root of the repository for details.

import (
	"reflect"
	"testing"
	"time"
)

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
		t.Errorf("db.NewQuery(%q).Threads()[0].GetID(): want %s got %s", "Essai accentué", want, got)
	}
}

func TestCount(t *testing.T) {
	db, err := Open(dbPath, DBReadOnly)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	qs := "subject:\"Introducing myself\" Hello"
	threads, err := db.NewQuery(qs).Threads()
	if err != nil {
		t.Fatalf("error getting the threads: %s", err)
	}
	thread := &Thread{}
	if !threads.Next(thread) {
		t.Fatalf("threads.Next(thread): unable to fetch the first and only thread")
	}
	if want, got := 3, thread.Count(); want != got {
		t.Errorf("db.NewQuery(%q).Threads()[0].Count(): want %d got %d", qs, want, got)
	}
	if want, got := 1, thread.CountMatched(); want != got {
		t.Errorf("db.NewQuery(%q).Threads()[0].CountMatched(): want %d got %d", qs, want, got)
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

func TestTopLevelMessages(t *testing.T) {
	db, err := Open(dbPath, DBReadOnly)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	qs := "subject:\"Introducing myself\""
	threads, err := db.NewQuery(qs).Threads()
	if err != nil {
		t.Fatalf("error getting the threads: %s", err)
	}
	thread := &Thread{}
	if !threads.Next(thread) {
		t.Fatalf("threads.Next(thread): unable to fetch the first and only thread")
	}
	msgs := thread.TopLevelMessages()
	if want, got := thread.cptr, msgs.thread.cptr; want != got {
		t.Errorf("thread.Message().thread: want %s got %s", want, got)
	}
	message := &Message{}
	var count int
	for msgs.Next(message) {
		count++
	}
	if want, got := 1, count; want != got {
		t.Errorf("db.NewQuery(%q).Threads()[0].TopLevelMessages(): want %d got %d", qs, want, got)
	}
}

func TestMessages(t *testing.T) {
	db, err := Open(dbPath, DBReadOnly)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	qs := "subject:\"Introducing myself\""
	threads, err := db.NewQuery(qs).Threads()
	if err != nil {
		t.Fatalf("error getting the threads: %s", err)
	}
	thread := &Thread{}
	if !threads.Next(thread) {
		t.Fatalf("threads.Next(thread): unable to fetch the first and only thread")
	}
	msgs := thread.Messages()
	if want, got := thread.cptr, msgs.thread.cptr; want != got {
		t.Errorf("thread.Message().thread: want %s got %s", want, got)
	}
	message := &Message{}
	var count int
	for msgs.Next(message) {
		count++
	}
	if want, got := 3, count; want != got {
		t.Errorf("db.NewQuery(%q).Threads()[0].Messages(): want %d got %d", qs, want, got)
	}
}

func TestAuthors(t *testing.T) {
	tests := map[string][]struct {
		matched   []string
		unmatched []string
	}{
		"subject:\"Introducing myself\"": {
			0: {
				unmatched: []string{"Adrian Perez de Castro", "Keith Packard", "Carl Worth"},
			},
		},

		"from:Jan": {
			0: {
				unmatched: []string{"Jan Janak"},
			},

			1: {
				matched:   []string{"Jan Janak"},
				unmatched: []string{"Carl Worth"},
			},

			2: {
				matched:   []string{"Jan Janak"},
				unmatched: []string{"Carl Worth"},
			},
		},
	}

	db, err := Open(dbPath, DBReadOnly)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	for q, ress := range tests {
		threads, err := db.NewQuery(q).Threads()
		if err != nil {
			t.Fatalf("error getting the threads: %s", err)
		}
		thread := &Thread{}
		for i := 0; threads.Next(thread); i++ {
			matched, unmatched := thread.Authors()
			if want, got := ress[i].matched, matched; !reflect.DeepEqual(want, got) {
				t.Errorf("thread.Authors() matched: want %v got %v", want, got)
			}
			if want, got := ress[i].unmatched, unmatched; !reflect.DeepEqual(want, got) {
				t.Errorf("thread.Authors() unmatched: want %v got %v", want, got)
			}
		}
	}
}

func TestOldestDate(t *testing.T) {
	db, err := Open(dbPath, DBReadOnly)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	qs := "subject:\"Introducing myself\""
	threads, err := db.NewQuery(qs).Threads()
	if err != nil {
		t.Fatalf("error getting the threads: %s", err)
	}
	thread := &Thread{}
	if !threads.Next(thread) {
		t.Fatalf("threads.Next(thread): unable to fetch the first and only thread")
	}
	if want, got := time.Unix(1258500059, 0), thread.OldestDate(); want.Unix() != got.Unix() {
		t.Errorf("thread.OldestDate(): want %s got %s", want, got)
	}
}

func TestNewestDate(t *testing.T) {
	db, err := Open(dbPath, DBReadOnly)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	qs := "subject:\"Introducing myself\""
	threads, err := db.NewQuery(qs).Threads()
	if err != nil {
		t.Fatalf("error getting the threads: %s", err)
	}
	thread := &Thread{}
	if !threads.Next(thread) {
		t.Fatalf("threads.Next(thread): unable to fetch the first and only thread")
	}
	if want, got := time.Unix(1258542931, 0), thread.NewestDate(); want.Unix() != got.Unix() {
		t.Errorf("thread.NewestDate(): want %s got %s", want, got)
	}
}

func TestThreadTags(t *testing.T) {
	db, err := Open(dbPath, DBReadOnly)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	qs := "subject:\"Introducing myself\""
	threads, err := db.NewQuery(qs).Threads()
	if err != nil {
		t.Fatalf("error getting the threads: %s", err)
	}
	thread := &Thread{}
	if !threads.Next(thread) {
		t.Fatalf("threads.Next(thread): unable to fetch the first and only thread")
	}
	ts := thread.Tags()
	tag := &Tag{}
	var tags []string
	for ts.Next(tag) {
		tags = append(tags, tag.Value)
	}
	if want, got := []string{"inbox", "signed", "unread"}, tags; !reflect.DeepEqual(want, got) {
		t.Errorf("thread.Tags(): want %v got %v", want, got)
	}
}
