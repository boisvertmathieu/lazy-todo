package ui

import (
	"strings"

	"lazy-todo/internal/model"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// FormField represents the current focused field
type FormField int

const (
	FieldTitle FormField = iota
	FieldDescription
	FieldTags
	FieldPriority
	FieldStatus
	FieldSubmit
	FieldCancel
)

// TaskForm is the form for creating/editing tasks
type TaskForm struct {
	task          *model.Task
	isNew         bool
	focusedField  FormField
	titleInput    textinput.Model
	descInput     textinput.Model
	tagsInput     textinput.Model
	priorityIdx   int
	statusIdx     int
	styles        Styles
	width, height int
}

// NewTaskForm creates a new task form
func NewTaskForm(styles Styles) *TaskForm {
	titleInput := textinput.New()
	titleInput.Placeholder = "Titre de la tâche"
	titleInput.Focus()
	titleInput.CharLimit = 100
	titleInput.Width = 40

	descInput := textinput.New()
	descInput.Placeholder = "Description (optionnel)"
	descInput.CharLimit = 500
	descInput.Width = 40

	tagsInput := textinput.New()
	tagsInput.Placeholder = "Tags séparés par des virgules"
	tagsInput.CharLimit = 100
	tagsInput.Width = 40

	return &TaskForm{
		titleInput:   titleInput,
		descInput:    descInput,
		tagsInput:    tagsInput,
		focusedField: FieldTitle,
		priorityIdx:  1, // Medium
		statusIdx:    0, // Todo
		styles:       styles,
	}
}

// SetTask sets the task to edit (nil for new task)
func (f *TaskForm) SetTask(task *model.Task) {
	if task == nil {
		f.isNew = true
		f.task = nil
		f.titleInput.SetValue("")
		f.descInput.SetValue("")
		f.tagsInput.SetValue("")
		f.priorityIdx = 1
		f.statusIdx = 0
	} else {
		f.isNew = false
		f.task = task
		f.titleInput.SetValue(task.Title)
		f.descInput.SetValue(task.Description)
		f.tagsInput.SetValue(strings.Join(task.Tags, ", "))

		// Set priority index
		priorities := model.AllPriorities()
		for i, p := range priorities {
			if p == task.Priority {
				f.priorityIdx = i
				break
			}
		}

		// Set status index
		statuses := model.AllStatuses()
		for i, s := range statuses {
			if s == task.Status {
				f.statusIdx = i
				break
			}
		}
	}

	f.focusedField = FieldTitle
	f.titleInput.Focus()
	f.descInput.Blur()
	f.tagsInput.Blur()
}

// SetSize sets the form dimensions
func (f *TaskForm) SetSize(width, height int) {
	f.width = width
	f.height = height
	inputWidth := width - 20
	if inputWidth > 60 {
		inputWidth = 60
	}
	f.titleInput.Width = inputWidth
	f.descInput.Width = inputWidth
	f.tagsInput.Width = inputWidth
}

// Update handles input
func (f *TaskForm) Update(msg tea.Msg) (*TaskForm, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down":
			f.nextField()
			return f, nil
		case "shift+tab", "up":
			f.prevField()
			return f, nil
		case "left":
			if f.focusedField == FieldPriority {
				if f.priorityIdx > 0 {
					f.priorityIdx--
				}
			} else if f.focusedField == FieldStatus {
				if f.statusIdx > 0 {
					f.statusIdx--
				}
			}
			return f, nil
		case "right":
			if f.focusedField == FieldPriority {
				if f.priorityIdx < len(model.AllPriorities())-1 {
					f.priorityIdx++
				}
			} else if f.focusedField == FieldStatus {
				if f.statusIdx < len(model.AllStatuses())-1 {
					f.statusIdx++
				}
			}
			return f, nil
		}
	}

	// Update the focused text input
	switch f.focusedField {
	case FieldTitle:
		f.titleInput, cmd = f.titleInput.Update(msg)
	case FieldDescription:
		f.descInput, cmd = f.descInput.Update(msg)
	case FieldTags:
		f.tagsInput, cmd = f.tagsInput.Update(msg)
	}

	return f, cmd
}

// nextField moves focus to the next field
func (f *TaskForm) nextField() {
	f.titleInput.Blur()
	f.descInput.Blur()
	f.tagsInput.Blur()

	f.focusedField++
	if f.focusedField > FieldCancel {
		f.focusedField = FieldTitle
	}

	switch f.focusedField {
	case FieldTitle:
		f.titleInput.Focus()
	case FieldDescription:
		f.descInput.Focus()
	case FieldTags:
		f.tagsInput.Focus()
	}
}

// prevField moves focus to the previous field
func (f *TaskForm) prevField() {
	f.titleInput.Blur()
	f.descInput.Blur()
	f.tagsInput.Blur()

	if f.focusedField == FieldTitle {
		f.focusedField = FieldCancel
	} else {
		f.focusedField--
	}

	switch f.focusedField {
	case FieldTitle:
		f.titleInput.Focus()
	case FieldDescription:
		f.descInput.Focus()
	case FieldTags:
		f.tagsInput.Focus()
	}
}

