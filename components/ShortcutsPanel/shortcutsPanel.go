package ShortcutsPanel

import "github.com/charmbracelet/lipgloss"

type Model struct {
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

func (m Model) View() string {

	// generally, we might want to implement wordbreak
	// buut this might be a bit overkill for now

	renderedShortcuts := []string{}
	for _, shortcut := range m.shortcuts {
		renderedShortcuts = append(renderedShortcuts, m.styles.Shortcut.Render(shortcut))
	}

	return m.styles.Container.Render(lipgloss.JoinHorizontal(lipgloss.Top, renderedShortcuts...))
}

func (m *Model) SetShortcuts(shortcuts []string) {
	m.shortcuts = shortcuts
}
