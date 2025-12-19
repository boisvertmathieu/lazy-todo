package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// HelpPanel displays keyboard shortcuts and help
type HelpPanel struct {
	styles Styles
	width  int
	height int
}

// NewHelpPanel creates a new help panel
func NewHelpPanel(styles Styles) *HelpPanel {
	return &HelpPanel{
		styles: styles,
	}
}

// SetSize sets the panel dimensions
func (h *HelpPanel) SetSize(width, height int) {
	h.width = width
	h.height = height
}

// Render renders the help panel
func (h *HelpPanel) Render() string {
	title := h.styles.HelpPanelTitle.Render("Raccourcis Clavier")

	sections := []struct {
		title string
		items []struct {
			key  string
			desc string
		}
	}{
		{
			title: "Navigation",
			items: []struct {
				key  string
				desc string
			}{
				{"j / ↓", "Descendre"},
				{"k / ↑", "Monter"},
				{"h / ←", "Gauche (kanban)"},
				{"l / →", "Droite (kanban)"},
			},
		},
		{
			title: "Actions",
			items: []struct {
				key  string
				desc string
			}{
				{"a", "Ajouter une tâche"},
				{"e", "Éditer la tâche"},
				{"d", "Supprimer la tâche"},
				{"p", "Changer la priorité"},
				{"t", "Gérer les tags"},
				{"Enter", "Voir/Éditer détails"},
			},
		},
		{
			title: "États rapides",
			items: []struct {
				key  string
				desc string
			}{
				{"1", "À faire"},
				{"2", "En cours"},
				{"3", "Bloqué"},
				{"4", "Terminé"},
			},
		},
		{
			title: "Kanban",
			items: []struct {
				key  string
				desc string
			}{
				{"H / Shift+←", "Déplacer tâche à gauche"},
				{"L / Shift+→", "Déplacer tâche à droite"},
			},
		},
		{
			title: "Général",
			items: []struct {
				key  string
				desc string
			}{
				{"Tab", "Changer de vue"},
				{"g", "Changer le groupage"},
				{"/", "Rechercher"},
				{"o", "Ouvrir le fichier YAML"},
				{"r", "Rafraîchir"},
				{"?", "Afficher/Masquer l'aide"},
				{"q / Ctrl+C", "Quitter"},
			},
		},
		{
			title: "Formulaire",
			items: []struct {
				key  string
				desc string
			}{
				{"Tab", "Champ suivant"},
				{"Shift+Tab", "Champ précédent"},
				{"Enter", "Valider"},
				{"Esc", "Annuler"},
			},
		},
	}

	keyStyle := h.styles.HelpKey
	descStyle := h.styles.HelpValue
	sectionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#cba6f7")).
		Bold(true).
		MarginTop(1)

	var content []string
	content = append(content, title)
	content = append(content, "")

	for _, section := range sections {
		content = append(content, sectionStyle.Render(section.title))
		for _, item := range section.items {
			line := keyStyle.Render(padRight(item.key, 16)) + descStyle.Render(item.desc)
			content = append(content, line)
		}
	}

	panelContent := strings.Join(content, "\n")

	return h.styles.HelpPanel.
		Width(h.width - 4).
		Height(h.height - 4).
		Render(panelContent)
}

// padRight pads a string to the right
func padRight(s string, length int) string {
	if len(s) >= length {
		return s
	}
	return s + strings.Repeat(" ", length-len(s))
}

// RenderFooter renders the footer help bar
func RenderFooter(styles Styles, isKanban bool) string {
	var items []string

	addItem := func(key, desc string) {
		items = append(items, styles.HelpKey.Render(key)+styles.HelpSep.Render(":")+styles.HelpValue.Render(desc))
	}

	addItem("j/k", "nav")
	if isKanban {
		addItem("h/l", "colonnes")
		addItem("H/L", "déplacer")
	}
	addItem("a", "ajouter")
	addItem("d", "supprimer")
	addItem("1-4", "état")
	addItem("g", "grouper")
	addItem("Tab", "vue")
	addItem("?", "aide")
	addItem("q", "quitter")

	separator := styles.HelpSep.Render(" │ ")
	return styles.Footer.Render(strings.Join(items, separator))
}
