package ui

import (
	"fmt"
	"strings"

	"lazy-todo/internal/model"

	"github.com/charmbracelet/lipgloss"
)

// ListItem represents an item in the list (task or group header)
type ListItem struct {
	isHeader   bool
	headerText string
	taskIndex  int // index in the main tasks slice
}

// ListView represents the list view of tasks
type ListView struct {
	tasks    []model.Task
	cursor   int
	styles   Styles
	width    int
	height   int
	filter   string
	filtered []int      // indices of filtered tasks
	groupBy  model.GroupBy
	items    []ListItem // items to display (headers + tasks)
}

// NewListView creates a new list view
func NewListView(styles Styles) *ListView {
	return &ListView{
		tasks:    []model.Task{},
		cursor:   0,
		styles:   styles,
		filtered: []int{},
		groupBy:  model.GroupByNone,
		items:    []ListItem{},
	}
}

// SetTasks sets the tasks to display
func (l *ListView) SetTasks(tasks []model.Task) {
	l.tasks = tasks
	l.applyFilter()
	l.organizeItems()
	l.adjustCursor()
}

// SetGroupBy sets the grouping mode
func (l *ListView) SetGroupBy(groupBy model.GroupBy) {
	l.groupBy = groupBy
	l.organizeItems()
	l.adjustCursor()
}

// GetGroupBy returns the current grouping mode
func (l *ListView) GetGroupBy() model.GroupBy {
	return l.groupBy
}

// CycleGroupBy cycles to the next grouping mode
func (l *ListView) CycleGroupBy() {
	l.groupBy = l.groupBy.Next()
	l.organizeItems()
	l.adjustCursor()
}

// adjustCursor ensures cursor is on a valid task item
func (l *ListView) adjustCursor() {
	if len(l.items) == 0 {
		l.cursor = 0
		return
	}
	if l.cursor >= len(l.items) {
		l.cursor = len(l.items) - 1
	}
	// Move cursor to next task if on header
	for l.cursor < len(l.items) && l.items[l.cursor].isHeader {
		l.cursor++
	}
	// If we went past the end, move back to last task
	if l.cursor >= len(l.items) {
		for l.cursor = len(l.items) - 1; l.cursor >= 0; l.cursor-- {
			if !l.items[l.cursor].isHeader {
				break
			}
		}
	}
	if l.cursor < 0 {
		l.cursor = 0
	}
}

// organizeItems builds the items list based on groupBy setting
func (l *ListView) organizeItems() {
	l.items = []ListItem{}

	if l.groupBy == model.GroupByNone {
		// No grouping - just add all filtered tasks
		for _, idx := range l.filtered {
			l.items = append(l.items, ListItem{taskIndex: idx})
		}
		return
	}

	// Group tasks
	groups := make(map[string][]int)
	groupOrder := []string{}

	for _, idx := range l.filtered {
		task := l.tasks[idx]
		var key string

		switch l.groupBy {
		case model.GroupByStatus:
			key = task.Status.Label()
		case model.GroupByPriority:
			key = task.Priority.Label()
		case model.GroupByTag:
			if len(task.Tags) > 0 {
				key = task.Tags[0] // Group by first tag
			} else {
				key = "Sans tag"
			}
		}

		if _, exists := groups[key]; !exists {
			groupOrder = append(groupOrder, key)
		}
		groups[key] = append(groups[key], idx)
	}

	// Sort groups by their natural order for status and priority
	if l.groupBy == model.GroupByStatus {
		orderedKeys := []string{}
		for _, s := range model.AllStatuses() {
			if _, exists := groups[s.Label()]; exists {
				orderedKeys = append(orderedKeys, s.Label())
			}
		}
		groupOrder = orderedKeys
	} else if l.groupBy == model.GroupByPriority {
		orderedKeys := []string{}
		for _, p := range model.AllPriorities() {
			if _, exists := groups[p.Label()]; exists {
				orderedKeys = append(orderedKeys, p.Label())
			}
		}
		groupOrder = orderedKeys
	}

	// Build items with headers
	for _, groupKey := range groupOrder {
		taskIndices := groups[groupKey]
		// Add header
		l.items = append(l.items, ListItem{
			isHeader:   true,
			headerText: groupKey + " (" + itoa(len(taskIndices)) + ")",
		})
		// Add tasks
		for _, idx := range taskIndices {
			l.items = append(l.items, ListItem{taskIndex: idx})
		}
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
	l.organizeItems()
	l.cursor = 0
	l.adjustCursor()
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
		// Skip headers
		for l.cursor > 0 && l.items[l.cursor].isHeader {
			l.cursor--
		}
		// If we landed on a header at the top, go to next task
		if l.cursor >= 0 && l.cursor < len(l.items) && l.items[l.cursor].isHeader {
			for l.cursor < len(l.items) && l.items[l.cursor].isHeader {
				l.cursor++
			}
		}
	}
}

