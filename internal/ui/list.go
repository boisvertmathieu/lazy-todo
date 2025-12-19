package ui

import (
	"fmt"
	"strings"

	"lazy-todo/internal/model"

	"github.com/charmbracelet/lipgloss"
)

// ListView represents the list view of tasks
type ListView struct {
	tasks    []model.Task
	cursor   int
	styles   Styles
	width    int
	height   int
	filter   string
	filtered []int // indices of filtered tasks
}

// NewListView creates a new list view
func NewListView(styles Styles) *ListView {
	return &ListView{
		tasks:    []model.Task{},
		cursor:   0,
		styles:   styles,
		filtered: []int{},
	}
}

// SetTasks sets the tasks to display
func (l *ListView) SetTasks(tasks []model.Task) {
	l.tasks = tasks
	l.applyFilter()
	if l.cursor >= len(l.filtered) {
		l.cursor = max(0, len(l.filtered)-1)
	}
}

// SetSize sets the view dimensions
func (l *ListView) SetSize(width, height int) {
	l.width = width
	l.height = height
}

// SetFilter sets the search filter
func (l *ListView) SetFilter(filter string) {
	l.filter = strings.ToLower(filter)
	l.applyFilter()
	l.cursor = 0
}

// applyFilter filters tasks based on the current filter
func (l *ListView) applyFilter() {
	l.filtered = []int{}
	for i, task := range l.tasks {
		if l.matchesFilter(task) {
			l.filtered = append(l.filtered, i)
		}
	}
}

// matchesFilter checks if a task matches the current filter
func (l *ListView) matchesFilter(task model.Task) bool {
	if l.filter == "" {
		return true
	}

	// Check title
	if strings.Contains(strings.ToLower(task.Title), l.filter) {
		return true
	}

	// Check description
	if strings.Contains(strings.ToLower(task.Description), l.filter) {
		return true
	}

	// Check tags
	for _, tag := range task.Tags {
		if strings.Contains(strings.ToLower(tag), l.filter) {
			return true
		}
	}

	return false
}

// MoveUp moves the cursor up
func (l *ListView) MoveUp() {
	if l.cursor > 0 {
		l.cursor--
	}
}

// MoveDown moves the cursor down
func (l *ListView) MoveDown() {
	if l.cursor < len(l.filtered)-1 {
		l.cursor++
	}
}

// SelectedTask returns the currently selected task
func (l *ListView) SelectedTask() *model.Task {
	if len(l.filtered) == 0 {
		return nil
	}
	if l.cursor >= 0 && l.cursor < len(l.filtered) {
		return &l.tasks[l.filtered[l.cursor]]
	}
	return nil
}

// SelectedIndex returns the index of the selected task in the original slice
func (l *ListView) SelectedIndex() int {
	if len(l.filtered) == 0 {
		return -1
	}
	if l.cursor >= 0 && l.cursor < len(l.filtered) {
		return l.filtered[l.cursor]
	}
	return -1
}

// Render renders the list view
func (l *ListView) Render() string {
	if len(l.filtered) == 0 {
		emptyMsg := "Aucune tâche"
		if l.filter != "" {
			emptyMsg = fmt.Sprintf("Aucun résultat pour \"%s\"", l.filter)
		}
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6c7086")).
			Italic(true).
			Padding(1, 2).
			Render(emptyMsg)
	}

	var lines []string
	visibleHeight := l.height - 2 // Account for padding

	// Calculate scroll offset
	scrollOffset := 0
	if l.cursor >= visibleHeight {
		scrollOffset = l.cursor - visibleHeight + 1
	}

	// Render visible items
	for i := scrollOffset; i < len(l.filtered) && i < scrollOffset+visibleHeight; i++ {
		taskIdx := l.filtered[i]
		task := l.tasks[taskIdx]
		isSelected := i == l.cursor

		line := l.renderTaskLine(task, isSelected)
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

// renderTaskLine renders a single task line
func (l *ListView) renderTaskLine(task model.Task, selected bool) string {
	// Priority icon
	priorityIcon := PriorityIcon(task.Priority)
	priorityStyle := l.styles.PriorityStyle(task.Priority)

	// Status icon
	statusIcon := StatusIcon(task.Status)
	statusStyle := l.styles.StatusStyle(task.Status)

	// Tags
	var tagStr string
	if len(task.Tags) > 0 {
		var tags []string
		for _, tag := range task.Tags {
			tags = append(tags, l.styles.Tag.Render(tag))
		}
		tagStr = " " + strings.Join(tags, " ")
	}

	// Build the line
	content := fmt.Sprintf(
		"%s %s %s%s",
		priorityStyle.Render(priorityIcon),
		statusStyle.Render(statusIcon),
		task.Title,
		tagStr,
	)

	// Truncate if too long
	maxWidth := l.width - 4
	if maxWidth > 0 && lipgloss.Width(content) > maxWidth {
		content = truncate(content, maxWidth)
	}

	// Apply selection style
	if selected {
		return l.styles.ListItemSelected.Width(l.width - 2).Render(content)
	}
	return l.styles.ListItem.Width(l.width - 2).Render(content)
}

// truncate truncates a string to a maximum width
func truncate(s string, maxWidth int) string {
	if lipgloss.Width(s) <= maxWidth {
		return s
	}

	// Simple truncation - could be improved for ANSI sequences
	runes := []rune(s)
	for i := len(runes) - 1; i >= 0; i-- {
		truncated := string(runes[:i]) + "…"
		if lipgloss.Width(truncated) <= maxWidth {
			return truncated
		}
	}
	return "…"
}

// Count returns the number of visible tasks
func (l *ListView) Count() int {
	return len(l.filtered)
}

// TotalCount returns the total number of tasks
func (l *ListView) TotalCount() int {
	return len(l.tasks)
}
