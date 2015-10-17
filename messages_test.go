package notmuch

import (
	"reflect"
	"testing"
)

func TestMessagesTags(t *testing.T) {
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
	tags := msgs.Tags().slice()
	if want, got := []string{"inbox", "signed", "unread"}, tags; !reflect.DeepEqual(want, got) {
		t.Errorf("thread.Tags(): want %v got %v", want, got)
	}
}
