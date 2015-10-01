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

func TestMessageReplies(t *testing.T) {
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

	replies, err := msg.Replies()
	if err != nil {
		t.Fatalf("msg.Replies(): unexpected error: %s", err)
	}
	var count int
	for replies.Next(msg) {
		count++
	}
	if want, got := 2, count; want != got {
		t.Errorf("msg.Replies(): want %d replies got %d", want, got)
	}

	// msg is now the last message and it shouldn't have any replies
	if _, err := msg.Replies(); err != ErrNoRepliesOrPointerNotFromThread {
		t.Errorf("msg.Replies() on the last message: expecting error %q got %q", ErrNoRepliesOrPointerNotFromThread, err)
	}
}
