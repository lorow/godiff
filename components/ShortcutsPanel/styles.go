package ShortcutsPanel

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	Container lipgloss.Style
	Shortcut  lipgloss.Style
}

func DefaultStyles() (s Styles) {
	s.Container = lipgloss.NewStyle()
	s.Shortcut = lipgloss.NewStyle().PaddingLeft(5)

	return s
}
