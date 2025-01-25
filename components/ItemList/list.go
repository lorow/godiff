package ItemList

import (
	"errors"
	tea "github.com/charmbracelet/bubbletea"
)

// Item defines a common interface that represents a given item in the list
type Item interface{}

// ItemRenderer defines a common interface for rendering items in the list
type ItemRenderer interface {
	Render(item Item, model Model, index int) string
	Height() int
	Width() int
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
	styles       Styles
	width        int
	height       int
	cursor       int
	state        State
	items        []Item
	itemRenderer ItemRenderer
}

func New() Model {
	return Model{
		items:  []Item{},
		state:  Default,
		cursor: -1,
	}
}

// todo render the list with itemRenderer
func (m Model) View() string {
	return ""
}

// todo think if this should be handled by the component or by the one using it
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyDown:
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
			return m, nil
		case tea.KeyUp:
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil
		default:
			return m, nil
		}
	}

	return m, nil
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

func (m Model) GetCurrentSelection() (Item, error) {
	if m.cursor >= 0 {
		return m.items[m.cursor], nil
	}
	return nil, errors.New("no item selected")
}
