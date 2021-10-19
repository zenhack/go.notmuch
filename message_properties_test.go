package notmuch

import (
	"testing"
)

func TestMessagesProperties(t *testing.T) {
	db, err := Open(dbPath, DBReadWrite)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	qs := "subject:\"Introducing myself\""
	messages, err := db.NewQuery(qs).Messages()
	if err != nil {
		t.Fatalf("error getting the messages: %s", err)
	}

	first := &Message{}
	found := messages.Next(&first)
	if !found {
		t.Fatalf("couldn't get the first message: %s", err)
	}

	if err := first.AddProperty("go-notmuch-test", "success"); err != nil {
		t.Fatalf("couldn't add property: %s", err)
	}

	properties := first.Properties("go-notmuch-test", true)
	property := &MessageProperty{}
	for properties.Next(&property) {
		if property.Key == "go-notmuch-test" && property.Value == "success" {
			return
		}
	}

	t.Fatalf("couldn't find expected property")
}
