package library

import "allinone/internal/models"

// Service is the Wails-bound facade the frontend calls.
type Service struct {
	store *Store
}

func NewService(store *Store) *Service {
	return &Service{store: store}
}

func (s *Service) ListGames(query string, includeHidden bool) ([]models.Game, error) {
	return s.store.List(ListOptions{Query: query, IncludeHidden: includeHidden})
}

func (s *Service) SetFavorite(id uint, favorite bool) error {
	return s.store.SetFavorite(id, favorite)
}

func (s *Service) SetHidden(id uint, hidden bool) error {
	return s.store.SetHidden(id, hidden)
}
