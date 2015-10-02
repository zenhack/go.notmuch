package notmuch

import (
	"path"
	"reflect"
	"runtime"
	"testing"
	"time"
)

func TestMessageID(t *testing.T) {
	db, err := Open(dbPath, DBReadOnly)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	thread, err := firstThread(db, "subject:\"Introducing myself\"")
	if err != nil {
		t.Fatal(err)
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

	thread, err := firstThread(db, "subject:\"Introducing myself\"")
	if err != nil {
		t.Fatal(err)
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

	thread, err := firstThread(db, "subject:\"Introducing myself\"")
	if err != nil {
		t.Fatal(err)
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

		// invoke the GC to make sure it's running smoothly.
		if count%2 == 0 {
			runtime.GC()
		}
	}
	if want, got := 2, count; want != got {
		t.Errorf("msg.Replies(): want %d replies got %d", want, got)
	}

	// msg is now the last message and it shouldn't have any replies
	if _, err := msg.Replies(); err != ErrNoRepliesOrPointerNotFromThread {
		t.Errorf("msg.Replies() on the last message: expecting error %q got %q", ErrNoRepliesOrPointerNotFromThread, err)
	}
}

func TestMessageFilename(t *testing.T) {
	db, err := Open(dbPath, DBReadOnly)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	thread, err := firstThread(db, "subject:\"Introducing myself\"")
	if err != nil {
		t.Fatal(err)
	}
	msgs := thread.Messages()
	msg := &Message{}
	if !msgs.Next(msg) {
		t.Fatalf("msgs.Next(msg): unable to fetch the first message in the thread")
	}

	var fn string
	fns := msg.Filenames()
	if !fns.Next(&fn) {
		t.Fatalf("msg.Filename: unable to fetch a filename but it's known to have 2")
	}
	if want, got := path.Join(dbPath, "bar/cur/20:2,"), fn; want != got {
		t.Errorf("msg.Filename(): want %s got %s", want, got)
	}
}

func TestMessageFilenames(t *testing.T) {
	db, err := Open(dbPath, DBReadOnly)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	thread, err := firstThread(db, "id:20091117232137.GA7669@griffis1.net")
	if err != nil {
		t.Fatal(err)
	}
	msgs := thread.Messages()
	msg := &Message{}
	if !msgs.Next(msg) {
		t.Fatalf("msgs.Next(msg): unable to fetch the first message in the thread")
	}

	var count int
	var fn string
	fns := msg.Filenames()
	for fns.Next(&fn) {
		count++

		// invoke the GC to make sure it's running smoothly.
		if count%2 == 0 {
			runtime.GC()
		}
	}

	if want, got := 2, count; want != got {
		t.Errorf("msg.Filenames(): want %d filename got %d", want, got)
	}
}

func TestMessageDate(t *testing.T) {
	db, err := Open(dbPath, DBReadOnly)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	thread, err := firstThread(db, "id:20091117232137.GA7669@griffis1.net")
	if err != nil {
		t.Fatal(err)
	}
	msgs := thread.Messages()
	msg := &Message{}
	if !msgs.Next(msg) {
		t.Fatalf("msgs.Next(msg): unable to fetch the first message in the thread")
	}
	if want, got := time.Unix(1258500098, 0), msg.Date(); want.Unix() != got.Unix() {
		t.Errorf("msg.Date(): want %s got %s", want, got)
	}
}

func TestMessageHeader(t *testing.T) {
	db, err := Open(dbPath, DBReadOnly)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	thread, err := firstThread(db, "id:1258471718-6781-2-git-send-email-dottedmag@dottedmag.net")
	if err != nil {
		t.Fatal(err)
	}
	msgs := thread.Messages()
	msg := &Message{}
	for msgs.Next(msg) {
		if msg.ID() == "1258471718-6781-2-git-send-email-dottedmag@dottedmag.net" {
			break
		}

		// invoke the GC to make sure it's running smoothly.
		runtime.GC()
	}
	for _, hn := range []string{"References", "references"} {
		if want, got := "<1258471718-6781-1-git-send-email-dottedmag@dottedmag.net>", msg.Header(hn); want != got {
			t.Errorf("msg.Header(%q): want %s got %s", hn, want, got)
		}
	}
}

func TestMessageTags(t *testing.T) {
	db, err := Open(dbPath, DBReadOnly)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	thread, err := firstThread(db, "id:1258471718-6781-2-git-send-email-dottedmag@dottedmag.net")
	if err != nil {
		t.Fatal(err)
	}
	msgs := thread.Messages()
	msg := &Message{}
	for msgs.Next(msg) {
		if msg.ID() == "1258471718-6781-2-git-send-email-dottedmag@dottedmag.net" {
			break
		}
		// invoke the GC to make sure it's running smoothly.
		runtime.GC()
	}

	ts := msg.Tags()
	tag := &Tag{}
	var tags []string
	for ts.Next(tag) {
		tags = append(tags, tag.Value)
		// invoke the GC to make sure it's running smoothly.
		runtime.GC()
	}
	if want, got := []string{"inbox", "unread"}, tags; !reflect.DeepEqual(want, got) {
		t.Errorf("msg.Tags(): want %v got %v", want, got)
	}
}

func TestMessageAddRemoveTagReadonlyDB(t *testing.T) {
	db, err := Open(dbPath, DBReadOnly)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	thread, err := firstThread(db, "id:1258471718-6781-2-git-send-email-dottedmag@dottedmag.net")
	if err != nil {
		t.Fatal(err)
	}
	msgs := thread.Messages()
	msg := &Message{}
	for msgs.Next(msg) {
		if msg.ID() == "1258471718-6781-2-git-send-email-dottedmag@dottedmag.net" {
			break
		}
		// invoke the GC to make sure it's running smoothly.
		runtime.GC()
	}

	ts := msg.Tags()
	tag := &Tag{}
	var tags []string
	for ts.Next(tag) {
		tags = append(tags, tag.Value)
		// invoke the GC to make sure it's running smoothly.
		runtime.GC()
	}
	if want, got := []string{"inbox", "unread"}, tags; !reflect.DeepEqual(want, got) {
		t.Errorf("msg.Tags(): want %v got %v", want, got)
	}
	tn := "newtag"
	if err := msg.AddTag(tn); err != ErrReadOnlyDB {
		t.Errorf("msg.AddTag(%q): want error %s got %s", tn, ErrReadOnlyDB, err)
	}
	if want, got := []string{"inbox", "unread"}, tags; !reflect.DeepEqual(want, got) {
		t.Errorf("msg.Tags(): want %v got %v", want, got)
	}
	if err := msg.RemoveTag(tn); err != ErrReadOnlyDB {
		t.Errorf("msg.RemoveTag(%q): want error %s got %s", tn, ErrReadOnlyDB, err)
	}
	if want, got := []string{"inbox", "unread"}, tags; !reflect.DeepEqual(want, got) {
		t.Errorf("msg.Tags(): want %v got %v", want, got)
	}
}

func TestMessageAddRemoveTag(t *testing.T) {
	db, err := Open(dbPath, DBReadWrite)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	thread, err := firstThread(db, "id:1258471718-6781-2-git-send-email-dottedmag@dottedmag.net")
	if err != nil {
		t.Fatal(err)
	}
	msgs := thread.Messages()
	msg := &Message{}
	for msgs.Next(msg) {
		if msg.ID() == "1258471718-6781-2-git-send-email-dottedmag@dottedmag.net" {
			break
		}
		// invoke the GC to make sure it's running smoothly.
		runtime.GC()
	}

	ts := msg.Tags()
	tag := &Tag{}
	var tags []string
	for ts.Next(tag) {
		tags = append(tags, tag.Value)
		// invoke the GC to make sure it's running smoothly.
		runtime.GC()
	}
	if want, got := []string{"inbox", "unread"}, tags; !reflect.DeepEqual(want, got) {
		t.Errorf("msg.Tags(): want %v got %v", want, got)
	}

	tn := "newtag"
	if err := msg.AddTag(tn); err != nil {
		t.Fatalf("msg.AddTag(%q): got error: %s", tn, err)
	}
	ts = msg.Tags()
	tags = []string{}
	for ts.Next(tag) {
		tags = append(tags, tag.Value)
		// invoke the GC to make sure it's running smoothly.
		runtime.GC()
	}
	if want, got := []string{"inbox", tn, "unread"}, tags; !reflect.DeepEqual(want, got) {
		t.Errorf("msg.Tags(): want %v got %v", want, got)
	}

	if err := msg.RemoveTag(tn); err != nil {
		t.Fatalf("msg.RemoveTag(%q): got error: %s", tn, err)
	}
	ts = msg.Tags()
	tags = []string{}
	for ts.Next(tag) {
		tags = append(tags, tag.Value)
		// invoke the GC to make sure it's running smoothly.
		runtime.GC()
	}
	if want, got := []string{"inbox", "unread"}, tags; !reflect.DeepEqual(want, got) {
		t.Errorf("msg.Tags(): want %v got %v", want, got)
	}

	if err := msg.RemoveAllTags(); err != nil {
		t.Fatalf("msg.RemoveAllTag(): got error: %s", err)
	}
	ts = msg.Tags()
	tags = []string{}
	for ts.Next(tag) {
		tags = append(tags, tag.Value)
		// invoke the GC to make sure it's running smoothly.
		runtime.GC()
	}
	if want, got := []string{}, tags; !reflect.DeepEqual(want, got) {
		t.Errorf("msg.Tags(): want %v got %v", want, got)
	}

	// return the DB to the pristine condition.
	msg.AddTag("inbox")
	msg.AddTag("unread")
}

func TestMessageAtomic(t *testing.T) {
	db, err := Open(dbPath, DBReadWrite)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	thread, err := firstThread(db, "id:1258471718-6781-2-git-send-email-dottedmag@dottedmag.net")
	if err != nil {
		t.Fatal(err)
	}
	msgs := thread.Messages()
	msg := &Message{}
	for msgs.Next(msg) {
		if msg.ID() == "1258471718-6781-2-git-send-email-dottedmag@dottedmag.net" {
			break
		}
		// invoke the GC to make sure it's running smoothly.
		runtime.GC()
	}

	ts := msg.Tags()
	tag := &Tag{}
	var tags []string
	for ts.Next(tag) {
		tags = append(tags, tag.Value)
		// invoke the GC to make sure it's running smoothly.
		runtime.GC()
	}
	if want, got := []string{"inbox", "unread"}, tags; !reflect.DeepEqual(want, got) {
		t.Errorf("msg.Tags(): want %v got %v", want, got)
	}

	tn := "newtag"
	msg.Atomic(func(mymsg *Message) {
		if err := mymsg.AddTag(tn); err != nil {
			t.Fatalf("msg.AddTag(%q): got error: %s", tn, err)
		}
	})
	ts = msg.Tags()
	tags = []string{}
	for ts.Next(tag) {
		tags = append(tags, tag.Value)
		// invoke the GC to make sure it's running smoothly.
		runtime.GC()
	}
	if want, got := []string{"inbox", tn, "unread"}, tags; !reflect.DeepEqual(want, got) {
		t.Errorf("msg.Tags(): want %v got %v", want, got)
	}

	msg.Atomic(func(mymsg *Message) {
		if err := mymsg.RemoveTag(tn); err != nil {
			t.Fatalf("msg.RemoveTag(%q): got error: %s", tn, err)
		}
	})
	ts = msg.Tags()
	tags = []string{}
	for ts.Next(tag) {
		tags = append(tags, tag.Value)
		// invoke the GC to make sure it's running smoothly.
		runtime.GC()
	}
	if want, got := []string{"inbox", "unread"}, tags; !reflect.DeepEqual(want, got) {
		t.Errorf("msg.Tags(): want %v got %v", want, got)
	}
}
