package ItemList

import (
	"errors"
	"godiff/components/FocusChain"
	"godiff/components/scrollablePaginator"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
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
	renderedTitle   string
	noItemsText     string
	styles          Styles
	focusedStyles   Styles
	paddingTop      int
	width           int
	height          int
	availableHeight int
	cursor          int
	amountOfItems   int
	items           []Item
	paginator       scrollablePaginator.Model
	itemRenderer    ItemRenderer
	isFocused       bool
	onSelection     func(Item) tea.Cmd
}

func New(opts ...func(*Model)) *Model {
	model := Model{
		noItemsText:     "No items",
		styles:          DefaultStyles(),
		focusedStyles:   FocusedStyles(),
		paddingTop:      1,
		amountOfItems:   0,
		width:           0,
		height:          0,
		availableHeight: 0,
		cursor:          0,
		itemRenderer:    NewDefaultItemRenderer(),
		paginator:       scrollablePaginator.New(),
		isFocused:       false,
	}

	for _, opt := range opts {
		opt(&model)
	}

	// we need to recalculate some parts for items to line up correctly
	model.SetTitle(model.title)
	model.SetPaddingTop(model.paddingTop)

	return &model
}

func WithItems(items []Item) func(*Model) {
	return func(model *Model) {
		model.SetItems(items)
	}
}

func WithItemRenderer(itemRenderer ItemRenderer) func(*Model) {
	return func(model *Model) {
		model.SetItemRenderer(itemRenderer)
	}
}

func WithOnSelection(onSelection func(Item) tea.Cmd) func(*Model) {
	return func(model *Model) {
		model.onSelection = onSelection
	}
}

func WithTittle(title string) func(*Model) {
	return func(model *Model) {
		model.SetTitle(title)
	}
}

func WithNoItemsText(noItemsText string) func(*Model) {
	return func(model *Model) {
		model.noItemsText = noItemsText
	}
}

func WithPaddingTop(padding int) func(model *Model) {
	if padding <= 0 {
		padding = 1
	}

	return func(model *Model) {
		model.SetPaddingTop(padding)
	}
}

func (m Model) View() string {
	var sections []string
	styles := m.getStyles()
	container := styles.Container.Width(m.width).Height(m.height)
	sections = append(sections, m.renderedTitle)

	content := lipgloss.NewStyle().Height(m.getAvailableHeight()).Render(m.renderItems())
	sections = append(sections, content)

	return container.Render(lipgloss.JoinVertical(lipgloss.Top, sections...))
}

func (m Model) renderItems() string {
	var view strings.Builder
	styles := m.getStyles()
	availableHeight := m.getAvailableHeight()
	totalItemHeight := m.getItemHeight()

	visibleItems := m.VisibleItems()
	itemsCount := len(visibleItems)

	maxVisibleItems := max(0, availableHeight/totalItemHeight)

	if itemsCount == 0 {
		return styles.NoItems.Render(m.noItemsText)
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

func (m Model) getStyles() Styles {
	if m.isFocused {
		return m.focusedStyles
	}
	return m.styles
}

func (m *Model) Focus() {
	m.isFocused = true
}

func (m *Model) Blur() {
	m.isFocused = false
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyDown:
			if m.IsCursorAtTheBottom() {
				return FocusChain.SwitchFocusCmd(FocusChain.FocusDown)
			}
			m.CursorDown()
		case tea.KeyUp:
			if m.IsCursorAtTheTop() {
				return FocusChain.SwitchFocusCmd(FocusChain.FocusUp)
			}
			m.CursorUp()
		case tea.KeyEnter:
			if m.onSelection == nil {
				return nil
			}

			currentSelection, _ := m.GetCurrentSelection()
			return m.onSelection(currentSelection)
		}
	}

	return nil
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
	styles := m.getStyles()
	m.title = title
	renderedTItle := styles.Title.Render(title)
	m.renderedTitle = lipgloss.JoinVertical(lipgloss.Top, renderedTItle, strings.Repeat("\n", m.paddingTop-1))
	m.recalculateAvailableHeight()
	m.recalculateViewBounds()
}

func (m *Model) SetPaddingTop(paddingTop int) {
	m.renderedTitle = lipgloss.JoinVertical(lipgloss.Top, m.renderedTitle, strings.Repeat("\n", paddingTop-1))
	m.recalculateAvailableHeight()
	m.recalculateViewBounds()
}

func (m *Model) recalculateAvailableHeight() {
	m.availableHeight = m.height - lipgloss.Height(m.renderedTitle)
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

func (m Model) IsCursorAtTheTop() bool {
	return m.cursor == 0
}

func (m Model) IsCursorAtTheBottom() bool {
	return m.cursor == m.amountOfItems-1
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
