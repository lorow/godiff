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
	selected          int
}

type SelectedProject int

func NewLandingPage() LandingPageModel {
	return LandingPageModel{
		selected: -1,
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
			m.selected = m.cursor
			return m, RouteTo("editor", SelectedProject(m.projects[m.cursor].id))
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

		if i == m.selected {
			checked = "█"
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
