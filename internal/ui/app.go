package ui

import (
	"fmt"
	"strings"
	"time"

	"lazy-todo/internal/keys"
	"lazy-todo/internal/model"
	"lazy-todo/internal/storage"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ViewMode represents the current view mode
type ViewMode int

const (
	ViewList ViewMode = iota
	ViewKanban
)

// AppState represents the current app state
type AppState int

const (
	StateNormal AppState = iota
	StateForm
	StateHelp
	StateSearch
	StateConfirmDelete
	StateTagInput
)

// App is the main application model
type App struct {
	storage    *storage.Storage
	tasks      []model.Task
	styles     Styles
	keys       keys.KeyMap
	viewMode   ViewMode
	state      AppState
	listView   *ListView
	kanbanView *KanbanView
	taskForm   *TaskForm
	helpPanel  *HelpPanel
	searchInput textinput.Model
	tagInput    textinput.Model
	width      int
	height     int
	err        error
	message    string
	messageTime time.Time
}

// NewApp creates a new App instance
func NewApp(store *storage.Storage) *App {
	styles := DefaultStyles()
	keyMap := keys.DefaultKeyMap()

	searchInput := textinput.New()
	searchInput.Placeholder = "Rechercher..."
	searchInput.CharLimit = 50

	tagInput := textinput.New()
	tagInput.Placeholder = "Nouveau tag..."
	tagInput.CharLimit = 30

	app := &App{
		storage:     store,
		tasks:       []model.Task{},
		styles:      styles,
		keys:        keyMap,
		viewMode:    ViewList,
		state:       StateNormal,
		listView:    NewListView(styles),
		kanbanView:  NewKanbanView(styles),
		taskForm:    NewTaskForm(styles),
		helpPanel:   NewHelpPanel(styles),
		searchInput: searchInput,
		tagInput:    tagInput,
	}

	return app
}

// Init initializes the app
func (a *App) Init() tea.Cmd {
	return tea.Batch(
		a.loadTasks,
		tea.EnterAltScreen,
	)
}

// loadTasks loads tasks from storage
func (a *App) loadTasks() tea.Msg {
	tasks, err := a.storage.Load()
	if err != nil {
		return errMsg{err}
	}
	return tasksLoadedMsg{tasks}
}

// Messages
type errMsg struct{ error }
type tasksLoadedMsg struct{ tasks []model.Task }
type tasksSavedMsg struct{}
type editorClosedMsg struct{ err error }

// Update handles messages and updates the model
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.updateSizes()
		return a, nil

	case errMsg:
		a.err = msg.error
		a.setMessage("Erreur: " + msg.Error())
		return a, nil

	case tasksLoadedMsg:
		a.tasks = msg.tasks
		a.refreshViews()
		return a, nil

	case tasksSavedMsg:
		a.setMessage("Tâches sauvegardées")
		return a, nil

	case editorClosedMsg:
		if msg.err != nil {
			a.setMessage("Erreur lors de l'ouverture de l'éditeur")
		}
		return a, a.loadTasks

	case tea.KeyMsg:
		return a.handleKeyPress(msg)
	}

	// Handle form updates
	if a.state == StateForm {
		var cmd tea.Cmd
		a.taskForm, cmd = a.taskForm.Update(msg)
		return a, cmd
	}

	// Handle search input
	if a.state == StateSearch {
		var cmd tea.Cmd
		a.searchInput, cmd = a.searchInput.Update(msg)
		a.listView.SetFilter(a.searchInput.Value())
		return a, cmd
	}

	// Handle tag input
	if a.state == StateTagInput {
		var cmd tea.Cmd
		a.tagInput, cmd = a.tagInput.Update(msg)
		return a, cmd
	}

	return a, nil
}

// handleKeyPress handles key press events
func (a *App) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Global keys
	if key.Matches(msg, a.keys.Quit) && a.state == StateNormal {
		return a, tea.Quit
	}

	// State-specific handling
	switch a.state {
	case StateForm:
		return a.handleFormKeys(msg)
	case StateHelp:
		return a.handleHelpKeys(msg)
	case StateSearch:
		return a.handleSearchKeys(msg)
	case StateConfirmDelete:
		return a.handleDeleteConfirmKeys(msg)
	case StateTagInput:
		return a.handleTagInputKeys(msg)
	default:
		return a.handleNormalKeys(msg)
	}
}

