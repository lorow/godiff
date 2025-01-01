package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type LandingPageModel struct {
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
	return m, nil
}

func (m LandingPageModel) View() string {
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

	return s
}