// MoveDown moves the cursor down
func (l *ListView) MoveDown() {
	if l.cursor < len(l.items)-1 {
		l.cursor++
		// Skip headers
		for l.cursor < len(l.items) && l.items[l.cursor].isHeader {
			l.cursor++
		}
		// If we went past the end, go back to last task
		if l.cursor >= len(l.items) {
			for l.cursor = len(l.items) - 1; l.cursor >= 0; l.cursor-- {
				if !l.items[l.cursor].isHeader {
					break
				}
			}
		}
	}
}

// SelectedTask returns the currently selected task
func (l *ListView) SelectedTask() *model.Task {
	if len(l.items) == 0 {
		return nil
	}
	if l.cursor >= 0 && l.cursor < len(l.items) {
		item := l.items[l.cursor]
		if item.isHeader {
			return nil
		}
		return &l.tasks[item.taskIndex]
	}
	return nil
}

// SelectedIndex returns the index of the selected task in the original slice
func (l *ListView) SelectedIndex() int {
	if len(l.items) == 0 {
		return -1
	}
	if l.cursor >= 0 && l.cursor < len(l.items) {
		item := l.items[l.cursor]
		if item.isHeader {
			return -1
		}
		return item.taskIndex
	}
	return -1
}

// Render renders the list view
func (l *ListView) Render() string {
	if len(l.items) == 0 {
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
	for i := scrollOffset; i < len(l.items) && i < scrollOffset+visibleHeight; i++ {
		item := l.items[i]
		if item.isHeader {
			line := l.renderGroupHeader(item.headerText)
			lines = append(lines, line)
		} else {
			task := l.tasks[item.taskIndex]
			isSelected := i == l.cursor
			line := l.renderTaskLine(task, isSelected)
			lines = append(lines, line)
		}
	}

	return strings.Join(lines, "\n")
}

// renderGroupHeader renders a group header
func (l *ListView) renderGroupHeader(text string) string {
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#cba6f7")).
		Bold(true).
		Padding(0, 1).
		MarginTop(1)

	return headerStyle.Width(l.width - 2).Render("▸ " + text)
}

// renderTaskLine renders a single task line
func (l *ListView) renderTaskLine(task model.Task, selected bool) string {
	// Priority icon
	priorityIcon := PriorityIcon(task.Priority)
	priorityStyle := l.styles.PriorityStyle(task.Priority)

	// Status icon
	statusIcon := StatusIcon(task.Status)
	statusStyle := l.styles.StatusStyle(task.Status)

	// Status label for right side
	statusLabel := task.Status.Label()
	statusLabelRendered := statusStyle.Render(statusLabel)
	statusLabelWidth := lipgloss.Width(statusLabelRendered)

	// Tags
	var tagStr string
	if len(task.Tags) > 0 {
		var tags []string
		for _, tag := range task.Tags {
			tags = append(tags, l.styles.Tag.Render(tag))
		}
		tagStr = " " + strings.Join(tags, " ")
	}

	// Build the left part of the line
	leftContent := fmt.Sprintf(
		"%s %s %s%s",
		priorityStyle.Render(priorityIcon),
		statusStyle.Render(statusIcon),
		task.Title,
		tagStr,
	)

	// Calculate available width for left content
	lineWidth := l.width - 4
	rightPartWidth := statusLabelWidth + 2 // padding for status label
	maxLeftWidth := lineWidth - rightPartWidth

	// Truncate left content if needed
	if maxLeftWidth > 0 && lipgloss.Width(leftContent) > maxLeftWidth {
		leftContent = truncate(leftContent, maxLeftWidth)
	}

	// Calculate spacing between left and right
	leftWidth := lipgloss.Width(leftContent)
	spacing := lineWidth - leftWidth - statusLabelWidth
	if spacing < 1 {
		spacing = 1
	}

	// Build full line with status on the right
	content := leftContent + strings.Repeat(" ", spacing) + statusLabelRendered

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
