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
	title            string
	noItemsText      string
	styles           Styles
	paddingTop       int
	width            int
	height           int
	availableHeight  int
	cursor           int
	items            []Item
	viewWindowBounds [2]int
	itemRenderer     ItemRenderer
}

func New(title, noItemsText string, items []Item, paddingTop int) Model {
	styles := DefaultStyles()
	model := Model{
		noItemsText:      noItemsText,
		styles:           styles,
		items:            items,
		width:            0,
		height:           0,
		availableHeight:  0,
		cursor:           0,
		viewWindowBounds: [2]int{0, 0},
		itemRenderer:     NewDefaultItemRenderer(),
	}

	model.SetTitle(title)
	model.SetPaddingTop(paddingTop)

	return model
}

func (m Model) View() string {
	var sections []string

	container := m.styles.Container.Width(m.width).Height(m.height)
	sections = append(sections, m.title)

	content := lipgloss.NewStyle().Height(m.getAvailableHeight()).Render(m.renderItems())
	sections = append(sections, content)

	return container.Render(lipgloss.JoinVertical(lipgloss.Top, sections...))
}

func (m Model) renderItems() string {
	var view strings.Builder
	availableHeight := m.getAvailableHeight()
	totalItemHeight := m.getItemHeight()

	visibleItems := m.VisibleItems()
	itemsCount := len(visibleItems)

	maxVisibleItems := max(0, availableHeight/totalItemHeight)

	if itemsCount == 0 {
		return m.styles.NoItems.Render(m.noItemsText)
	}

	for i, item := range visibleItems[:min(itemsCount, maxVisibleItems)] {
		// since we've moved the view window, we need to update the index as well
		itemIndex := m.viewWindowBounds[0] + i
		view.WriteString(m.itemRenderer.Render(item, m, itemIndex))
		if i != itemsCount-1 {
			view.WriteString(strings.Repeat("\n", m.itemRenderer.Spacing()+1))
		}
		maxVisibleItems--
	}

	// if we didn't have enough visibleItems to fill the view, we need to fill it up with
	// new lines, otherwise drawing anything else might get drawn in the list
	linesToFill := maxVisibleItems * totalItemHeight
	view.WriteString(strings.Repeat("\n", linesToFill))

	return view.String()
}

func (m *Model) CursorUp() {
	if m.cursor > 0 {
		m.cursor--
		// we're at the bound of the view window, but not yet at the top
		if m.cursor < m.viewWindowBounds[0] {
			m.viewWindowBounds[0]--
			m.viewWindowBounds[1]--
		}
	}
}

func (m *Model) CursorDown() {
	ceiling := len(m.items) - 1
	if m.cursor < ceiling {
		m.cursor++

		if m.cursor == m.viewWindowBounds[1] && m.cursor <= ceiling {
			m.viewWindowBounds[0]++
			m.viewWindowBounds[1]++
		}
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
	m.recalculateAvailableHeight()
	m.recalculateViewBounds()
}

func (m *Model) SetWidth(width int) {
	m.width = width
}

func (m *Model) SetTitle(title string) {
	m.title = m.styles.Title.Render(title)
	m.recalculateAvailableHeight()
	m.recalculateViewBounds()
}

func (m *Model) SetPaddingTop(paddingTop int) {
	m.title = lipgloss.JoinVertical(lipgloss.Top, m.title, strings.Repeat("\n", paddingTop-1))
	m.recalculateAvailableHeight()
	m.recalculateViewBounds()
}

func (m *Model) recalculateAvailableHeight() {
	m.availableHeight = m.height - lipgloss.Height(m.title)
}

func (m *Model) recalculateViewBounds() {
	// we need to update the max visible area

	distanceBetweenCurrentBounds := m.viewWindowBounds[1] - m.viewWindowBounds[0]
	newMaxDistance := m.getAvailableHeight() / m.getItemHeight()
	if m.viewWindowBounds[0] == 0 && m.viewWindowBounds[1] == 0 {
		// we're just starting then, we only need the upper bound
		m.viewWindowBounds[1] = max(0, newMaxDistance)
		return
	}
	// otherwise we have to expand the list downwards, so
	expandViewBy := newMaxDistance - distanceBetweenCurrentBounds
	m.viewWindowBounds[1] += expandViewBy

}

func (m Model) getAvailableHeight() int {
	return m.availableHeight
}

func (m Model) getItemHeight() int {
	return m.itemRenderer.Height() + m.itemRenderer.Spacing()
}

func (m Model) GetCurrentSelection() (Item, error) {
	if m.cursor >= 0 {
		return m.items[m.cursor], nil
	}
	return nil, errors.New("no item selected")
}

func (m Model) VisibleItems() []Item {
	// visible items should return only the items that fit the current view window
	// or, if we're filtering - filtered items that would fit in the very same window
	return m.items[m.viewWindowBounds[0]:min(m.viewWindowBounds[1], len(m.items))]
}

func (m Model) GetIndex() int {
	return m.cursor
}
