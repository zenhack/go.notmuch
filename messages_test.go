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
	// FIXME: this leaks the query object. Ongoing discussion on issue #11 as to
	// whether we should bring the GC back or just require doing this in two
	// statements. This appears in several other places in the tests.
	threads, err := db.NewQuery(qs).Threads()
	if err != nil {
		t.Fatalf("error getting the threads: %s", err)
	}
	// This takes care of all of the thread's child objects, too.
	defer threads.Close()
	thread := &Thread{}
	if !threads.Next(thread) {
		t.Fatalf("threads.Next(thread): unable to fetch the first and only thread")
	}
	msgs := thread.Messages()
	ts := msgs.Tags()
	tag := &Tag{}
	var tags []string
	for ts.Next(tag) {
		tags = append(tags, tag.Value)
	}
	if want, got := []string{"inbox", "signed", "unread"}, tags; !reflect.DeepEqual(want, got) {
		t.Errorf("thread.Tags(): want %v got %v", want, got)
	}
}
