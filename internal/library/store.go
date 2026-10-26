// Package library manages the games persisted in the local database.
package library

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"allinone/internal/models"
)

var ErrGameNotFound = errors.New("game not found")

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{db: db}
}

type ListOptions struct {
	// Query filters by a case-insensitive substring of the game name.
	Query         string
	IncludeHidden bool
}

func (s *Store) Add(g *models.Game) error {
	if err := s.db.Create(g).Error; err != nil {
		return fmt.Errorf("add game %q (%s): %w", g.Name, g.ExePath, err)
	}
	return nil
}

// List returns games with favorites first, then alphabetically by name.
func (s *Store) List(opts ListOptions) ([]models.Game, error) {
	q := s.db.Model(&models.Game{})
	if !opts.IncludeHidden {
		q = q.Where("hidden = ?", false)
	}
	if opts.Query != "" {
		// SQLite LIKE is case-insensitive for ASCII.
		q = q.Where("name LIKE ?", "%"+opts.Query+"%")
	}

	games := []models.Game{}
	if err := q.Order("favorite DESC").Order("name COLLATE NOCASE ASC").Find(&games).Error; err != nil {
		return nil, fmt.Errorf("list games: %w", err)
	}
	return games, nil
}

func (s *Store) SetFavorite(id uint, favorite bool) error {
	return s.setFlag(id, "favorite", favorite)
}

func (s *Store) SetHidden(id uint, hidden bool) error {
	return s.setFlag(id, "hidden", hidden)
}

func (s *Store) setFlag(id uint, column string, value bool) error {
	res := s.db.Model(&models.Game{}).Where("id = ?", id).Update(column, value)
	if res.Error != nil {
		return fmt.Errorf("set %s=%t on game %d: %w", column, value, id, res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("set %s on game %d: %w", column, id, ErrGameNotFound)
	}
	return nil
}
