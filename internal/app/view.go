package app

import (
	Shared "godiff/pkg/common"
	"strings"
)

func (m Model) View() string {
	doc := strings.Builder{}

	switch m.screen {
	case Shared.LandingScreen:
		doc.WriteString(m.landingPage.View())
	case Shared.NewProjectWizardScreen:
		doc.WriteString(m.newProjectWizard.View())
	}

	return doc.String()
}
