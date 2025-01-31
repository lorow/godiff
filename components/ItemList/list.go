package ItemList

import (
	"errors"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

// Item defines a common interface that represents a given item in the list
type Item interface{}

// ItemRenderer defines a common interface for rendering items in the list
type ItemRenderer interface {
	Render(item Item, model Model, index int) string
	Height() int
	Spacing() int
}

type Model struct {
	title        string
	paddingTop   int
	noItemsText  string
	styles       Styles
	width        int
	height       int
	cursor       int
	items        []Item
	itemRenderer ItemRenderer
}

func New(title, noItemsText string, items []Item, paddingTop int) Model {
	return Model{
		title:        title,
		noItemsText:  noItemsText,
		styles:       DefaultStyles(),
		items:        items,
		cursor:       0,
		paddingTop:   paddingTop,
		itemRenderer: NewDefaultItemRenderer(),
	}
}

func (m Model) View() string {
	var (
		sections        []string
		availableHeight = m.height
	)

	container := m.styles.Container.Width(m.width).Height(m.height)

	title := m.styles.Title.Render(m.title)
	titleRendered := lipgloss.JoinVertical(lipgloss.Top, title, strings.Repeat("\n", m.paddingTop-1))
	availableHeight -= lipgloss.Height(titleRendered)
	sections = append(sections, titleRendered)

	content := lipgloss.NewStyle().Height(availableHeight).Render(m.renderItems(availableHeight))
	sections = append(sections, content)

	return container.Render(lipgloss.JoinVertical(lipgloss.Top, sections...))
}

func (m Model) renderItems(availableHeight int) string {
	var view strings.Builder
	items := m.VisibleItems()
	itemsCount := len(items)
	totalItemHeight := m.itemRenderer.Height() + m.itemRenderer.Spacing()
	maxVisibleItems := max(0, availableHeight/totalItemHeight)

	if itemsCount == 0 {
		return m.styles.NoItems.Render(m.noItemsText)
	}

	for i, item := range items[:min(itemsCount, maxVisibleItems)] {
		// todo add windowing mechanism here
		view.WriteString(m.itemRenderer.Render(item, m, i))
		if i != itemsCount-1 {
			view.WriteString(strings.Repeat("\n", m.itemRenderer.Spacing()+1))
		}
		maxVisibleItems--
	}

	// if we didn't have enough items to fill the view, we need to fill it up with
	// new lines, otherwise drawing anything else might get drawn in the list
	linesToFill := maxVisibleItems * totalItemHeight
	view.WriteString(strings.Repeat("\n", linesToFill))

	return view.String()
}

func (m *Model) CursorUp() {
	if m.cursor > 0 {
		m.cursor--
	}
}

func (m *Model) CursorDown() {
	if m.cursor < len(m.items)-1 {
		m.cursor++
	}
}

func (m *Model) SetItems(items []Item) {
	m.items = items
}

func (m *Model) SetItemRenderer(renderer ItemRenderer) {
	m.itemRenderer = renderer
}

func (m *Model) SetHeight(height int) {
	m.height = height
}

func (m *Model) SetWidth(width int) {
	m.width = width
}

func (m Model) GetCurrentSelection() (Item, error) {
	if m.cursor >= 0 {
		return m.items[m.cursor], nil
	}
	return nil, errors.New("no item selected")
}

func (m Model) VisibleItems() []Item {
	// prepared for adding support for filtering
	return m.items
}

func (m Model) GetIndex() int {
	return m.cursor
}
