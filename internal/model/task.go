package model

import (
	"time"

	"github.com/google/uuid"
)

// Priority represents the priority level of a task
type Priority string

const (
	PriorityLow      Priority = "low"
	PriorityMedium   Priority = "medium"
	PriorityHigh     Priority = "high"
	PriorityCritical Priority = "critical"
)

// Status represents the current state of a task
type Status string

const (
	StatusTodo       Status = "todo"
	StatusInProgress Status = "in_progress"
	StatusBlocked    Status = "blocked"
	StatusDone       Status = "done"
)

// GroupBy represents the grouping criteria for tasks
type GroupBy int

const (
	GroupByNone GroupBy = iota
	GroupByStatus
	GroupByPriority
	GroupByTag
)

// AllGroupBy returns all available grouping options
func AllGroupBy() []GroupBy {
	return []GroupBy{GroupByNone, GroupByStatus, GroupByPriority, GroupByTag}
}

// Label returns the French label for a grouping option
func (g GroupBy) Label() string {
	switch g {
	case GroupByNone:
		return "Aucun"
	case GroupByStatus:
		return "État"
	case GroupByPriority:
		return "Priorité"
	case GroupByTag:
		return "Tag"
	default:
		return "Aucun"
	}
}

// Next cycles to the next grouping option
func (g GroupBy) Next() GroupBy {
	switch g {
	case GroupByNone:
		return GroupByStatus
	case GroupByStatus:
		return GroupByPriority
	case GroupByPriority:
		return GroupByTag
	case GroupByTag:
		return GroupByNone
	default:
		return GroupByNone
	}
}

// Task represents a single todo item
type Task struct {
	ID          string    `yaml:"id"`
	Title       string    `yaml:"title"`
	Description string    `yaml:"description,omitempty"`
	Priority    Priority  `yaml:"priority"`
	Status      Status    `yaml:"status"`
	Tags        []string  `yaml:"tags,omitempty"`
	CreatedAt   time.Time `yaml:"created_at"`
	UpdatedAt   time.Time `yaml:"updated_at"`
}

// TaskStore represents the root structure of the YAML file
type TaskStore struct {
	Tasks []Task `yaml:"tasks"`
}

// NewTask creates a new task with default values
func NewTask(title string) Task {
	now := time.Now()
	return Task{
		ID:        uuid.New().String(),
		Title:     title,
		Priority:  PriorityMedium,
		Status:    StatusTodo,
		Tags:      []string{},
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// AllPriorities returns all available priorities
func AllPriorities() []Priority {
	return []Priority{PriorityLow, PriorityMedium, PriorityHigh, PriorityCritical}
}

// AllStatuses returns all available statuses
func AllStatuses() []Status {
	return []Status{StatusTodo, StatusInProgress, StatusBlocked, StatusDone}
}

// PriorityLabel returns the French label for a priority
func (p Priority) Label() string {
	switch p {
	case PriorityLow:
		return "Basse"
	case PriorityMedium:
		return "Moyenne"
	case PriorityHigh:
		return "Haute"
	case PriorityCritical:
		return "Critique"
	default:
		return string(p)
	}
}

// StatusLabel returns the French label for a status
func (s Status) Label() string {
	switch s {
	case StatusTodo:
		return "À faire"
	case StatusInProgress:
		return "En cours"
	case StatusBlocked:
		return "Bloqué"
	case StatusDone:
		return "Terminé"
	default:
		return string(s)
	}
}

// StatusIndex returns the index of the status (for kanban columns)
func (s Status) Index() int {
	switch s {
	case StatusTodo:
		return 0
	case StatusInProgress:
		return 1
	case StatusBlocked:
		return 2
	case StatusDone:
		return 3
	default:
		return 0
	}
}

// StatusFromIndex returns the status for a given index
func StatusFromIndex(i int) Status {
	statuses := AllStatuses()
	if i >= 0 && i < len(statuses) {
		return statuses[i]
	}
	return StatusTodo
}

// NextPriority cycles to the next priority
func (p Priority) Next() Priority {
	switch p {
	case PriorityLow:
		return PriorityMedium
	case PriorityMedium:
		return PriorityHigh
	case PriorityHigh:
		return PriorityCritical
	case PriorityCritical:
		return PriorityLow
	default:
		return PriorityMedium
	}
}

// Index returns the index of the priority
func (p Priority) Index() int {
	switch p {
	case PriorityLow:
		return 0
	case PriorityMedium:
		return 1
	case PriorityHigh:
		return 2
	case PriorityCritical:
		return 3
	default:
		return 1
	}
}
