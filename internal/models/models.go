// Package models holds the GORM models persisted in the local SQLite database.
package models

import "time"

// Source identifies how a game entered the library.
type Source string

const (
	SourceFolder    Source = "folder"
	SourceRegistry  Source = "registry"
	SourceStartMenu Source = "startmenu"
	SourceManual    Source = "manual"
)

type Game struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Name      string `gorm:"not null;index" json:"name"`
	ExePath   string `gorm:"not null;uniqueIndex" json:"exePath"`
	WorkDir   string `json:"workDir"`
	Args      string `json:"args"`
	Source    Source `gorm:"not null" json:"source"`
	CoverPath string `json:"coverPath"`
	IconPath  string `json:"iconPath"`
	Favorite  bool   `gorm:"not null;default:false" json:"favorite"`
	Hidden    bool   `gorm:"not null;default:false" json:"hidden"`
	// Tags is stored as JSON so tag support (post-v1) needs no schema migration.
	Tags      []string  `gorm:"serializer:json" json:"tags"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ScanFolder is a user-configured folder scanned for game executables.
type ScanFolder struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Path      string    `gorm:"not null;uniqueIndex" json:"path"`
	CreatedAt time.Time `json:"createdAt"`
}

// Setting is a key/value app setting (e.g. the SteamGridDB API key).
type Setting struct {
	Key   string `gorm:"primaryKey"`
	Value string `gorm:"not null"`
}

// All lists every model, in migration order.
func All() []any {
	return []any{&Game{}, &ScanFolder{}, &Setting{}}
}
