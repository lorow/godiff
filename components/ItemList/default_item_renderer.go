package ItemList

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

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
// due to List requiring render() to take an item, this interface is mostly for validation
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

func (d DefaultItemRenderer) Render(item Item, model Model, index int) string {

	// if for some reason the list has zero width, we don't have to anything
	if model.width <= 0 {
		return ""
	}

	defaultItem, isDefaultItem := item.(DefaultItem)
	if !isDefaultItem {
		return ""
	}

	title := defaultItem.Title()
	description := defaultItem.Description()

	// we also have to make sure the text doesn't overflow the list
	textWidth := model.width - d.Styles.NormalTitle.GetPaddingLeft() - d.Styles.NormalTitle.GetPaddingRight()
	title = ansi.Truncate(title, textWidth, "...")
	description = ansi.Truncate(description, textWidth, "...")

	isSelected := index == model.GetIndex()

	if isSelected {
		title = d.Styles.SelectedTitle.Render(title)
		description = d.Styles.SelectedDescription.Render(description)
		return lipgloss.JoinVertical(lipgloss.Top, title, description)
	}

	title = d.Styles.NormalTitle.Render(title)
	description = d.Styles.NormalDescription.Render(description)
	return lipgloss.JoinVertical(lipgloss.Top, title, description)
}

func (d DefaultItemRenderer) Height() int {
	return d.height
}

func (d DefaultItemRenderer) Spacing() int {
	return d.spacing
}

func NewDefaultItemRenderer() DefaultItemRenderer {
	return DefaultItemRenderer{
		// update this when the widget height changes
		Styles:  NewDefaultItemStyles(),
		height:  2,
		spacing: 1,
	}
}
