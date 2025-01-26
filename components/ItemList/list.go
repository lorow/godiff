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
	Width() int
	Spacing() int
}

type State int

const (
	Default State = iota
	Loading
	Loaded
	Error
)

type Model struct {
	title        string
	noItemsText  string
	styles       Styles
	width        int
	height       int
	cursor       int
	state        State
	items        []Item
	itemRenderer ItemRenderer
}

func New(title, noItemsText string) Model {
	return Model{
		title:       title,
		noItemsText: noItemsText,
		styles:      DefaultStyles(),
		items:       []Item{},
		state:       Default,
		cursor:      -1,
	}
}

func (m Model) View() string {
	var (
		sections        []string
		availableHeight = m.height
	)

	container := m.styles.Container.Width(m.width).Height(m.height)

	titleRendered := m.styles.Title.Render(m.title)
	availableHeight -= lipgloss.Height(titleRendered)
	sections = append(sections, titleRendered)

	content := lipgloss.NewStyle().Height(availableHeight).Render(m.renderItems())
	sections = append(sections, content)

	return container.Render(lipgloss.JoinVertical(lipgloss.Top, sections...))
}

func (m Model) renderItems() string {
	items := m.VisibleItems()
	itemsCount := len(items)

	var view strings.Builder

	if itemsCount == 0 {
		return m.styles.NoItems.Render(m.noItemsText)
	}

	for i, item := range items {
		view.WriteString(m.itemRenderer.Render(item, m, i))
		if i != itemsCount-1 {
			view.WriteString(strings.Repeat("\n", m.itemRenderer.Spacing()+1))
		}
	}

	// todo see if this will break if we don't have enough items

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

func (m *Model) SetState(state State) {
	m.state = state
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