// GetTask returns the task with form values
func (f *TaskForm) GetTask() model.Task {
	var task model.Task
	if f.task != nil {
		task = *f.task
	} else {
		task = model.NewTask(f.titleInput.Value())
	}

	task.Title = f.titleInput.Value()
	task.Description = f.descInput.Value()

	// Parse tags
	tagStr := f.tagsInput.Value()
	if tagStr != "" {
		tags := strings.Split(tagStr, ",")
		task.Tags = make([]string, 0, len(tags))
		for _, t := range tags {
			t = strings.TrimSpace(t)
			if t != "" {
				task.Tags = append(task.Tags, t)
			}
		}
	} else {
		task.Tags = []string{}
	}

	priorities := model.AllPriorities()
	task.Priority = priorities[f.priorityIdx]

	statuses := model.AllStatuses()
	task.Status = statuses[f.statusIdx]

	return task
}

// IsValid returns true if the form is valid
func (f *TaskForm) IsValid() bool {
	return strings.TrimSpace(f.titleInput.Value()) != ""
}

// IsFocusedOnSubmit returns true if submit button is focused
func (f *TaskForm) IsFocusedOnSubmit() bool {
	return f.focusedField == FieldSubmit
}

// IsFocusedOnCancel returns true if cancel button is focused
func (f *TaskForm) IsFocusedOnCancel() bool {
	return f.focusedField == FieldCancel
}

// Render renders the form
func (f *TaskForm) Render() string {
	title := "Nouvelle tâche"
	if !f.isNew {
		title = "Modifier la tâche"
	}

	titleStyle := f.styles.DialogTitle
	labelStyle := f.styles.FormLabel

	var sections []string

	// Title
	sections = append(sections, titleStyle.Render(title))
	sections = append(sections, "")

	// Title field
	sections = append(sections, labelStyle.Render("Titre:"))
	sections = append(sections, f.renderInput(f.titleInput.View(), f.focusedField == FieldTitle))

	// Description field
	sections = append(sections, labelStyle.Render("Description:"))
	sections = append(sections, f.renderInput(f.descInput.View(), f.focusedField == FieldDescription))

	// Tags field
	sections = append(sections, labelStyle.Render("Tags:"))
	sections = append(sections, f.renderInput(f.tagsInput.View(), f.focusedField == FieldTags))

	// Priority selector
	sections = append(sections, labelStyle.Render("Priorité:"))
	sections = append(sections, f.renderPrioritySelector())

	// Status selector
	sections = append(sections, labelStyle.Render("État:"))
	sections = append(sections, f.renderStatusSelector())

	// Buttons
	sections = append(sections, "")
	sections = append(sections, f.renderButtons())

	content := strings.Join(sections, "\n")

	return f.styles.Dialog.Render(content)
}

// renderInput renders an input field
func (f *TaskForm) renderInput(view string, focused bool) string {
	if focused {
		return f.styles.FormInputFocus.Render(view)
	}
	return f.styles.FormInput.Render(view)
}

// renderPrioritySelector renders the priority selector
func (f *TaskForm) renderPrioritySelector() string {
	priorities := model.AllPriorities()
	var items []string

	for i, p := range priorities {
		icon := PriorityIcon(p)
		label := p.Label()
		style := f.styles.PriorityStyle(p)

		item := style.Render(icon + " " + label)
		if i == f.priorityIdx && f.focusedField == FieldPriority {
			item = lipgloss.NewStyle().
				Background(lipgloss.Color("#45475a")).
				Render("[" + icon + " " + label + "]")
		} else if i == f.priorityIdx {
			item = "[" + item + "]"
		}

		items = append(items, item)
	}

	return strings.Join(items, "  ")
}

// renderStatusSelector renders the status selector
func (f *TaskForm) renderStatusSelector() string {
	statuses := model.AllStatuses()
	var items []string

	for i, s := range statuses {
		icon := StatusIcon(s)
		label := s.Label()
		style := f.styles.StatusStyle(s)

		item := style.Render(icon + " " + label)
		if i == f.statusIdx && f.focusedField == FieldStatus {
			item = lipgloss.NewStyle().
				Background(lipgloss.Color("#45475a")).
				Render("[" + icon + " " + label + "]")
		} else if i == f.statusIdx {
			item = "[" + item + "]"
		}

		items = append(items, item)
	}

	return strings.Join(items, "  ")
}

// renderButtons renders the form buttons
func (f *TaskForm) renderButtons() string {
	submitStyle := f.styles.FormButton
	cancelStyle := f.styles.FormButton

	if f.focusedField == FieldSubmit {
		submitStyle = f.styles.FormButtonFocus
	}
	if f.focusedField == FieldCancel {
		cancelStyle = f.styles.FormButtonFocus
	}

	submit := submitStyle.Render("Valider")
	cancel := cancelStyle.Render("Annuler")

	return submit + "  " + cancel
}
