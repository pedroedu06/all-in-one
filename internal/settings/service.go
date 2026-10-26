package settings

import (
	"context"
	"errors"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/application"

	"allinone/internal/models"
	"allinone/internal/steamgriddb"
)

// Service is the Wails-bound facade for the settings screen. Errors are
// returned in Portuguese because the frontend shows them as-is.
type Service struct {
	store *Store
	sgdb  *steamgriddb.Client
}

func NewService(store *Store, sgdb *steamgriddb.Client) *Service {
	return &Service{store: store, sgdb: sgdb}
}

func (s *Service) GetSteamGridDBKey() (string, error) {
	return s.store.Get(KeySteamGridDB)
}

// SetSteamGridDBKey saves the key; an empty key clears it.
func (s *Service) SetSteamGridDBKey(key string) error {
	return s.store.Set(KeySteamGridDB, strings.TrimSpace(key))
}

// TestSteamGridDBKey returns nil when SteamGridDB accepts the key.
func (s *Service) TestSteamGridDBKey(ctx context.Context, key string) error {
	err := s.sgdb.ValidateKey(ctx, key)
	switch {
	case err == nil:
		return nil
	case errors.Is(err, steamgriddb.ErrInvalidKey):
		return errors.New("Chave inválida. Confira em steamgriddb.com/profile/preferences/api")
	case errors.Is(err, steamgriddb.ErrUnreachable):
		return errors.New("Não foi possível conectar ao SteamGridDB. Verifique sua internet")
	default:
		return errors.New("Erro ao testar a chave: " + err.Error())
	}
}

func (s *Service) ListFolders() ([]models.ScanFolder, error) {
	return s.store.ListFolders()
}

func (s *Service) AddFolder(path string) (models.ScanFolder, error) {
	f, err := s.store.AddFolder(path)
	if errors.Is(err, ErrFolderExists) {
		return f, errors.New("Essa pasta já foi adicionada")
	}
	return f, err
}

func (s *Service) RemoveFolder(id uint) error {
	return s.store.RemoveFolder(id)
}

// PickFolder opens the native folder picker and returns "" if cancelled.
func (s *Service) PickFolder() (string, error) {
	return application.Get().Dialog.OpenFile().
		CanChooseDirectories(true).
		CanChooseFiles(false).
		SetTitle("Escolher pasta de jogos").
		PromptForSingleSelection()
}
