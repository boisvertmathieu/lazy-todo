package ui

import (
	"strings"

	"lazy-todo/internal/model"

	"github.com/charmbracelet/lipgloss"
)

// KanbanItem represents an item in a column (task or group header)
type KanbanItem struct {
	isHeader   bool
	headerText string
	taskIndex  int // index in the main tasks slice
}

// KanbanColumn represents a single column in the kanban board
type KanbanColumn struct {
	status model.Status
	tasks  []int        // indices in the main tasks slice
	items  []KanbanItem // items to display (headers + tasks)
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
	groupBy     model.GroupBy
}

// NewKanbanView creates a new kanban view
func NewKanbanView(styles Styles) *KanbanView {
	return &KanbanView{
		tasks:     []model.Task{},
		activeCol: 0,
		styles:    styles,
		groupBy:   model.GroupByNone,
		columns: [4]KanbanColumn{
			{status: model.StatusTodo, tasks: []int{}, items: []KanbanItem{}, cursor: 0},
			{status: model.StatusInProgress, tasks: []int{}, items: []KanbanItem{}, cursor: 0},
			{status: model.StatusBlocked, tasks: []int{}, items: []KanbanItem{}, cursor: 0},
			{status: model.StatusDone, tasks: []int{}, items: []KanbanItem{}, cursor: 0},
		},
	}
}

// SetGroupBy sets the grouping mode
func (k *KanbanView) SetGroupBy(groupBy model.GroupBy) {
	k.groupBy = groupBy
	k.organizeItems()
	k.adjustCursors()
}

// GetGroupBy returns the current grouping mode
func (k *KanbanView) GetGroupBy() model.GroupBy {
	return k.groupBy
}

// CycleGroupBy cycles to the next grouping mode
func (k *KanbanView) CycleGroupBy() {
	k.groupBy = k.groupBy.Next()
	k.organizeItems()
	k.adjustCursors()
}

// adjustCursors ensures cursors are on valid task items in all columns
func (k *KanbanView) adjustCursors() {
	for i := range k.columns {
		k.adjustColumnCursor(i)
	}
}

// adjustColumnCursor ensures cursor is on a valid task item in a column
func (k *KanbanView) adjustColumnCursor(colIdx int) {
	col := &k.columns[colIdx]
	if len(col.items) == 0 {
		col.cursor = 0
		return
	}
	if col.cursor >= len(col.items) {
		col.cursor = len(col.items) - 1
	}
	// Move cursor to next task if on header
	for col.cursor < len(col.items) && col.items[col.cursor].isHeader {
		col.cursor++
	}
	// If we went past the end, move back to last task
	if col.cursor >= len(col.items) {
		for col.cursor = len(col.items) - 1; col.cursor >= 0; col.cursor-- {
			if !col.items[col.cursor].isHeader {
				break
			}
		}
	}
	if col.cursor < 0 {
		col.cursor = 0
	}
}

// organizeItems organizes items within each column based on groupBy
func (k *KanbanView) organizeItems() {
	for i := range k.columns {
		k.organizeColumnItems(i)
	}
}

// organizeColumnItems organizes items in a single column
func (k *KanbanView) organizeColumnItems(colIdx int) {
	col := &k.columns[colIdx]
	col.items = []KanbanItem{}

	if k.groupBy == model.GroupByNone || k.groupBy == model.GroupByStatus {
		// No grouping within column - just add all tasks
		for _, idx := range col.tasks {
			col.items = append(col.items, KanbanItem{taskIndex: idx})
		}
		return
	}

	// Group tasks within the column
	groups := make(map[string][]int)
	groupOrder := []string{}

	for _, idx := range col.tasks {
		task := k.tasks[idx]
		var key string

		switch k.groupBy {
		case model.GroupByPriority:
			key = task.Priority.Label()
		case model.GroupByTag:
			if len(task.Tags) > 0 {
				key = task.Tags[0]
			} else {
				key = "Sans tag"
			}
		}

		if _, exists := groups[key]; !exists {
			groupOrder = append(groupOrder, key)
		}
		groups[key] = append(groups[key], idx)
	}

	// Sort groups by their natural order for priority
	if k.groupBy == model.GroupByPriority {
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
		col.items = append(col.items, KanbanItem{
			isHeader:   true,
			headerText: groupKey,
		})
		// Add tasks
		for _, idx := range taskIndices {
			col.items = append(col.items, KanbanItem{taskIndex: idx})
		}
	}
}

