package ShortcutsPanel

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	Container           lipgloss.Style
	ShortcutContainer   lipgloss.Style
	ShortcutKey         lipgloss.Style
	ShortcutDescription lipgloss.Style
}

func DefaultStyles() (s Styles) {
	s.Container = lipgloss.NewStyle()
	s.ShortcutContainer = lipgloss.NewStyle().PaddingLeft(2)
	s.ShortcutKey = lipgloss.NewStyle().Foreground(lipgloss.Color("#61AFEF"))
	s.ShortcutDescription = lipgloss.NewStyle()

	return s
}
