package newProjectWizard

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	windowContainer := lipgloss.NewStyle().Width(m.width).Height(m.height)
	doc := strings.Builder{}

	title := m.titlePanel.View()
	middleSpacer := lipgloss.NewStyle().Width(m.width - lipgloss.Width(title)).Render("")
	renderedTitle := lipgloss.NewStyle().PaddingTop(1).PaddingBottom(1).Render(lipgloss.JoinHorizontal(lipgloss.Top, title, middleSpacer))

	projectNameView := m.projectName.View()
	projectDescriptionView := m.projectDescription.View()
	shortcutsView := m.shortcutsPanel.View()

	spaceFiller := lipgloss.NewStyle().Width(m.width).Height(m.height - lipgloss.Height(renderedTitle) - lipgloss.Height(projectNameView) - lipgloss.Height(projectDescriptionView) - lipgloss.Height(shortcutsView) - 1).Render("")

	doc.WriteString(
		lipgloss.JoinVertical(
			lipgloss.Top,
			renderedTitle,
			projectNameView,
			projectDescriptionView,
			spaceFiller,
			shortcutsView,
		),
	)
	return windowContainer.Render(doc.String())
}
