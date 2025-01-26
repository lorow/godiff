package ItemList

import "github.com/charmbracelet/lipgloss"

type DefaultItemStyles struct {
	NormalTitle         lipgloss.Style
	NormalDescription   lipgloss.Style
	SelectedTitle       lipgloss.Style
	SelectedDescription lipgloss.Style
}

func NewDefaultItemStyles() (s DefaultItemStyles) {
	s.NormalTitle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"}).
		Padding(0, 0, 0, 2) //nolint:mnd

	s.NormalDescription = s.NormalTitle.
		Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"})

	s.SelectedTitle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}).
		Foreground(lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"}).
		Padding(0, 0, 0, 1)

	s.SelectedDescription = s.SelectedTitle.
		Foreground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"})

	return s
}

// DefaultItem describes an item that can be rendered by the default item renderer
// it expands on the Item type and requests to provide Title() and Description()
type DefaultItem interface {
	Item
	Title() string
	Description() string
}

type DefaultItemRenderer struct {
	Styles  DefaultItemStyles
	height  int
	spacing int
}

func NewDefaultItemRenderer() DefaultItemRenderer {
	return DefaultItemRenderer{}
}
