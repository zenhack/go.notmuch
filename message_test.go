package notmuch

import "testing"

func TestMessageID(t *testing.T) {
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
	msg := &Message{}
	if !msgs.Next(msg) {
		t.Fatalf("msgs.Next(msg): unable to fetch the first message in the thread")
	}
	if want, got := "20091118002059.067214ed@hikari", msg.ID(); want != got {
		t.Errorf("msg.ID(): want %s got %s", want, got)
	}
}

func TestMessageThreadID(t *testing.T) {
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
	msg := &Message{}
	if !msgs.Next(msg) {
		t.Fatalf("msgs.Next(msg): unable to fetch the first message in the thread")
	}
	if want, got := thread.ID(), msg.ThreadID(); want != got {
		t.Errorf("msg.ThreadID(): want %s got %s", want, got)
	}
}
