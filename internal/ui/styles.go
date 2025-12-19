package ui

import (
	"lazy-todo/internal/model"

	"github.com/charmbracelet/lipgloss"
)

// Colors - using a catppuccin-inspired palette
var (
	colorRosewater = lipgloss.Color("#f5e0dc")
	colorFlamingo  = lipgloss.Color("#f2cdcd")
	colorPink      = lipgloss.Color("#f5c2e7")
	colorMauve     = lipgloss.Color("#cba6f7")
	colorRed       = lipgloss.Color("#f38ba8")
	colorMaroon    = lipgloss.Color("#eba0ac")
	colorPeach     = lipgloss.Color("#fab387")
	colorYellow    = lipgloss.Color("#f9e2af")
	colorGreen     = lipgloss.Color("#a6e3a1")
	colorTeal      = lipgloss.Color("#94e2d5")
	colorSky       = lipgloss.Color("#89dceb")
	colorSapphire  = lipgloss.Color("#74c7ec")
	colorBlue      = lipgloss.Color("#89b4fa")
	colorLavender  = lipgloss.Color("#b4befe")
	colorText      = lipgloss.Color("#cdd6f4")
	colorSubtext1  = lipgloss.Color("#bac2de")
	colorSubtext0  = lipgloss.Color("#a6adc8")
	colorOverlay2  = lipgloss.Color("#9399b2")
	colorOverlay1  = lipgloss.Color("#7f849c")
	colorOverlay0  = lipgloss.Color("#6c7086")
	colorSurface2  = lipgloss.Color("#585b70")
	colorSurface1  = lipgloss.Color("#45475a")
	colorSurface0  = lipgloss.Color("#313244")
	colorBase      = lipgloss.Color("#1e1e2e")
	colorMantle    = lipgloss.Color("#181825")
	colorCrust     = lipgloss.Color("#11111b")
)

// Styles holds all the application styles
type Styles struct {
	// App
	App lipgloss.Style

	// Header
	Header       lipgloss.Style
	HeaderTitle  lipgloss.Style
	HeaderTab    lipgloss.Style
	HeaderTabSel lipgloss.Style

	// List view
	ListItem         lipgloss.Style
	ListItemSelected lipgloss.Style
	ListItemTitle    lipgloss.Style
	ListItemDesc     lipgloss.Style

	// Kanban view
	KanbanColumn         lipgloss.Style
	KanbanColumnSelected lipgloss.Style
	KanbanColumnTitle    lipgloss.Style
	KanbanCard           lipgloss.Style
	KanbanCardSelected   lipgloss.Style

	// Priority colors
	PriorityLow      lipgloss.Style
	PriorityMedium   lipgloss.Style
	PriorityHigh     lipgloss.Style
	PriorityCritical lipgloss.Style

	// Status colors
	StatusTodo       lipgloss.Style
	StatusInProgress lipgloss.Style
	StatusBlocked    lipgloss.Style
	StatusDone       lipgloss.Style

	// Tags
	Tag lipgloss.Style

	// Footer/Help
	Footer    lipgloss.Style
	HelpKey   lipgloss.Style
	HelpValue lipgloss.Style
	HelpSep   lipgloss.Style

	// Form
	FormLabel       lipgloss.Style
	FormInput       lipgloss.Style
	FormInputFocus  lipgloss.Style
	FormButton      lipgloss.Style
	FormButtonFocus lipgloss.Style

	// Help panel
	HelpPanel      lipgloss.Style
	HelpPanelTitle lipgloss.Style

	// Dialog
	Dialog      lipgloss.Style
	DialogTitle lipgloss.Style

	// Borders
	Border lipgloss.Border
}

