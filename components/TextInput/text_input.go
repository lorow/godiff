package TextInput

import (
	"strings"

	"godiff/components/Cursor"
	"godiff/components/FocusChain"

	tea "github.com/charmbracelet/bubbletea"
)

type CursorPosition struct {
	X int
}

type Model struct {
	prompt         string
	inputText      []rune
	sanitizer      *Sanitizer
	cursor         Cursor.Model
	cursorPosition CursorPosition
	defaultStyle   Styles
	focusedStyle   Styles
	onSubmit       func(string) tea.Cmd
	width          int
	isFocused      bool
}

func New(opts ...func(*Model)) *Model {
	model := &Model{
		prompt:         ">",
		cursor:         Cursor.New(),
		sanitizer:      NewSanitizer(),
		inputText:      make([]rune, 0),
		cursorPosition: CursorPosition{X: 1},
		defaultStyle:   DefaultStyles(),
		focusedStyle:   focusedStyles(),
		width:          0,
		isFocused:      false,
	}

	for _, opt := range opts {
		opt(model)
	}

	return model
}

func (m *Model) Init() tea.Cmd {
	return Cursor.Blink
}

func WithPrompt(prompt string) func(model *Model) {
	return func(model *Model) {
		model.prompt = prompt
	}
}

func WIthRegularStyles(styles Styles) func(model *Model) {
	return func(model *Model) {
		model.defaultStyle = styles
	}
}

func WithFocusedStyles(styles Styles) func(model *Model) {
	return func(model *Model) {
		model.focusedStyle = styles
	}
}

func WIthOnSubmit(onSubmit func(string) tea.Cmd) func(model *Model) {
	return func(model *Model) {
		model.onSubmit = onSubmit
	}
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	// todo add keymaps
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyDown:
			return FocusChain.SwitchFocusCmd(FocusChain.FocusDown)
		case tea.KeyEnter:
			if m.onSubmit != nil {
				return m.onSubmit(string(m.inputText))
			}
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

	return tea.Batch(cmds...)
}

func (m Model) View() string {
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

	style := m.getStyle()

	return style.Container.Width(m.width).Render(v.String())
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

	sanitized_runes := m.sanitizer.Sanitize(values)
	for _, value := range sanitized_runes {
		head = append(head, value)
	}

	m.cursorPosition.X = m.cursorPosition.X + len(sanitized_runes)
	m.inputText = append(head, tail...)
}

func (m *Model) SetWidth(width int) {
	m.width = width
}

func (m *Model) Reset() {
	m.cursorPosition.X = 1
	m.inputText = make([]rune, 0)
}

func (m *Model) Focus() tea.Cmd {
	cmd := m.cursor.Focus()
	m.isFocused = true
	return cmd
}

func (m *Model) Blur() tea.Cmd {
	m.cursor.Blur()
	m.isFocused = false
	return nil
}

func (m Model) getStyle() Styles {
	if m.isFocused {
		return m.focusedStyle
	}
	return m.defaultStyle
}
