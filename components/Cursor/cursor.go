package Cursor

import (
	"context"
	"errors"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var lastId = 0

type CursorMode int

type InitCursorBlinkMsg struct{}

type CursorBlinkMsg struct {
	id int
}

type blinkCanceled struct{}

type cursorBlinkCtx struct {
	ctx    context.Context
	cancel context.CancelFunc
}

const (
	CursorBlink CursorMode = iota
	CursorStatic
	CursorHide
)

type Model struct {
	cursorStyle lipgloss.Style
	textStyle   lipgloss.Style
	mode        CursorMode
	blinkSpeed  time.Duration
	char        string
	id          int
	focus       bool
	blink       bool
	blinkCtx    *cursorBlinkCtx
}

func NewCursorModel() Model {
	cursorModel := Model{
		id:         lastId,
		mode:       CursorBlink,
		blinkSpeed: time.Millisecond * 500,
		blink:      true,
		blinkCtx: &cursorBlinkCtx{
			ctx: context.Background(),
		},
	}

	lastId++

	return cursorModel
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case InitCursorBlinkMsg:
		if m.mode != CursorBlink || !m.focus {
			return m, nil
		}

		cmd := m.BlinkCmd()
		return m, cmd

	case CursorBlinkMsg:
		// we only want to blink for the current cursor
		if msg.id != m.id {
			return m, nil
		}

		if m.mode != CursorBlink || !m.focus {
			return m, nil
		}

		if m.mode == CursorBlink {
			m.blink = !m.blink
			return m, m.BlinkCmd()
		}
		return m, nil

	case blinkCanceled:
		return m, nil
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
		return BlinkCursor
	}

	return nil
}

func (m Model) Mode() CursorMode {
	return m.mode
}

func (m *Model) TurnBlinkOff() {
	m.blink = false
}

func (m *Model) SetChar(char string) {
	m.char = char
}

func (m *Model) Focus() {
	m.focus = true
}

func (m *Model) Blur() {
	m.focus = false
}

func BlinkCursor() tea.Msg {
	return InitCursorBlinkMsg{}
}

// todo refactor this to use tea commands instead of channels
func (m *Model) BlinkCmd() tea.Cmd {
	if m.mode != CursorBlink {
		return nil
	}

	if m.blinkCtx != nil && m.blinkCtx.cancel != nil {
		m.blinkCtx.cancel()
	}

	ctx, cancel := context.WithTimeout(m.blinkCtx.ctx, m.blinkSpeed)
	m.blinkCtx.cancel = cancel

	return func() tea.Msg {
		defer cancel()

		<-ctx.Done()
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return CursorBlinkMsg{
				id: m.id,
			}
		}
		// since we're still waiting
		// for the next blink step to happen
		// we can send a noop
		return blinkCanceled{}
	}
}
