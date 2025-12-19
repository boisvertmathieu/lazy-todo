package ui

import (
	"strings"

	"lazy-todo/internal/model"

	"github.com/charmbracelet/lipgloss"
)

// KanbanColumn represents a single column in the kanban board
type KanbanColumn struct {
	status model.Status
	tasks  []int // indices in the main tasks slice
	cursor int
}

// KanbanView represents the kanban board view
type KanbanView struct {
	tasks       []model.Task
	columns     [4]KanbanColumn
	activeCol   int
	styles      Styles
	width       int
	height      int
	columnWidth int
}

// NewKanbanView creates a new kanban view
func NewKanbanView(styles Styles) *KanbanView {
	return &KanbanView{
		tasks:     []model.Task{},
		activeCol: 0,
		styles:    styles,
		columns: [4]KanbanColumn{
			{status: model.StatusTodo, tasks: []int{}, cursor: 0},
			{status: model.StatusInProgress, tasks: []int{}, cursor: 0},
			{status: model.StatusBlocked, tasks: []int{}, cursor: 0},
			{status: model.StatusDone, tasks: []int{}, cursor: 0},
		},
	}
}

// SetTasks sets the tasks to display
func (k *KanbanView) SetTasks(tasks []model.Task) {
	k.tasks = tasks
	k.organizeTasks()
}

// organizeTasks organizes tasks into columns
func (k *KanbanView) organizeTasks() {
	// Reset columns
	for i := range k.columns {
		k.columns[i].tasks = []int{}
	}

	// Distribute tasks to columns
	for i, task := range k.tasks {
		colIdx := task.Status.Index()
		if colIdx >= 0 && colIdx < 4 {
			k.columns[colIdx].tasks = append(k.columns[colIdx].tasks, i)
		}
	}

	// Adjust cursors
	for i := range k.columns {
		if k.columns[i].cursor >= len(k.columns[i].tasks) {
			k.columns[i].cursor = max(0, len(k.columns[i].tasks)-1)
		}
	}
}

// SetSize sets the view dimensions
func (k *KanbanView) SetSize(width, height int) {
	k.width = width
	k.height = height
	// Calculate column width (4 columns with gaps)
	k.columnWidth = (width - 12) / 4
	if k.columnWidth < 20 {
		k.columnWidth = 20
	}
}

// MoveUp moves the cursor up in the current column
func (k *KanbanView) MoveUp() {
	col := &k.columns[k.activeCol]
	if col.cursor > 0 {
		col.cursor--
	}
}

// MoveDown moves the cursor down in the current column
func (k *KanbanView) MoveDown() {
	col := &k.columns[k.activeCol]
	if col.cursor < len(col.tasks)-1 {
		col.cursor++
	}
}

// MoveLeft moves to the previous column
func (k *KanbanView) MoveLeft() {
	if k.activeCol > 0 {
		k.activeCol--
	}
}

// MoveRight moves to the next column
func (k *KanbanView) MoveRight() {
	if k.activeCol < 3 {
		k.activeCol++
	}
}

// MoveTaskLeft moves the selected task to the previous column
func (k *KanbanView) MoveTaskLeft() *model.Task {
	if k.activeCol == 0 {
		return nil
	}

	task := k.SelectedTask()
	if task == nil {
		return nil
	}

	newStatus := model.StatusFromIndex(k.activeCol - 1)
	task.Status = newStatus
	return task
}

// MoveTaskRight moves the selected task to the next column
func (k *KanbanView) MoveTaskRight() *model.Task {
	if k.activeCol >= 3 {
		return nil
	}

	task := k.SelectedTask()
	if task == nil {
		return nil
	}

	newStatus := model.StatusFromIndex(k.activeCol + 1)
	task.Status = newStatus
	return task
}

// SelectedTask returns the currently selected task
func (k *KanbanView) SelectedTask() *model.Task {
	col := k.columns[k.activeCol]
	if len(col.tasks) == 0 {
		return nil
	}
	if col.cursor >= 0 && col.cursor < len(col.tasks) {
		return &k.tasks[col.tasks[col.cursor]]
	}
	return nil
}

