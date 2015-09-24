package notmuch

import (
	"os"
	"path"
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

func TestOpenNotFound(t *testing.T) {
	_, err := Open("/not-found", DBReadOnly)
	if err == nil {
		t.Errorf("Open(%q): expected error got nil", "/not-found")
	}
}

func TestCreate(t *testing.T) {
	tmp := os.TempDir()
	tmpDbPath := path.Join(tmp, ".notmuch")
	defer func() {
		os.RemoveAll(tmpDbPath)
	}()

	db, err := Create(tmp)
	if err != nil {
		t.Fatal(err)
	}
	if want, got := 1, db.Version(); got < want {
		t.Errorf("db.Version(): want at least %d got %d", want, got)
	}
}

func TestOpen(t *testing.T) {
	db, err := Open(dbPath, DBReadOnly)
	if err != nil {
		t.Fatalf("Open(%q): unexpected error: %s", dbPath, err)
	}
	defer db.Close()
	if want, got := 1, db.Version(); got < want {
		t.Errorf("db.Version(): want at least %d got %d", want, got)
	}
}
