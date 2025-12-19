package main

import (
	"flag"
	"fmt"
	"os"

	"lazy-todo/internal/storage"
	"lazy-todo/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

var version = "0.2.0"

func main() {
	// Command line flags
	filePath := flag.String("file", "", "Chemin vers le fichier de tâches (défaut: ~/.local/share/lazy-todo/tasks.yaml)")
	showVersion := flag.Bool("version", false, "Afficher la version")
	flag.Parse()

	if *showVersion {
		fmt.Printf("lazy-todo v%s\n", version)
		os.Exit(0)
	}

	// Determine file path
	path := *filePath
	if path == "" {
		path = storage.DefaultFilePath()
	}

	// Create storage
	store := storage.NewStorage(path)

	// Create and run the app
	app := ui.NewApp(store)

	p := tea.NewProgram(app, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Erreur: %v\n", err)
		os.Exit(1)
	}
}
