package Cursor

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var lastId = 0

type CursorBlinkMsg int

type CursorMode int

const (
	CursorBlink CursorMode = iota
	CursorStatic
	CursorHide
)

type initialBlinkMsg struct{}

type BlinkMsg struct {
	id     int
	tag_id int
}

type blinkCanceled struct{}

type blinkCtx struct {
	ctx    context.Context
	cancel context.CancelFunc
}

type Model struct {
	cursorStyle lipgloss.Style
	textStyle   lipgloss.Style
	mode        CursorMode
	blinkSpeed  time.Duration
	char        string
	id          int
	tag_id      int
	focus       bool
	blink       bool
	blinkCtx    *blinkCtx
}

func New() Model {
	cursorModel := Model{
		blinkSpeed: time.Millisecond * 500,
		mode:       CursorBlink,
		blink:      true,
		id:         lastId,
		blinkCtx: &blinkCtx{
			ctx: context.Background(),
		},
	}

	lastId++

	return cursorModel
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {

	switch msg := msg.(type) {
	case initialBlinkMsg:
		if m.mode != CursorBlink || !m.focus {
			return m, nil
		}

		return m, m.BlinkCursorCmd()

	case BlinkMsg:
		if m.mode != CursorBlink || !m.focus {
			return m, nil
		}

		if msg.id != m.id || msg.tag_id != m.tag_id {
			return m, nil
		}

		var cmd tea.Cmd
		if m.mode == CursorBlink {
			m.blink = !m.blink
			cmd = m.BlinkCursorCmd()
		}

		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	if m.blink {
		return m.textStyle.Inline(true).Render(m.char)
	}
	return m.cursorStyle.Inline(true).Reverse(true).Render(m.char)
}

func (m *Model) SetMode(mode CursorMode) tea.Cmd {
	if mode < CursorBlink || mode > CursorHide {
		return nil
	}

	m.mode = mode
	m.blink = m.mode == CursorHide || !m.focus

	if mode == CursorBlink {
		return m.BlinkCursorCmd()
	}

	return nil
}

func (m Model) Mode() CursorMode {
	return m.mode
}

func (m *Model) SetChar(char string) {
	m.char = char
}

func (m *Model) Focus() tea.Cmd {
	m.focus = true
	m.blink = m.mode != CursorHide

	if m.mode == CursorBlink {
		return m.BlinkCursorCmd()
	}

	return nil
}

func (m *Model) Blur() {
	m.focus = false
	m.blink = false
}

func (m *Model) IsFocused() bool {
	return m.focus
}

func (m *Model) BlinkCursorCmd() tea.Cmd {
	if m.mode != CursorBlink {
		return nil
	}

	if m.blinkCtx != nil && m.blinkCtx.cancel != nil {
		m.blinkCtx.cancel()
	}

	ctx, cancel := context.WithTimeout(m.blinkCtx.ctx, m.blinkSpeed)
	m.blinkCtx.cancel = cancel

	m.tag_id++
	blinkMsg := BlinkMsg{id: m.id, tag_id: m.tag_id}

	return func() tea.Msg {
		defer cancel()
		<-ctx.Done()
		if ctx.Err() == context.DeadlineExceeded {
			return blinkMsg
		}

		return blinkCanceled{}
	}
}

// Initialize blinking, this is accepted by any and all cursors, regardless of their focus.
// if the current view is no receiving input currently, you may need to refire the command
func Blink() tea.Msg {
	return initialBlinkMsg{}
}
