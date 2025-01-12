package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type LandingPageModel struct {
	width  int
	height int
	// replace this with enum
	hasLoadedProjects bool
	loadingProjects   bool
	projects          []Project
	cursor            int
	// generally, this will be removed
	// instead, we will be sending a RouteToEditor
	// command with the selected project id
	// which we will then load up and display
	selected map[int]struct{}
}

type SelectedProject int

func NewLandingPage() LandingPageModel {
	return LandingPageModel{
		selected: make(map[int]struct{}),
	}
}

func (m LandingPageModel) Init() tea.Cmd {
	return nil
}

func (m LandingPageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// we need to first handle the incoming messages
	// only then we can send our own IO commands

	switch msg := msg.(type) {
	case SetNewSizeMsg:
		m.width = msg.width
		m.height = msg.height
		return m, nil
	case LoadedProjectsMsg:
		m.hasLoadedProjects = true
		m.loadingProjects = false
		m.projects = msg.projects
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil
		case "down":
			if m.cursor < len(m.projects)-1 {
				m.cursor++
			}
			return m, nil
		case "enter":
			//return m, RouteTo("editor", SelectedProject(m.projects[m.cursor].id))

			if _, ok := m.selected[m.cursor]; ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
			return m, nil
		}
	}

	if !m.hasLoadedProjects && !m.loadingProjects {
		m.loadingProjects = true
		return m, LoadProjectsCmd(10, 0)
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

	for i, choice := range m.projects {
		cursor := " "

		if i == m.cursor {
			cursor = ">"
		}

		checked := " "

		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice.name)
	}

	doc.WriteString(s + "\n\n")
	doc.WriteString(fmt.Sprintf("height: %d \n", m.height))
	return windowContainer.Render(doc.String())
}

type LoadedProjectsMsg struct {
	projects []Project
}

func LoadProjectsCmd(limit, offset int) tea.Cmd {
	return func() tea.Msg {
		projects := GetProjects(limit, offset)
		return LoadedProjectsMsg{projects}
	}
}
