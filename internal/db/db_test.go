package db

import (
	"os"
	"path/filepath"
	"testing"

	"gorm.io/gorm"

	"allinone/internal/models"
)

// The app ships as a single CGO-free .exe; this proves the pure-Go driver can
// create a fresh database and the full schema from nothing.
func TestOpenCreatesDatabaseAndSchema(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "allinone.db")

	gdb, err := Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer mustClose(t, gdb)

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("database file not created: %v", err)
	}
	for _, m := range models.All() {
		if !gdb.Migrator().HasTable(m) {
			t.Errorf("table for %T was not created", m)
		}
	}
}

// The app opens the same database on every launch; reopening must not fail
// or lose data.
func TestOpenIsIdempotentAndKeepsData(t *testing.T) {
	path := filepath.Join(t.TempDir(), "allinone.db")

	first, err := Open(path)
	if err != nil {
		t.Fatalf("first Open: %v", err)
	}
	if err := first.Create(&models.Setting{Key: "k", Value: "v"}).Error; err != nil {
		t.Fatalf("insert setting: %v", err)
	}
	mustClose(t, first)

	second, err := Open(path)
	if err != nil {
		t.Fatalf("second Open: %v", err)
	}
	defer mustClose(t, second)

	var s models.Setting
	if err := second.First(&s, "key = ?", "k").Error; err != nil {
		t.Fatalf("setting lost after reopen: %v", err)
	}
}

// A bad path must surface an error instead of silently running without storage.
func TestOpenFailsLoudOnUnusablePath(t *testing.T) {
	blocker := filepath.Join(t.TempDir(), "file")
	if err := os.WriteFile(blocker, nil, 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := Open(filepath.Join(blocker, "allinone.db")); err == nil {
		t.Fatal("expected error when parent path is a file, got nil")
	}
}

func mustClose(t *testing.T, gdb *gorm.DB) {
	t.Helper()
	if err := Close(gdb); err != nil {
		t.Fatalf("Close: %v", err)
	}
}
