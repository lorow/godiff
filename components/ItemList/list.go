package ItemList

import (
	"errors"
	"godiff/components/scrollablePaginator"
	"strings"

	"github.com/charmbracelet/lipgloss"
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
	title           string
	noItemsText     string
	styles          Styles
	paddingTop      int
	width           int
	height          int
	availableHeight int
	cursor          int
	amountOfItems   int
	items           []Item
	paginator       scrollablePaginator.Model
	itemRenderer    ItemRenderer
}

func New(title, noItemsText string, items []Item, itemRenderer ItemRenderer, paddingTop int) (Model, error) {
	if paddingTop <= 0 {
		return Model{}, errors.New("padding must be greater than 0")
	}

	styles := DefaultStyles()
	model := Model{
		noItemsText:     noItemsText,
		styles:          styles,
		paddingTop:      paddingTop,
		amountOfItems:   len(items),
		items:           items,
		width:           0,
		height:          0,
		availableHeight: 0,
		cursor:          0,
		paginator:       scrollablePaginator.New(),
		itemRenderer:    itemRenderer,
	}

	model.SetTitle(title)
	model.SetPaddingTop(paddingTop)

	return model, nil
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
		itemIndex := m.paginator.GetLowerBound() + i
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
		if m.cursor < m.paginator.GetLowerBound() {
			m.paginator.MoveViewWindowUp()
		}
	}
}

func (m *Model) CursorDown() {
	ceiling := m.amountOfItems - 1
	if m.cursor < ceiling {
		m.cursor++

		if m.cursor == m.paginator.GetUpperBound() && m.cursor <= ceiling {
			m.paginator.MoveViewWindowDown()
		}
	}
}

func (m *Model) SetItems(items []Item) {
	m.amountOfItems = len(items)
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
	renderedTItle := m.styles.Title.Render(title)
	m.title = lipgloss.JoinVertical(lipgloss.Top, renderedTItle, strings.Repeat("\n", m.paddingTop-1))
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
	m.paginator.RecalculateViewBounds(m.getAvailableHeight(), m.getItemHeight())
}

func (m Model) getAvailableHeight() int {
	return m.availableHeight
}

func (m Model) getItemHeight() int {
	return m.itemRenderer.Height() + m.itemRenderer.Spacing()
}

func (m Model) GetItemsCount() int {
	return m.amountOfItems
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
	viewBounds := m.paginator.GetViewWindowBounds()
	return m.items[viewBounds[0]:min(viewBounds[1], m.amountOfItems)]
}

func (m Model) GetIndex() int {
	return m.cursor
}
