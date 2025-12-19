# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

lazy-todo is a TUI (Terminal User Interface) todo list application written in Go, inspired by lazygit's interface style. It uses the Bubble Tea framework for the TUI, Lip Gloss for styling, and YAML for task storage.

## Build & Run Commands

```bash
# Build the application
go build -o lazy-todo .

# Run the application
./lazy-todo

# Run with custom task file
./lazy-todo --file path/to/tasks.yaml

# Check version
./lazy-todo --version

# Build for multiple platforms (for releases)
GOOS=linux GOARCH=amd64 go build -o dist/lazy-todo-linux-amd64 .
GOOS=darwin GOARCH=amd64 go build -o dist/lazy-todo-darwin-amd64 .
GOOS=darwin GOARCH=arm64 go build -o dist/lazy-todo-darwin-arm64 .
GOOS=windows GOARCH=amd64 go build -o dist/lazy-todo-windows-amd64.exe .
```

## Architecture

### Bubble Tea Pattern
This application follows the Elm Architecture via Bubble Tea:
- **Model** (`internal/ui/app.go`): The `App` struct holds all application state
- **Update** (`App.Update()`): Handles all messages (key presses, task operations, window resizing)
- **View** (`App.View()`): Renders the current state to the terminal

### Key Components

**State Management**:
- `AppState`: Tracks current mode (Normal, Form, Help, Search, ConfirmDelete, TagInput)
- `ViewMode`: Switches between List and Kanban views
- Tasks are loaded/saved through the storage layer, with UI state synchronized on updates

**View Layer** (`internal/ui/`):
- `ListView`: Renders tasks as a scrollable list with filtering
- `KanbanView`: Renders tasks in 4 columns (todo, in_progress, blocked, done)
- `TaskForm`: Modal form for creating/editing tasks with tab navigation
- `HelpPanel`: Full keyboard shortcut reference
- Views receive tasks and maintain their own cursor/selection state

**Data Flow**:
1. User input → `App.Update()` → tea.Cmd
2. tea.Cmd executes storage operation → returns tea.Msg
3. Message updates model → `refreshViews()` syncs UI state
4. `App.View()` renders current state

### Storage Layer
- `storage.Storage`: Handles YAML file I/O at `~/.local/share/lazy-todo/tasks.yaml` or `./tasks.yaml`
- Operations (Add/Update/Delete) reload tasks and return the full updated slice
- Opening in editor uses `$EDITOR` or `$VISUAL` environment variable

### Styling
- Uses Catppuccin color palette (defined in `internal/ui/styles.go`)
- Priority and status have dedicated styles and icons
- Styles are passed down to all components for consistency

## Version Updates

To release a new version:
1. Update `version` variable in `main.go`
2. Build binaries for all platforms
3. Create git tag: `git tag -a vX.Y.Z -m "Release message"`
4. Push tag: `git push origin vX.Y.Z`
5. Create GitHub release with binaries: `gh release create vX.Y.Z dist/* --title "..." --notes "..."`

## Task File Format

Tasks are stored in YAML with this structure:
```yaml
tasks:
  - id: "uuid"
    title: "Task title"
    description: "Optional description"
    priority: low|medium|high|critical
    status: todo|in_progress|blocked|done
    tags: ["tag1", "tag2"]
    created_at: "2025-12-19T10:00:00Z"
    updated_at: "2025-12-19T14:00:00Z"
```

The file location is determined by `storage.DefaultFilePath()` which checks for `./tasks.yaml` first, then falls back to XDG data directory.
