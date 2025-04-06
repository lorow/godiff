package ShortcutsPanel

import "github.com/charmbracelet/lipgloss"

type Model struct {
	width     int
	styles    Styles
	shortcuts []string
}

func New(ops ...func(*Model)) *Model {
	model := &Model{
		styles:    DefaultStyles(),
		shortcuts: []string{},
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

func WithShortcuts(shortcuts []string) func(model *Model) {
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
		renderedShortcuts = append(renderedShortcuts, m.styles.Shortcut.Render(shortcut))
	}

	shortcutsBar := m.styles.Container.Render(lipgloss.JoinHorizontal(lipgloss.Top, renderedShortcuts...))
	filler := m.styles.Container.Background(lipgloss.Color("#E06C75")).Width(m.width - lipgloss.Width(shortcutsBar)).Render(" ")

	return lipgloss.JoinHorizontal(lipgloss.Top, shortcutsBar, filler)
}

func (m *Model) SetShortcuts(shortcuts []string) {
	m.shortcuts = shortcuts
}
