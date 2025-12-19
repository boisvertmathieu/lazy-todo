package storage

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"lazy-todo/internal/model"

	"gopkg.in/yaml.v3"
)

// Storage handles persistence of tasks to YAML file
type Storage struct {
	FilePath string
}

// NewStorage creates a new Storage instance
func NewStorage(filePath string) *Storage {
	return &Storage{FilePath: filePath}
}

// DefaultFilePath returns the default path for the tasks file
func DefaultFilePath() string {
	// First, check if tasks.yaml exists in current directory
	if _, err := os.Stat("tasks.yaml"); err == nil {
		absPath, _ := filepath.Abs("tasks.yaml")
		return absPath
	}

	// Otherwise, use XDG data directory or home directory
	dataDir := os.Getenv("XDG_DATA_HOME")
	if dataDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "tasks.yaml"
		}
		dataDir = filepath.Join(home, ".local", "share")
	}

	appDir := filepath.Join(dataDir, "lazy-todo")
	return filepath.Join(appDir, "tasks.yaml")
}

// Load reads tasks from the YAML file
func (s *Storage) Load() ([]model.Task, error) {
	data, err := os.ReadFile(s.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []model.Task{}, nil
		}
		return nil, err
	}

	var store model.TaskStore
	if err := yaml.Unmarshal(data, &store); err != nil {
		return nil, err
	}

	return store.Tasks, nil
}

// Save writes tasks to the YAML file
func (s *Storage) Save(tasks []model.Task) error {
	// Ensure directory exists
	dir := filepath.Dir(s.FilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	store := model.TaskStore{Tasks: tasks}
	data, err := yaml.Marshal(&store)
	if err != nil {
		return err
	}

	return os.WriteFile(s.FilePath, data, 0644)
}

// AddTask adds a new task and saves
func (s *Storage) AddTask(task model.Task) ([]model.Task, error) {
	tasks, err := s.Load()
	if err != nil {
		return nil, err
	}

	tasks = append(tasks, task)
	if err := s.Save(tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

// UpdateTask updates an existing task
func (s *Storage) UpdateTask(task model.Task) ([]model.Task, error) {
	tasks, err := s.Load()
	if err != nil {
		return nil, err
	}

	task.UpdatedAt = time.Now()

	for i, t := range tasks {
		if t.ID == task.ID {
			tasks[i] = task
			break
		}
	}

	if err := s.Save(tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

// DeleteTask removes a task by ID
func (s *Storage) DeleteTask(id string) ([]model.Task, error) {
	tasks, err := s.Load()
	if err != nil {
		return nil, err
	}

	var newTasks []model.Task
	for _, t := range tasks {
		if t.ID != id {
			newTasks = append(newTasks, t)
		}
	}

	if err := s.Save(newTasks); err != nil {
		return nil, err
	}

	return newTasks, nil
}

// OpenInEditor opens the YAML file in the default editor
func (s *Storage) OpenInEditor() error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = os.Getenv("VISUAL")
	}
	if editor == "" {
		// Default editors based on OS
		switch runtime.GOOS {
		case "windows":
			editor = "notepad"
		case "darwin":
			editor = "open"
		default:
			editor = "nano"
		}
	}

	cmd := exec.Command(editor, s.FilePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// GetFilePath returns the current file path
func (s *Storage) GetFilePath() string {
	return s.FilePath
}
