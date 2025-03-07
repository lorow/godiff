package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"strings"
)

// I know I should have used bubbles text input component here
// I just wanted to see how I'd write my own, so I did

type CursorPosition struct {
	X int
}

type TextInputModel struct {
	focused        bool
	prompt         string
	inputText      []rune
	cursor         CursorModel
	cursorPosition CursorPosition
	width          int
}

func NewTextInput() TextInputModel {
	return TextInputModel{
		focused:        false,
		prompt:         ">",
		cursor:         NewCursorModel(),
		inputText:      make([]rune, 0),
		cursorPosition: CursorPosition{X: 1},
		width:          0,
	}
}

func (m TextInputModel) Update(msg tea.Msg) (TextInputModel, tea.Cmd) {
	// todo add keymaps

	// check if we should initialize blinking
	if initMsg, ok := msg.(InitCursorBlinkMsg); ok {
		var cmd tea.Cmd
		m.cursor.Focus()
		m.cursor, cmd = m.cursor.Update(initMsg)
		m.cursor.Blur()
		return m, cmd
	}

	if !m.focused {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyBackspace:
			m.DeleteCharBackwards()
		case tea.KeyLeft:
			if m.cursorPosition.X > 1 {
				m.cursorPosition.X--
			}
		case tea.KeyRight:
			if m.cursorPosition.X < len(m.inputText)+1 {
				m.cursorPosition.X++
			}
		case tea.KeyCtrlV:
			m.insertRuneFromUserInput(msg.Runes)
		default:
			m.insertRuneFromUserInput(msg.Runes)
		}
	}

	var cmds []tea.Cmd
	var cmd tea.Cmd

	m.cursor, cmd = m.cursor.Update(msg)
	cmds = append(cmds, cmd)

	if m.cursor.Mode() == CursorBlink {
		cmds = append(cmds, m.cursor.BlinkCmd())
	}

	return m, tea.Batch(cmds...)
}

func (m TextInputModel) View() string {
	v := strings.Builder{}
	inputTextLength := len(m.inputText)

	// cursor is "in" the text,
	// so we should grab the char under it
	// and display it as cursor
	if m.cursorPosition.X <= inputTextLength {
		m.cursor.SetChar(string(m.inputText[m.cursorPosition.X-1]))
	} else {
		m.cursor.SetChar("█")
	}

	v.WriteString(m.prompt)
	v.WriteString(" ")

	if m.cursorPosition.X > inputTextLength {
		v.WriteString(string(m.inputText))
		v.WriteString(m.cursor.View())
	} else {
		offset := max(0, m.cursorPosition.X-1)
		v.WriteString(string(m.inputText[:offset]))
		v.WriteString(m.cursor.View())
		v.WriteString(string(m.inputText[offset+1:]))
	}
	return v.String()
}

func (m *TextInputModel) DeleteCharBackwards() {
	head := m.inputText[:max(0, m.cursorPosition.X-2)]
	tail := m.inputText[m.cursorPosition.X-1:]

	if m.cursorPosition.X > 1 {
		m.cursorPosition.X--
	}

	m.inputText = append(head, tail...)
}

func (m *TextInputModel) insertRuneFromUserInput(values []rune) {
	head := m.inputText[:m.cursorPosition.X-1]
	tail := m.inputText[m.cursorPosition.X-1:]

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
	m.cursorPosition.X = 1
	m.inputText = make([]rune, 0)
}

func (m *TextInputModel) Focus() {
	m.cursor.Focus()
	m.focused = true
}

func (m TextInputModel) Focused() bool {
	return m.focused
}

func (m *TextInputModel) Blur() {
	m.cursor.Blur()
	m.focused = false
}

func TextInputBlink() tea.Msg {
	return BlinkCursor()
}