// handleNormalKeys handles keys in normal state
func (a *App) handleNormalKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	// Navigation
	case key.Matches(msg, a.keys.Up):
		a.moveUp()
	case key.Matches(msg, a.keys.Down):
		a.moveDown()
	case key.Matches(msg, a.keys.Left):
		if a.viewMode == ViewKanban {
			a.kanbanView.MoveLeft()
		}
	case key.Matches(msg, a.keys.Right):
		if a.viewMode == ViewKanban {
			a.kanbanView.MoveRight()
		}

	// Move task between columns
	case key.Matches(msg, a.keys.MoveLeft):
		if a.viewMode == ViewKanban {
			if task := a.kanbanView.MoveTaskLeft(); task != nil {
				return a, a.updateTask(*task)
			}
		}
	case key.Matches(msg, a.keys.MoveRight):
		if a.viewMode == ViewKanban {
			if task := a.kanbanView.MoveTaskRight(); task != nil {
				return a, a.updateTask(*task)
			}
		}

	// Actions
	case key.Matches(msg, a.keys.Add):
		a.taskForm.SetTask(nil)
		a.taskForm.SetSize(a.width, a.height)
		a.state = StateForm
	case key.Matches(msg, a.keys.Edit), key.Matches(msg, a.keys.Enter):
		if task := a.selectedTask(); task != nil {
			a.taskForm.SetTask(task)
			a.taskForm.SetSize(a.width, a.height)
			a.state = StateForm
		}
	case key.Matches(msg, a.keys.Delete):
		if a.selectedTask() != nil {
			a.state = StateConfirmDelete
		}
	case key.Matches(msg, a.keys.Priority):
		if task := a.selectedTask(); task != nil {
			task.Priority = task.Priority.Next()
			return a, a.updateTask(*task)
		}
	case key.Matches(msg, a.keys.Tag):
		if a.selectedTask() != nil {
			a.tagInput.SetValue("")
			a.tagInput.Focus()
			a.state = StateTagInput
		}

	// Quick status change
	case key.Matches(msg, a.keys.StatusTodo):
		return a, a.setTaskStatus(model.StatusTodo)
	case key.Matches(msg, a.keys.StatusInProgress):
		return a, a.setTaskStatus(model.StatusInProgress)
	case key.Matches(msg, a.keys.StatusBlocked):
		return a, a.setTaskStatus(model.StatusBlocked)
	case key.Matches(msg, a.keys.StatusDone):
		return a, a.setTaskStatus(model.StatusDone)

	// Views
	case key.Matches(msg, a.keys.ToggleView):
		if a.viewMode == ViewList {
			a.viewMode = ViewKanban
			// Sync selection
			if task := a.listView.SelectedTask(); task != nil {
				a.kanbanView.SetActiveColumn(task.Status.Index())
			}
		} else {
			a.viewMode = ViewList
		}
	case key.Matches(msg, a.keys.Search):
		a.searchInput.SetValue("")
		a.searchInput.Focus()
		a.state = StateSearch
	case key.Matches(msg, a.keys.Help):
		a.state = StateHelp
	case key.Matches(msg, a.keys.Refresh):
		return a, a.loadTasks
	case key.Matches(msg, a.keys.OpenEditor):
		return a, a.openEditor()
	}

	return a, nil
}

// handleFormKeys handles keys in form state
func (a *App) handleFormKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		a.state = StateNormal
		return a, nil
	case "enter":
		if a.taskForm.IsFocusedOnSubmit() {
			if a.taskForm.IsValid() {
				task := a.taskForm.GetTask()
				a.state = StateNormal
				if a.taskForm.isNew {
					return a, a.addTask(task)
				}
				return a, a.updateTask(task)
			}
		} else if a.taskForm.IsFocusedOnCancel() {
			a.state = StateNormal
			return a, nil
		}
	}

	var cmd tea.Cmd
	a.taskForm, cmd = a.taskForm.Update(msg)
	return a, cmd
}

// handleHelpKeys handles keys in help state
func (a *App) handleHelpKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, a.keys.Help), msg.String() == "esc", msg.String() == "q":
		a.state = StateNormal
	}
	return a, nil
}

// handleSearchKeys handles keys in search state
func (a *App) handleSearchKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		a.searchInput.SetValue("")
		a.listView.SetFilter("")
		a.state = StateNormal
		return a, nil
	case "enter":
		a.state = StateNormal
		return a, nil
	}

	var cmd tea.Cmd
	a.searchInput, cmd = a.searchInput.Update(msg)
	a.listView.SetFilter(a.searchInput.Value())
	return a, cmd
}

