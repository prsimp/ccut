package styles

import "github.com/charmbracelet/lipgloss"

var (
	Subtle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	Highlight = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	Special   = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	Warn      = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))

	Bold     = lipgloss.NewStyle().Bold(true)
	Title    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("99"))
	BadgeMax = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("0")).Background(lipgloss.Color("212")).Padding(0, 1)
	BadgePro = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("0")).Background(lipgloss.Color("86")).Padding(0, 1)

	Box = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(0, 1)

	ActiveTab = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("212")).
			Underline(true).
			Padding(0, 2)

	InactiveTab = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Padding(0, 2)

	StatValue = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86"))
	StatLabel = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

	TableHeader = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("99"))
	TableCell   = lipgloss.NewStyle()

	BarFull  = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	BarEmpty = lipgloss.NewStyle().Foreground(lipgloss.Color("236"))

	HeatLow    = lipgloss.NewStyle().Foreground(lipgloss.Color("236"))
	HeatMedLow = lipgloss.NewStyle().Foreground(lipgloss.Color("60"))
	HeatMed    = lipgloss.NewStyle().Foreground(lipgloss.Color("99"))
	HeatHigh   = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))

	DeltaUp   = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	DeltaDown = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
)
