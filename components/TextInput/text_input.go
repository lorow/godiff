package TextInput

import (
	"godiff/components/Cursor"
	"godiff/components/FocusChain"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type CursorPosition struct {
	X int
}

type Model struct {
	prompt    string
	inputText []rune
	// cursor         Cursor.Model
	cursorPosition CursorPosition
	width          int
	isFocused      bool
}

func New(opts ...func(*Model)) *Model {
	model := &Model{
		prompt: ">",
		// cursor:         Cursor.NewCursorModel(),
		inputText:      make([]rune, 0),
		cursorPosition: CursorPosition{X: 1},
		width:          0,
		isFocused:      false,
	}

	for _, opt := range opts {
		opt(model)
	}

	return model
}

func WithPrompt(prompt string) func(model *Model) {
	return func(model *Model) {
		model.prompt = prompt
	}
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	// todo add keymaps

	// check if we should initialize blinking
	// if initMsg, ok := msg.(Cursor.InitCursorBlinkMsg); ok {
	// 	var cmd tea.Cmds
	// 	// rethink cursor blinking
	// 	// m.cursor.Focus()
	// 	// m.cursor, cmd = m.cursor.Update(initMsg)
	// 	// m.cursor.Blur()
	// 	return cmd
	// }

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyDown:
			return FocusChain.SwitchFocusCmd(FocusChain.FocusDown)
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

	// var cmds []tea.Cmd
	// var cmd tea.Cmd

	// m.cursor, cmd = m.cursor.Update(msg)
	// cmds = append(cmds, cmd)

	// if m.cursor.Mode() == Cursor.CursorBlink {
	// 	cmds = append(cmds, m.cursor.BlinkCmd())
	// }

	// return tea.Batch(cmds...)
	return nil
}

func (m Model) View() string {
	v := strings.Builder{}
	inputTextLength := len(m.inputText)

	// cursor is "in" the text,
	// so we should grab the char under it
	// and display it as cursor
	if m.cursorPosition.X <= inputTextLength {
		// m.cursor.SetChar(string(m.inputText[m.cursorPosition.X-1]))
	} else {
		// m.cursor.SetChar("█")
	}

	v.WriteString(m.prompt)
	v.WriteString(" ")

	if m.isFocused {
		v.WriteString("Focused")
	}

	if m.cursorPosition.X > inputTextLength {
		v.WriteString(string(m.inputText))
		// v.WriteString(m.cursor.View())
	} else {
		offset := max(0, m.cursorPosition.X-1)
		v.WriteString(string(m.inputText[:offset]))
		// v.WriteString(m.cursor.View())
		v.WriteString(string(m.inputText[offset+1:]))
	}
	return v.String()
}

func (m *Model) DeleteCharBackwards() {
	head := m.inputText[:max(0, m.cursorPosition.X-2)]
	tail := m.inputText[m.cursorPosition.X-1:]

	if m.cursorPosition.X > 1 {
		m.cursorPosition.X--
	}

	m.inputText = append(head, tail...)
}

func (m *Model) insertRuneFromUserInput(values []rune) {
	head := m.inputText[:m.cursorPosition.X-1]
	tail := m.inputText[m.cursorPosition.X-1:]

	for _, value := range values {
		head = append(head, value)
	}

	m.cursorPosition.X = m.cursorPosition.X + len(values)
	m.inputText = append(head, tail...)
}

func (m *Model) SetWidth(width int) {
	m.width = width
}

func (m *Model) Reset() {
	m.cursorPosition.X = 1
	m.inputText = make([]rune, 0)
}

func (m *Model) Focus() {
	// m.cursor.Focus()
	m.isFocused = true
}

func (m *Model) Blur() {
	// m.cursor.Blur()
	m.isFocused = false
}

func TextInputBlink() tea.Msg {
	return Cursor.BlinkCursor()
}