// handleDeleteConfirmKeys handles delete confirmation
func (a *App) handleDeleteConfirmKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		a.state = StateNormal
		return a, a.deleteSelectedTask()
	case "n", "N", "esc":
		a.state = StateNormal
	}
	return a, nil
}

// handleTagInputKeys handles tag input
func (a *App) handleTagInputKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		a.state = StateNormal
		return a, nil
	case "enter":
		tag := strings.TrimSpace(a.tagInput.Value())
		if tag != "" {
			if task := a.selectedTask(); task != nil {
				// Toggle tag
				found := false
				newTags := []string{}
				for _, t := range task.Tags {
					if t == tag {
						found = true
					} else {
						newTags = append(newTags, t)
					}
				}
				if !found {
					newTags = append(newTags, tag)
				}
				task.Tags = newTags
				a.state = StateNormal
				return a, a.updateTask(*task)
			}
		}
		a.state = StateNormal
		return a, nil
	}

	var cmd tea.Cmd
	a.tagInput, cmd = a.tagInput.Update(msg)
	return a, cmd
}

// selectedTask returns the currently selected task
func (a *App) selectedTask() *model.Task {
	if a.viewMode == ViewList {
		return a.listView.SelectedTask()
	}
	return a.kanbanView.SelectedTask()
}

// selectedIndex returns the index of the selected task
func (a *App) selectedIndex() int {
	if a.viewMode == ViewList {
		return a.listView.SelectedIndex()
	}
	return a.kanbanView.SelectedIndex()
}

// moveUp moves selection up
func (a *App) moveUp() {
	if a.viewMode == ViewList {
		a.listView.MoveUp()
	} else {
		a.kanbanView.MoveUp()
	}
}

// moveDown moves selection down
func (a *App) moveDown() {
	if a.viewMode == ViewList {
		a.listView.MoveDown()
	} else {
		a.kanbanView.MoveDown()
	}
}

// updateSizes updates component sizes
func (a *App) updateSizes() {
	contentHeight := a.height - 4 // Header + Footer
	a.listView.SetSize(a.width, contentHeight)
	a.kanbanView.SetSize(a.width, contentHeight)
	a.taskForm.SetSize(a.width, a.height)
	a.helpPanel.SetSize(a.width-10, a.height-10)
}

// refreshViews refreshes all views with current tasks
func (a *App) refreshViews() {
	a.listView.SetTasks(a.tasks)
	a.kanbanView.SetTasks(a.tasks)
}

// setMessage sets a temporary status message
func (a *App) setMessage(msg string) {
	a.message = msg
	a.messageTime = time.Now()
}

// Task operations

func (a *App) addTask(task model.Task) tea.Cmd {
	return func() tea.Msg {
		tasks, err := a.storage.AddTask(task)
		if err != nil {
			return errMsg{err}
		}
		return tasksLoadedMsg{tasks}
	}
}

func (a *App) updateTask(task model.Task) tea.Cmd {
	return func() tea.Msg {
		tasks, err := a.storage.UpdateTask(task)
		if err != nil {
			return errMsg{err}
		}
		return tasksLoadedMsg{tasks}
	}
}

func (a *App) deleteSelectedTask() tea.Cmd {
	task := a.selectedTask()
	if task == nil {
		return nil
	}
	return func() tea.Msg {
		tasks, err := a.storage.DeleteTask(task.ID)
		if err != nil {
			return errMsg{err}
		}
		return tasksLoadedMsg{tasks}
	}
}

func (a *App) setTaskStatus(status model.Status) tea.Cmd {
	task := a.selectedTask()
	if task == nil {
		return nil
	}
	task.Status = status
	return a.updateTask(*task)
}

func (a *App) openEditor() tea.Cmd {
	return func() tea.Msg {
		err := a.storage.OpenInEditor()
		return editorClosedMsg{err}
	}
}

// View renders the app
func (a *App) View() string {
	if a.width == 0 || a.height == 0 {
		return "Chargement..."
	}

	var content string

	switch a.state {
	case StateForm:
		content = a.renderFormOverlay()
	case StateHelp:
		content = a.renderHelpOverlay()
	case StateConfirmDelete:
		content = a.renderDeleteConfirm()
	case StateTagInput:
		content = a.renderTagInput()
	default:
		content = a.renderMainView()
	}

	return content
}