// DefaultStyles returns the default application styles
func DefaultStyles() Styles {
	s := Styles{}

	// App container
	s.App = lipgloss.NewStyle().
		Background(colorBase)

	// Header
	s.Header = lipgloss.NewStyle().
		Background(colorMantle).
		Foreground(colorText).
		Padding(0, 1).
		Bold(true)

	s.HeaderTitle = lipgloss.NewStyle().
		Foreground(colorMauve).
		Bold(true)

	s.HeaderTab = lipgloss.NewStyle().
		Foreground(colorOverlay1).
		Padding(0, 1)

	s.HeaderTabSel = lipgloss.NewStyle().
		Foreground(colorMauve).
		Background(colorSurface0).
		Padding(0, 1).
		Bold(true)

	// List
	s.ListItem = lipgloss.NewStyle().
		Padding(0, 1).
		Foreground(colorText)

	s.ListItemSelected = lipgloss.NewStyle().
		Padding(0, 1).
		Background(colorSurface0).
		Foreground(colorText).
		Bold(true)

	s.ListItemTitle = lipgloss.NewStyle().
		Foreground(colorText)

	s.ListItemDesc = lipgloss.NewStyle().
		Foreground(colorSubtext0).
		Italic(true)

	// Kanban
	s.KanbanColumn = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorSurface2).
		Padding(0, 1)

	s.KanbanColumnSelected = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorMauve).
		Padding(0, 1)

	s.KanbanColumnTitle = lipgloss.NewStyle().
		Foreground(colorText).
		Bold(true).
		Padding(0, 0, 1, 0)

	s.KanbanCard = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorSurface1).
		Padding(0, 1).
		Margin(0, 0, 1, 0)

	s.KanbanCardSelected = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorLavender).
		Background(colorSurface0).
		Padding(0, 1).
		Margin(0, 0, 1, 0)

	// Priorities
	s.PriorityLow = lipgloss.NewStyle().
		Foreground(colorGreen)

	s.PriorityMedium = lipgloss.NewStyle().
		Foreground(colorBlue)

	s.PriorityHigh = lipgloss.NewStyle().
		Foreground(colorPeach)

	s.PriorityCritical = lipgloss.NewStyle().
		Foreground(colorRed).
		Bold(true)

	// Statuses
	s.StatusTodo = lipgloss.NewStyle().
		Foreground(colorSubtext0)

	s.StatusInProgress = lipgloss.NewStyle().
		Foreground(colorBlue)

	s.StatusBlocked = lipgloss.NewStyle().
		Foreground(colorRed)

	s.StatusDone = lipgloss.NewStyle().
		Foreground(colorGreen)

	// Tags
	s.Tag = lipgloss.NewStyle().
		Foreground(colorCrust).
		Background(colorMauve).
		Padding(0, 1)

	// Footer
	s.Footer = lipgloss.NewStyle().
		Background(colorMantle).
		Foreground(colorSubtext0).
		Padding(0, 1)

	s.HelpKey = lipgloss.NewStyle().
		Foreground(colorMauve).
		Bold(true)

	s.HelpValue = lipgloss.NewStyle().
		Foreground(colorSubtext0)

	s.HelpSep = lipgloss.NewStyle().
		Foreground(colorSurface2)

	// Form
	s.FormLabel = lipgloss.NewStyle().
		Foreground(colorText).
		Bold(true)

	s.FormInput = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorSurface2).
		Padding(0, 1)

	s.FormInputFocus = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorMauve).
		Padding(0, 1)

	s.FormButton = lipgloss.NewStyle().
		Foreground(colorText).
		Background(colorSurface1).
		Padding(0, 2)

	s.FormButtonFocus = lipgloss.NewStyle().
		Foreground(colorCrust).
		Background(colorMauve).
		Padding(0, 2).
		Bold(true)

	// Help panel
	s.HelpPanel = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorMauve).
		Padding(1, 2)

	s.HelpPanelTitle = lipgloss.NewStyle().
		Foreground(colorMauve).
		Bold(true)

	// Dialog
	s.Dialog = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorMauve).
		Padding(1, 2).
		Background(colorSurface0)

	s.DialogTitle = lipgloss.NewStyle().
		Foreground(colorMauve).
		Bold(true)

	s.Border = lipgloss.RoundedBorder()

	return s
}

// PriorityStyle returns the style for a given priority
func (s Styles) PriorityStyle(p model.Priority) lipgloss.Style {
	switch p {
	case model.PriorityLow:
		return s.PriorityLow
	case model.PriorityMedium:
		return s.PriorityMedium
	case model.PriorityHigh:
		return s.PriorityHigh
	case model.PriorityCritical:
		return s.PriorityCritical
	default:
		return s.PriorityMedium
	}
}

// StatusStyle returns the style for a given status
func (s Styles) StatusStyle(st model.Status) lipgloss.Style {
	switch st {
	case model.StatusTodo:
		return s.StatusTodo
	case model.StatusInProgress:
		return s.StatusInProgress
	case model.StatusBlocked:
		return s.StatusBlocked
	case model.StatusDone:
		return s.StatusDone
	default:
		return s.StatusTodo
	}
}

// PriorityIcon returns an icon for the priority
func PriorityIcon(p model.Priority) string {
	switch p {
	case model.PriorityLow:
		return "○"
	case model.PriorityMedium:
		return "◐"
	case model.PriorityHigh:
		return "●"
	case model.PriorityCritical:
		return "◉"
	default:
		return "○"
	}
}

// StatusIcon returns an icon for the status
func StatusIcon(s model.Status) string {
	switch s {
	case model.StatusTodo:
		return "☐"
	case model.StatusInProgress:
		return "◷"
	case model.StatusBlocked:
		return "⊘"
	case model.StatusDone:
		return "☑"
	default:
		return "☐"
	}
}
