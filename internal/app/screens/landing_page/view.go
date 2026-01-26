package landingPage

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	windowContainer := lipgloss.NewStyle().Width(m.width).Height(m.height)
	doc := strings.Builder{}

	title := m.titlePanel.View()
	quitText := lipgloss.NewStyle().PaddingRight(2).Render("Press Q to quit")
	middleSpacer := lipgloss.NewStyle().Width(m.width - lipgloss.Width(title) - lipgloss.Width(quitText)).Render("")
	renderedTitle := lipgloss.NewStyle().PaddingTop(1).PaddingBottom(1).Render(lipgloss.JoinHorizontal(lipgloss.Top, title, middleSpacer, quitText))

	doc.WriteString(
		lipgloss.JoinVertical(
			lipgloss.Top,
			renderedTitle,
			m.searchInput.View(),
			m.itemList.View(),
			m.shortcutsPanel.View(),
		),
	)

	return windowContainer.Render(doc.String())
}

func (p RenderableProject) Title() string {
	return fmt.Sprintf("%d - %s", p.Project.ID, p.Project.Name)
}

func (p RenderableProject) Description() string {
	return "Dummy description"
}
