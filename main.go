package main

import (
	"embed"
	"log"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v3/pkg/application"

	"allinone/internal/db"
	"allinone/internal/library"
	"allinone/internal/settings"
	"allinone/internal/steamgriddb"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	dbPath, err := databasePath()
	if err != nil {
		log.Fatalf("resolve database path: %v", err)
	}
	gdb, err := db.Open(dbPath)
	if err != nil {
		log.Fatalf("open database at %s: %v", dbPath, err)
	}
	defer func() {
		if err := db.Close(gdb); err != nil {
			log.Printf("close database: %v", err)
		}
	}()

	app := application.New(application.Options{
		Name:        "AllInOne",
		Description: "Game library launcher",
		Services: []application.Service{
			application.NewService(library.NewService(library.NewStore(gdb))),
			application.NewService(settings.NewService(settings.NewStore(gdb), steamgriddb.New())),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            "AllInOne",
		Width:            1200,
		Height:           760,
		BackgroundColour: application.NewRGB(6, 7, 15),
		URL:              "/",
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

// databasePath returns %APPDATA%\AllInOne\allinone.db on Windows (and the
// equivalent user config dir elsewhere).
func databasePath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "AllInOne", "allinone.db"), nil
}
