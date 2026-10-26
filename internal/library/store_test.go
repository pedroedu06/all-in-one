package library

import (
	"errors"
	"path/filepath"
	"slices"
	"testing"

	"allinone/internal/db"
	"allinone/internal/models"
)

func newStore(t *testing.T) *Store {
	t.Helper()
	gdb, err := db.Open(filepath.Join(t.TempDir(), "test.db"))
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

func addGame(t *testing.T, s *Store, name string, mods ...func(*models.Game)) *models.Game {
	t.Helper()
	g := &models.Game{Name: name, ExePath: `C:\Games\` + name + `\game.exe`, Source: models.SourceManual}
	for _, m := range mods {
		m(g)
	}
	if err := s.Add(g); err != nil {
		t.Fatalf("add %q: %v", name, err)
	}
	return g
}

func names(games []models.Game) []string {
	out := make([]string, len(games))
	for i, g := range games {
		out[i] = g.Name
	}
	return out
}

// Tags are in the schema now so post-v1 tagging needs no migration; they must
// survive a round-trip through SQLite unchanged.
func TestTagsRoundTrip(t *testing.T) {
	s := newStore(t)
	addGame(t, s, "Hades", func(g *models.Game) { g.Tags = []string{"roguelike", "indie"} })

	games, err := s.List(ListOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if got := games[0].Tags; !slices.Equal(got, []string{"roguelike", "indie"}) {
		t.Fatalf("tags = %v, want [roguelike indie]", got)
	}
}

// Discovery may find the same executable from several sources (folder scan,
// registry, Start Menu); the library must never hold it twice.
func TestAddRejectsDuplicateExePath(t *testing.T) {
	s := newStore(t)
	first := addGame(t, s, "Celeste")

	dup := &models.Game{Name: "Celeste (Start Menu)", ExePath: first.ExePath, Source: models.SourceStartMenu}
	if err := s.Add(dup); err == nil {
		t.Fatal("expected error adding duplicate exe path, got nil")
	}
}

// Hiding a game removes it from the library view without deleting it.
func TestListExcludesHiddenUnlessRequested(t *testing.T) {
	s := newStore(t)
	addGame(t, s, "Visible")
	hidden := addGame(t, s, "Hidden")
	if err := s.SetHidden(hidden.ID, true); err != nil {
		t.Fatal(err)
	}

	def, err := s.List(ListOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if got := names(def); !slices.Equal(got, []string{"Visible"}) {
		t.Fatalf("default list = %v, want [Visible]", got)
	}

	all, err := s.List(ListOptions{IncludeHidden: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 2 {
		t.Fatalf("IncludeHidden list has %d games, want 2", len(all))
	}
}

// Users type partial names in any case in the search box.
func TestListSearchIsPartialAndCaseInsensitive(t *testing.T) {
	s := newStore(t)
	addGame(t, s, "Hollow Knight")
	addGame(t, s, "Stardew Valley")

	got, err := s.List(ListOptions{Query: "KNIG"})
	if err != nil {
		t.Fatal(err)
	}
	if n := names(got); !slices.Equal(n, []string{"Hollow Knight"}) {
		t.Fatalf("search = %v, want [Hollow Knight]", n)
	}
}

// Favorites are pinned to the top; the rest is alphabetical regardless of case.
func TestListOrdersFavoritesFirstThenByName(t *testing.T) {
	s := newStore(t)
	addGame(t, s, "zelda-like")
	addGame(t, s, "Apex")
	fav := addGame(t, s, "Terraria")
	if err := s.SetFavorite(fav.ID, true); err != nil {
		t.Fatal(err)
	}

	got, err := s.List(ListOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if n := names(got); !slices.Equal(n, []string{"Terraria", "Apex", "zelda-like"}) {
		t.Fatalf("order = %v, want [Terraria Apex zelda-like]", n)
	}
}

// Toggling a flag on a game that no longer exists must fail loudly so the UI
// can tell the user, instead of pretending it worked.
func TestSetFlagOnMissingGameReturnsNotFound(t *testing.T) {
	s := newStore(t)

	if err := s.SetFavorite(999, true); !errors.Is(err, ErrGameNotFound) {
		t.Fatalf("SetFavorite err = %v, want ErrGameNotFound", err)
	}
	if err := s.SetHidden(999, true); !errors.Is(err, ErrGameNotFound) {
		t.Fatalf("SetHidden err = %v, want ErrGameNotFound", err)
	}
}