// SelectedIndex returns the index of the selected task in the original slice
func (k *KanbanView) SelectedIndex() int {
	col := k.columns[k.activeCol]
	if len(col.tasks) == 0 {
		return -1
	}
	if col.cursor >= 0 && col.cursor < len(col.tasks) {
		return col.tasks[col.cursor]
	}
	return -1
}

// Render renders the kanban board
func (k *KanbanView) Render() string {
	var columns []string

	for i := 0; i < 4; i++ {
		col := k.renderColumn(i)
		columns = append(columns, col)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, columns...)
}

// renderColumn renders a single column
func (k *KanbanView) renderColumn(colIdx int) string {
	col := k.columns[colIdx]
	isActive := colIdx == k.activeCol

	// Column title
	title := col.status.Label()
	count := len(col.tasks)
	titleText := k.styles.KanbanColumnTitle.Render(title + " (" + itoa(count) + ")")

	// Render cards
	var cards []string
	cardHeight := k.height - 6 // Account for title and borders

	// Calculate visible range
	visibleCards := cardHeight / 4 // Approximate cards per column
	if visibleCards < 1 {
		visibleCards = 1
	}

	scrollOffset := 0
	if col.cursor >= visibleCards {
		scrollOffset = col.cursor - visibleCards + 1
	}

	for i := scrollOffset; i < len(col.tasks) && i < scrollOffset+visibleCards; i++ {
		taskIdx := col.tasks[i]
		task := k.tasks[taskIdx]
		isSelected := isActive && i == col.cursor

		card := k.renderCard(task, isSelected)
		cards = append(cards, card)
	}

	content := titleText + "\n" + strings.Join(cards, "\n")

	// Apply column style
	var colStyle lipgloss.Style
	if isActive {
		colStyle = k.styles.KanbanColumnSelected.Width(k.columnWidth).Height(k.height - 4)
	} else {
		colStyle = k.styles.KanbanColumn.Width(k.columnWidth).Height(k.height - 4)
	}

	return colStyle.Render(content)
}

// renderCard renders a single task card
func (k *KanbanView) renderCard(task model.Task, selected bool) string {
	// Priority icon
	priorityIcon := PriorityIcon(task.Priority)
	priorityStyle := k.styles.PriorityStyle(task.Priority)

	// Title (truncated)
	title := task.Title
	maxTitleLen := k.columnWidth - 8
	if len(title) > maxTitleLen {
		title = title[:maxTitleLen-1] + "…"
	}

	// Tags (first 2 only)
	var tagStr string
	if len(task.Tags) > 0 {
		maxTags := 2
		if len(task.Tags) < maxTags {
			maxTags = len(task.Tags)
		}
		tagStr = strings.Join(task.Tags[:maxTags], ", ")
		if len(task.Tags) > 2 {
			tagStr += "…"
		}
	}

	// Build card content
	var lines []string
	lines = append(lines, priorityStyle.Render(priorityIcon)+" "+title)
	if tagStr != "" {
		tagLine := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6c7086")).
			Italic(true).
			Render(tagStr)
		lines = append(lines, tagLine)
	}

	content := strings.Join(lines, "\n")

	// Apply card style
	cardWidth := k.columnWidth - 4
	if selected {
		return k.styles.KanbanCardSelected.Width(cardWidth).Render(content)
	}
	return k.styles.KanbanCard.Width(cardWidth).Render(content)
}

// SetActiveColumn sets the active column
func (k *KanbanView) SetActiveColumn(col int) {
	if col >= 0 && col < 4 {
		k.activeCol = col
	}
}

// ActiveColumn returns the active column index
func (k *KanbanView) ActiveColumn() int {
	return k.activeCol
}

// itoa converts int to string without importing strconv
func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	if i < 0 {
		return "-" + itoa(-i)
	}
	var result []byte
	for i > 0 {
		result = append([]byte{byte('0' + i%10)}, result...)
		i /= 10
	}
	return string(result)
}
