package ShortcutsPanel

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Shortcut struct {
	Key          string
	Description  string
	comboToCheck string
	command      tea.Cmd
}

type Model struct {
	width     int
	styles    Styles
	shortcuts []Shortcut
}

func NewShortcut(key, comboToCheck, description string, command tea.Cmd) Shortcut {
	return Shortcut{
		Key:          key,
		Description:  description,
		comboToCheck: comboToCheck,
		command:      command,
	}
}

func (s *Shortcut) HasComboHit(combo string) bool {
	return s.comboToCheck == combo
}

func (s *Shortcut) GetCommand() tea.Cmd {
	return s.command
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
		shortcutRendered := m.styles.ShortcutContainer.Render(lipgloss.JoinHorizontal(lipgloss.Top, keyRendered, m.styles.Container.Render(" - "), descriptionRendered))
		renderedShortcuts = append(renderedShortcuts, shortcutRendered)
	}

	shortcutsBar := m.styles.Container.Render(lipgloss.JoinHorizontal(lipgloss.Top, renderedShortcuts...))
	filler := m.styles.Container.Width(m.width - lipgloss.Width(shortcutsBar)).Render(" ")

	return lipgloss.JoinHorizontal(lipgloss.Top, shortcutsBar, filler)
}

func (m *Model) SetShortcuts(shortcuts []Shortcut) {
	m.shortcuts = shortcuts
}

func (m *Model) GetActiveShortcuts() []Shortcut {
	return m.shortcuts
}

func (m *Model) CheckIfShortcutHit(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		for _, shortcut := range m.GetActiveShortcuts() {
			if shortcut.HasComboHit(msg.String()) {
				return shortcut.GetCommand()
			}
		}
	}

	return nil
}
