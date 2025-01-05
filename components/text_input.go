package components

import tea "github.com/charmbracelet/bubbletea"

// I know I should have used bubbles text input component here
// I just wanted to see how I'd write my own, so I did

type CursorPosition struct {
	X int
}

type TextInputModel struct {
	focused        bool
	prompt         string
	inputText      []rune
	cursorPosition CursorPosition
	width          int
}

func NewTextInput() TextInputModel {
	return TextInputModel{
		focused:        false,
		prompt:         ">",
		inputText:      make([]rune, 0),
		cursorPosition: CursorPosition{X: 0},
		width:          0,
	}
}

func (m TextInputModel) Update(msg tea.Msg) (TextInputModel, tea.Cmd) {
	// todo add keymaps

	if !m.focused {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyBackspace:
			m.DeleteCharBackwards()
			return m, nil
		case tea.KeyLeft:
			if m.cursorPosition.X > 0 {
				m.cursorPosition.X--
			}
		case tea.KeyRight:
			if m.cursorPosition.X < len(m.inputText) {
				m.cursorPosition.X++
			}
		case tea.KeyCtrlV:
			m.insertRuneFromUserInput(msg.Runes)
		default:
			m.insertRuneFromUserInput(msg.Runes)
			return m, nil
		}
	}

	return m, nil
}

func (m TextInputModel) View() string {
	return m.prompt + " " + string(m.inputText)
}

func (m *TextInputModel) DeleteCharBackwards() {
	head := m.inputText[:max(0, m.cursorPosition.X-1)]
	tail := m.inputText[m.cursorPosition.X:]

	if m.cursorPosition.X > 0 {
		m.cursorPosition.X--
	}

	m.inputText = append(head, tail...)
}

func (m *TextInputModel) insertRuneFromUserInput(values []rune) {
	head := m.inputText[:m.cursorPosition.X]
	tail := m.inputText[m.cursorPosition.X:]

	for _, value := range values {
		head = append(head, value)
	}

	m.cursorPosition.X = m.cursorPosition.X + len(values)
	m.inputText = append(head, tail...)
}

func (m *TextInputModel) setWidth(width int) {
	m.width = width
}

func (m *TextInputModel) Reset() {
	m.cursorPosition.X = 0
	m.inputText = make([]rune, 0)
}

func (m *TextInputModel) Focus() {
	m.focused = true
}

func (m TextInputModel) Focused() bool {
	return m.focused
}

func (m *TextInputModel) Blur() {
	m.focused = false
}