// renderMainView renders the main view
func (a *App) renderMainView() string {
	var sections []string

	// Header
	sections = append(sections, a.renderHeader())

	// Content
	contentHeight := a.height - 4
	var viewContent string
	if a.viewMode == ViewList {
		viewContent = a.listView.Render()
	} else {
		viewContent = a.kanbanView.Render()
	}

	// Add search bar if searching
	if a.state == StateSearch {
		searchBar := a.styles.FormInputFocus.Render("/ " + a.searchInput.View())
		viewContent = searchBar + "\n" + viewContent
	}

	contentStyle := lipgloss.NewStyle().
		Height(contentHeight).
		Width(a.width)
	sections = append(sections, contentStyle.Render(viewContent))

	// Footer
	sections = append(sections, RenderFooter(a.styles, a.viewMode == ViewKanban))

	return strings.Join(sections, "\n")
}

// renderHeader renders the header
func (a *App) renderHeader() string {
	title := a.styles.HeaderTitle.Render("lazy-todo")

	// File path
	filePath := a.storage.GetFilePath()
	if len(filePath) > 40 {
		filePath = "..." + filePath[len(filePath)-37:]
	}
	fileInfo := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6c7086")).
		Italic(true).
		Render(filePath)

	// View tabs
	listTab := a.styles.HeaderTab
	kanbanTab := a.styles.HeaderTab
	if a.viewMode == ViewList {
		listTab = a.styles.HeaderTabSel
	} else {
		kanbanTab = a.styles.HeaderTabSel
	}

	tabs := listTab.Render("Liste") + " " + kanbanTab.Render("Kanban")

	// Task count
	count := fmt.Sprintf("%d tâches", len(a.tasks))
	countStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#a6adc8"))

	leftSide := title + "  " + fileInfo
	rightSide := countStyle.Render(count) + "  " + tabs

	// Calculate spacing
	gap := a.width - lipgloss.Width(leftSide) - lipgloss.Width(rightSide) - 2
	if gap < 1 {
		gap = 1
	}

	return a.styles.Header.Width(a.width).Render(
		leftSide + strings.Repeat(" ", gap) + rightSide,
	)
}

// renderFormOverlay renders the form overlay
func (a *App) renderFormOverlay() string {
	formView := a.taskForm.Render()

	// Center the form
	formWidth := lipgloss.Width(formView)
	formHeight := lipgloss.Height(formView)

	x := (a.width - formWidth) / 2
	y := (a.height - formHeight) / 2

	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}

	// Create overlay
	overlay := lipgloss.Place(
		a.width, a.height,
		lipgloss.Center, lipgloss.Center,
		formView,
	)

	return overlay
}

// renderHelpOverlay renders the help overlay
func (a *App) renderHelpOverlay() string {
	helpView := a.helpPanel.Render()

	return lipgloss.Place(
		a.width, a.height,
		lipgloss.Center, lipgloss.Center,
		helpView,
	)
}

// renderDeleteConfirm renders the delete confirmation dialog
func (a *App) renderDeleteConfirm() string {
	task := a.selectedTask()
	if task == nil {
		return a.renderMainView()
	}

	title := a.styles.DialogTitle.Render("Supprimer la tâche?")
	taskTitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#cdd6f4")).
		Render(task.Title)

	buttons := a.styles.FormButton.Render("(Y)es") + "  " +
		a.styles.FormButtonFocus.Render("(N)o")

	content := title + "\n\n" + taskTitle + "\n\n" + buttons

	dialog := a.styles.Dialog.Render(content)

	return lipgloss.Place(
		a.width, a.height,
		lipgloss.Center, lipgloss.Center,
		dialog,
	)
}

// renderTagInput renders the tag input dialog
func (a *App) renderTagInput() string {
	task := a.selectedTask()
	if task == nil {
		return a.renderMainView()
	}

	title := a.styles.DialogTitle.Render("Ajouter/Retirer un tag")

	// Show current tags
	var tagList string
	if len(task.Tags) > 0 {
		tags := strings.Join(task.Tags, ", ")
		tagList = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#a6adc8")).
			Italic(true).
			Render("Tags actuels: " + tags)
	} else {
		tagList = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6c7086")).
			Italic(true).
			Render("Aucun tag")
	}

	input := a.styles.FormInputFocus.Render(a.tagInput.View())

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6c7086")).
		Render("Enter: ajouter/retirer, Esc: annuler")

	content := title + "\n\n" + tagList + "\n\n" + input + "\n\n" + help

	dialog := a.styles.Dialog.Render(content)

	return lipgloss.Place(
		a.width, a.height,
		lipgloss.Center, lipgloss.Center,
		dialog,
	)
}
