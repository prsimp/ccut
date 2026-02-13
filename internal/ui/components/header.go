package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/prs/ccut/internal/ui/styles"
)

var tabs = []string{"Dashboard", "Sessions", "Live"}

func Header(activeTab int, width int) string {
	title := styles.Title.Render("⚡ ccut")

	var renderedTabs []string
	for i, t := range tabs {
		if i == activeTab {
			renderedTabs = append(renderedTabs, styles.ActiveTab.Render(t))
		} else {
			renderedTabs = append(renderedTabs, styles.InactiveTab.Render(t))
		}
	}
	tabBar := lipgloss.JoinHorizontal(lipgloss.Bottom, renderedTabs...)

	row := lipgloss.JoinHorizontal(lipgloss.Center, title, "  ", tabBar)
	line := styles.Subtle.Render(strings.Repeat("─", max(0, width-lipgloss.Width(row))))
	return lipgloss.JoinHorizontal(lipgloss.Bottom, row, line)
}
