package app

import (
	LandingPage "godiff/internal/app/screens/landing_page"
	NewProjectWizard "godiff/internal/app/screens/new_project_wizard"
	Shared "godiff/pkg/common"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	size             tea.WindowSizeMsg
	screen           Shared.Screen
	landingPage      LandingPage.Model
	newProjectWizard NewProjectWizard.Model
}

func New() Model {
	m := Model{
		screen:           Shared.LandingScreen,
		landingPage:      LandingPage.New(),
		newProjectWizard: NewProjectWizard.New(),
	}
	return m
}
