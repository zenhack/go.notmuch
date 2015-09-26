package notmuch

import "testing"

func TestMessageGetThreadID(t *testing.T) {
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
	if want, got := thread.GetID(), msg.GetThreadID(); want != got {
		t.Errorf("msg.GetThreadID(): want %s got %s", want, got)
	}
}
