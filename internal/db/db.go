// Package db opens the local SQLite database through GORM using a pure-Go
// driver (modernc.org/sqlite), so the app builds with CGO_ENABLED=0.
package db

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"allinone/internal/models"
)

const pragmas = "?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)"

// Open creates the parent directory if needed, opens the database at path and
// migrates the schema. It is safe to call on an existing database.
func Open(path string) (*gorm.DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("create database directory for %q: %w", path, err)
	}

	gdb, err := gorm.Open(sqlite.Open(path+pragmas), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("open database %q: %w", path, err)
	}

	if err := gdb.AutoMigrate(models.All()...); err != nil {
		return nil, fmt.Errorf("migrate database %q: %w", path, err)
	}
	return gdb, nil
}

// Close releases the underlying connection pool.
func Close(gdb *gorm.DB) error {
	sqlDB, err := gdb.DB()
	if err != nil {
		return fmt.Errorf("get sql.DB: %w", err)
	}
	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("close database: %w", err)
	}
	return nil
}