// SetTasks sets the tasks to display
func (k *KanbanView) SetTasks(tasks []model.Task) {
	k.tasks = tasks
	k.organizeTasks()
	k.organizeItems()
	k.adjustCursors()
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
		// Skip headers
		for col.cursor > 0 && col.items[col.cursor].isHeader {
			col.cursor--
		}
		// If we landed on a header at the top, go to next task
		if col.cursor >= 0 && col.cursor < len(col.items) && col.items[col.cursor].isHeader {
			for col.cursor < len(col.items) && col.items[col.cursor].isHeader {
				col.cursor++
			}
		}
	}
}

// MoveDown moves the cursor down in the current column
func (k *KanbanView) MoveDown() {
	col := &k.columns[k.activeCol]
	if col.cursor < len(col.items)-1 {
		col.cursor++
		// Skip headers
		for col.cursor < len(col.items) && col.items[col.cursor].isHeader {
			col.cursor++
		}
		// If we went past the end, go back to last task
		if col.cursor >= len(col.items) {
			for col.cursor = len(col.items) - 1; col.cursor >= 0; col.cursor-- {
				if !col.items[col.cursor].isHeader {
					break
				}
			}
		}
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
	if len(col.items) == 0 {
		return nil
	}
	if col.cursor >= 0 && col.cursor < len(col.items) {
		item := col.items[col.cursor]
		if item.isHeader {
			return nil
		}
		return &k.tasks[item.taskIndex]
	}
	return nil
}

// SelectedIndex returns the index of the selected task in the original slice
func (k *KanbanView) SelectedIndex() int {
	col := k.columns[k.activeCol]
	if len(col.items) == 0 {
		return -1
	}
	if col.cursor >= 0 && col.cursor < len(col.items) {
		item := col.items[col.cursor]
		if item.isHeader {
			return -1
		}
		return item.taskIndex
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

	// Render items (cards and headers)
	var items []string
	cardHeight := k.height - 6 // Account for title and borders

	// Calculate visible range
	visibleItems := cardHeight / 3 // Approximate items per column (headers are smaller)
	if visibleItems < 1 {
		visibleItems = 1
	}

	scrollOffset := 0
	if col.cursor >= visibleItems {
		scrollOffset = col.cursor - visibleItems + 1
	}

	for i := scrollOffset; i < len(col.items) && i < scrollOffset+visibleItems; i++ {
		item := col.items[i]
		if item.isHeader {
			header := k.renderGroupHeader(item.headerText)
			items = append(items, header)
		} else {
			task := k.tasks[item.taskIndex]
			isSelected := isActive && i == col.cursor
			card := k.renderCard(task, isSelected)
			items = append(items, card)
		}
	}

	content := titleText + "\n" + strings.Join(items, "\n")

	// Apply column style
	var colStyle lipgloss.Style
	if isActive {
		colStyle = k.styles.KanbanColumnSelected.Width(k.columnWidth).Height(k.height - 4)
	} else {
		colStyle = k.styles.KanbanColumn.Width(k.columnWidth).Height(k.height - 4)
	}

	return colStyle.Render(content)
}

// renderGroupHeader renders a group header within a column
func (k *KanbanView) renderGroupHeader(text string) string {
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#cba6f7")).
		Bold(true).
		Italic(true)

	return headerStyle.Render("─ " + text + " ─")
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
