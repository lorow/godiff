package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type LandingPageModel struct {
	width    int
	height   int
	choices  []string
	cursor   int
	selected map[int]struct{}
}

func NewLandingPage() LandingPageModel {
	return LandingPageModel{
		choices:  []string{"Buy carrots", "Buy celery", "Buy kohlrabi"},
		selected: make(map[int]struct{}),
	}
}

func (m LandingPageModel) Init() tea.Cmd {
	return nil
}

func (m LandingPageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case SetNewSizeMsg:
		m.width = msg.width
		m.height = msg.height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil
		case "down":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
			return m, nil
		case "enter":
			if _, ok := m.selected[m.cursor]; ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
			return m, nil
		}
	}
	return m, nil
}

func (m LandingPageModel) View() string {
	windowContainer := lipgloss.NewStyle().Width(m.width).Height(m.height)
	doc := strings.Builder{}

	desc := lipgloss.JoinVertical(lipgloss.Left,
		descStyle.Render("Style Definitions for Nice Terminal Layouts"),
		infoStyle.Render("From Charm"+divider+url("https://github.com/charmbracelet/lipgloss")),
	)

	row := lipgloss.JoinHorizontal(lipgloss.Top, "", desc)
	doc.WriteString(row + "\n\n")

	s := ""

	for i, choice := range m.choices {
		cursor := " "

		if i == m.cursor {
			cursor = ">"
		}

		checked := " "

		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	doc.WriteString(s + "\n\n")
	doc.WriteString(fmt.Sprintf("height: %d \n", m.height))
	return windowContainer.Render(doc.String())
}
