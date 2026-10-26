// Package settings persists app settings and the folders scanned for games.
package settings

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"allinone/internal/models"
)

const KeySteamGridDB = "steamgriddb_api_key"

var (
	ErrFolderExists   = errors.New("folder already added")
	ErrFolderNotFound = errors.New("folder not found")
)

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{db: db}
}

// Get returns the value for key, or "" if it was never set.
func (s *Store) Get(key string) (string, error) {
	var st models.Setting
	err := s.db.Where("key = ?", key).Limit(1).Find(&st).Error
	if err != nil {
		return "", fmt.Errorf("get setting %q: %w", key, err)
	}
	return st.Value, nil
}

func (s *Store) Set(key, value string) error {
	err := s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).Create(&models.Setting{Key: key, Value: value}).Error
	if err != nil {
		return fmt.Errorf("set setting %q: %w", key, err)
	}
	return nil
}

func (s *Store) ListFolders() ([]models.ScanFolder, error) {
	folders := []models.ScanFolder{}
	if err := s.db.Order("path COLLATE NOCASE").Find(&folders).Error; err != nil {
		return nil, fmt.Errorf("list folders: %w", err)
	}
	return folders, nil
}

// AddFolder stores an existing absolute directory. Paths are compared
// case-insensitively because Windows paths are.
func (s *Store) AddFolder(path string) (models.ScanFolder, error) {
	path = filepath.Clean(path)
	if !filepath.IsAbs(path) {
		return models.ScanFolder{}, fmt.Errorf("folder %q must be an absolute path", path)
	}
	info, err := os.Stat(path)
	if err != nil {
		return models.ScanFolder{}, fmt.Errorf("folder %q: %w", path, err)
	}
	if !info.IsDir() {
		return models.ScanFolder{}, fmt.Errorf("%q is not a folder", path)
	}

	var count int64
	if err := s.db.Model(&models.ScanFolder{}).Where("LOWER(path) = LOWER(?)", path).Count(&count).Error; err != nil {
		return models.ScanFolder{}, fmt.Errorf("check folder %q: %w", path, err)
	}
	if count > 0 {
		return models.ScanFolder{}, fmt.Errorf("%q: %w", path, ErrFolderExists)
	}

	f := models.ScanFolder{Path: path}
	if err := s.db.Create(&f).Error; err != nil {
		return models.ScanFolder{}, fmt.Errorf("add folder %q: %w", path, err)
	}
	return f, nil
}

func (s *Store) RemoveFolder(id uint) error {
	res := s.db.Delete(&models.ScanFolder{}, id)
	if res.Error != nil {
		return fmt.Errorf("remove folder %d: %w", id, res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("remove folder %d: %w", id, ErrFolderNotFound)
	}
	return nil
}
