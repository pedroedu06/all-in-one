package settings

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"allinone/internal/db"
)

func openStore(t *testing.T, path string) *Store {
	t.Helper()
	gdb, err := db.Open(path)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Close(gdb); err != nil {
			t.Errorf("close db: %v", err)
		}
	})
	return NewStore(gdb)
}

func newStore(t *testing.T) *Store {
	t.Helper()
	return openStore(t, filepath.Join(t.TempDir(), "test.db"))
}

// The app reopens the database on every launch; the saved key must still be there.
func TestSettingSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.db")
	gdb, err := db.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := NewStore(gdb).Set(KeySteamGridDB, "abc"); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(gdb); err != nil {
		t.Fatal(err)
	}

	got, err := openStore(t, path).Get(KeySteamGridDB)
	if err != nil {
		t.Fatal(err)
	}
	if got != "abc" {
		t.Fatalf("Get = %q, want abc", got)
	}
}

// Saving the key again replaces it instead of failing on the primary key.
func TestSetOverwrites(t *testing.T) {
	s := newStore(t)
	for _, v := range []string{"old", "new"} {
		if err := s.Set(KeySteamGridDB, v); err != nil {
			t.Fatalf("Set(%q): %v", v, err)
		}
	}
	if got, _ := s.Get(KeySteamGridDB); got != "new" {
		t.Fatalf("Get = %q, want new", got)
	}
}

// An unset key is a normal state (the key is optional), not an error.
func TestGetMissingReturnsEmpty(t *testing.T) {
	got, err := newStore(t).Get("missing")
	if err != nil || got != "" {
		t.Fatalf("Get = %q, %v; want empty, nil", got, err)
	}
}

// The folder scanner trusts that stored folders are real absolute directories.
func TestAddFolderRejectsInvalidPaths(t *testing.T) {
	s := newStore(t)
	dir := t.TempDir()
	file := filepath.Join(dir, "game.exe")
	if err := os.WriteFile(file, nil, 0o644); err != nil {
		t.Fatal(err)
	}

	for name, path := range map[string]string{
		"missing":  filepath.Join(dir, "nope"),
		"file":     file,
		"relative": "Games",
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := s.AddFolder(path); err == nil {
				t.Fatalf("AddFolder(%q) succeeded, want error", path)
			}
		})
	}
	if folders, _ := s.ListFolders(); len(folders) != 0 {
		t.Fatalf("invalid folders were stored: %v", folders)
	}
}

// The same folder typed with different case or a trailing separator would be
// scanned twice; Windows treats them as the same path.
func TestAddFolderRejectsSameFolderWrittenDifferently(t *testing.T) {
	s := newStore(t)
	dir := t.TempDir()
	if _, err := s.AddFolder(dir); err != nil {
		t.Fatal(err)
	}

	for _, variant := range []string{strings.ToUpper(dir), dir + string(filepath.Separator)} {
		if _, err := s.AddFolder(variant); !errors.Is(err, ErrFolderExists) {
			t.Fatalf("AddFolder(%q) err = %v, want ErrFolderExists", variant, err)
		}
	}
}

func TestRemoveFolder(t *testing.T) {
	s := newStore(t)
	f, err := s.AddFolder(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	if err := s.RemoveFolder(f.ID); err != nil {
		t.Fatalf("RemoveFolder: %v", err)
	}
	if folders, _ := s.ListFolders(); len(folders) != 0 {
		t.Fatalf("folder still listed after removal: %v", folders)
	}
	if err := s.RemoveFolder(f.ID); !errors.Is(err, ErrFolderNotFound) {
		t.Fatalf("second RemoveFolder err = %v, want ErrFolderNotFound", err)
	}
}
