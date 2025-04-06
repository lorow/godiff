package ShortcutsPanel

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	Container           lipgloss.Style
	ShortcutContainer   lipgloss.Style
	ShortcutKey         lipgloss.Style
	ShortcutDescription lipgloss.Style
}

func DefaultStyles() (s Styles) {
	backgroundColor := lipgloss.Color("#282C34")
	s.Container = lipgloss.NewStyle().Background(backgroundColor)
	s.ShortcutContainer = lipgloss.NewStyle().Background(backgroundColor).PaddingLeft(2)
	s.ShortcutKey = lipgloss.NewStyle().Background(backgroundColor).Foreground(lipgloss.Color("#61AFEF"))
	s.ShortcutDescription = lipgloss.NewStyle().Background(backgroundColor)

	return s
}
