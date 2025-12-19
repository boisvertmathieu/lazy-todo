package keys

import "github.com/charmbracelet/bubbles/key"

// KeyMap contains all keybindings for the application
type KeyMap struct {
	// Navigation
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding

	// Actions
	Add       key.Binding
	Edit      key.Binding
	Delete    key.Binding
	Enter     key.Binding
	Priority  key.Binding
	Tag       key.Binding
	MoveLeft  key.Binding
	MoveRight key.Binding

	// Quick status change
	StatusTodo       key.Binding
	StatusInProgress key.Binding
	StatusBlocked    key.Binding
	StatusDone       key.Binding

	// Views
	ToggleView key.Binding
	Search     key.Binding
	OpenEditor key.Binding
	Help       key.Binding
	Refresh    key.Binding

	// Form
	Submit key.Binding
	Cancel key.Binding
	Next   key.Binding
	Prev   key.Binding

	// Global
	Quit key.Binding
}

// DefaultKeyMap returns the default keybindings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		// Navigation
		Up: key.NewBinding(
			key.WithKeys("k", "up"),
			key.WithHelp("k/↑", "monter"),
		),
		Down: key.NewBinding(
			key.WithKeys("j", "down"),
			key.WithHelp("j/↓", "descendre"),
		),
		Left: key.NewBinding(
			key.WithKeys("h", "left"),
			key.WithHelp("h/←", "gauche"),
		),
		Right: key.NewBinding(
			key.WithKeys("l", "right"),
			key.WithHelp("l/→", "droite"),
		),

		// Actions
		Add: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "ajouter"),
		),
		Edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "éditer"),
		),
		Delete: key.NewBinding(
			key.WithKeys("d", "delete"),
			key.WithHelp("d", "supprimer"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "sélectionner"),
		),
		Priority: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "priorité"),
		),
		Tag: key.NewBinding(
			key.WithKeys("t"),
			key.WithHelp("t", "tag"),
		),
		MoveLeft: key.NewBinding(
			key.WithKeys("H", "shift+left"),
			key.WithHelp("H", "déplacer ←"),
		),
		MoveRight: key.NewBinding(
			key.WithKeys("L", "shift+right"),
			key.WithHelp("L", "déplacer →"),
		),

		// Quick status
		StatusTodo: key.NewBinding(
			key.WithKeys("1"),
			key.WithHelp("1", "à faire"),
		),
		StatusInProgress: key.NewBinding(
			key.WithKeys("2"),
			key.WithHelp("2", "en cours"),
		),
		StatusBlocked: key.NewBinding(
			key.WithKeys("3"),
			key.WithHelp("3", "bloqué"),
		),
		StatusDone: key.NewBinding(
			key.WithKeys("4"),
			key.WithHelp("4", "terminé"),
		),

		// Views
		ToggleView: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "changer vue"),
		),
		Search: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "rechercher"),
		),
		OpenEditor: key.NewBinding(
			key.WithKeys("o"),
			key.WithHelp("o", "ouvrir fichier"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "aide"),
		),
		Refresh: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "rafraîchir"),
		),

		// Form
		Submit: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "valider"),
		),
		Cancel: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "annuler"),
		),
		Next: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "suivant"),
		),
		Prev: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "précédent"),
		),

		// Global
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quitter"),
		),
	}
}

// ShortHelp returns a short help string
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Up, k.Down, k.Add, k.Delete, k.ToggleView, k.Help, k.Quit,
	}
}

// FullHelp returns the full help for all keybindings
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Add, k.Edit, k.Delete, k.Priority},
		{k.StatusTodo, k.StatusInProgress, k.StatusBlocked, k.StatusDone},
		{k.ToggleView, k.Search, k.OpenEditor, k.Help},
		{k.MoveLeft, k.MoveRight, k.Refresh, k.Quit},
	}
}
