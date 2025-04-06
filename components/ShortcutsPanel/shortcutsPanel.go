package ShortcutsPanel

import "github.com/charmbracelet/lipgloss"

type Shortcut struct {
	Key         string
	Description string
}
type Model struct {
	width     int
	styles    Styles
	shortcuts []Shortcut
}

func New(ops ...func(*Model)) *Model {
	model := &Model{
		styles:    DefaultStyles(),
		shortcuts: []Shortcut{},
	}

	for _, op := range ops {
		op(model)
	}

	return model
}

func WithStyles(styles Styles) func(model *Model) {
	return func(model *Model) {
		model.styles = styles
	}
}

func WithShortcuts(shortcuts []Shortcut) func(model *Model) {
	return func(model *Model) {
		model.shortcuts = shortcuts
	}
}

func WithWidth(width int) func(model *Model) {
	return func(model *Model) {
		model.width = width
	}
}

func (m *Model) SetWidth(width int) {
	m.width = width
}

func (m Model) GetWidth() int {
	return m.width
}

func (m Model) View() string {
	renderedShortcuts := []string{}
	for _, shortcut := range m.shortcuts {
		keyRendered := m.styles.ShortcutKey.Render(shortcut.Key)
		descriptionRendered := m.styles.ShortcutDescription.Render(shortcut.Description)
		shortcutRendered := m.styles.ShortcutContainer.Render(lipgloss.JoinHorizontal(lipgloss.Top, keyRendered, " - ", descriptionRendered))
		renderedShortcuts = append(renderedShortcuts, shortcutRendered)
	}

	shortcutsBar := m.styles.Container.Render(lipgloss.JoinHorizontal(lipgloss.Top, renderedShortcuts...))
	filler := m.styles.Container.Width(m.width - lipgloss.Width(shortcutsBar)).Render(" ")

	return lipgloss.JoinHorizontal(lipgloss.Top, shortcutsBar, filler)
}

func (m *Model) SetShortcuts(shortcuts []Shortcut) {
	m.shortcuts = shortcuts
}
