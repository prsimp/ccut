package live

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/prs/ccut/internal/data"
	"github.com/prs/ccut/internal/ui/format"
	"github.com/prs/ccut/internal/ui/keys"
	"github.com/prs/ccut/internal/ui/styles"
)

type TickMsg time.Time

type Model struct {
	stats      *data.StatsCache
	prevStats  *data.StatsCache
	spinner    spinner.Model
	lastUpdate time.Time
	width      int
}

func New(stats *data.StatsCache) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styles.Special
	return Model{
		stats:      stats,
		spinner:    s,
		lastUpdate: data.StatsMtime(),
	}
}

func (m *Model) SetSize(w, _ int) {
	m.width = w
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, tickCmd())
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case data.StatsUpdatedMsg:
		m.prevStats = m.stats
		m.stats = msg.Stats
		m.lastUpdate = data.StatsMtime()
		return m, data.WatchStats()
	case TickMsg:
		return m, tickCmd()
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case tea.KeyMsg:
		if key.Matches(msg, keys.Map.Refresh) {
			stats, err := data.LoadStats()
			if err == nil {
				m.prevStats = m.stats
				m.stats = stats
				m.lastUpdate = data.StatsMtime()
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	if m.stats == nil {
		return styles.Subtle.Render("  Waiting for data...")
	}

	var lines []string

	elapsed := time.Since(m.lastUpdate).Truncate(time.Second)
	lines = append(lines, fmt.Sprintf("  %s Live Monitor   %s  Last updated: %s",
		m.spinner.View(),
		styles.Subtle.Render("watching stats-cache.json"),
		styles.Special.Render(elapsed.String()+" ago")))
	lines = append(lines, "")

	if len(m.stats.DailyActivity) > 0 {
		today := m.stats.DailyActivity[len(m.stats.DailyActivity)-1]

		var prevDay *data.DailyActivity
		if m.prevStats != nil && len(m.prevStats.DailyActivity) > 0 {
			p := m.prevStats.DailyActivity[len(m.prevStats.DailyActivity)-1]
			prevDay = &p
		}

		lines = append(lines, styles.Bold.Render("  Today"))
		lines = append(lines, m.statLine("Messages", today.MessageCount, prevDay, func(d *data.DailyActivity) int { return d.MessageCount }))
		lines = append(lines, m.statLine("Sessions", today.SessionCount, prevDay, func(d *data.DailyActivity) int { return d.SessionCount }))
		lines = append(lines, m.statLine("Tool Calls", today.ToolCallCount, prevDay, func(d *data.DailyActivity) int { return d.ToolCallCount }))

		totalTokens := 0
		if len(m.stats.DailyModelTokens) > 0 {
			last := m.stats.DailyModelTokens[len(m.stats.DailyModelTokens)-1]
			for _, v := range last.TokensByModel {
				totalTokens += v
			}
		}
		prevTokens := 0
		if m.prevStats != nil && len(m.prevStats.DailyModelTokens) > 0 {
			last := m.prevStats.DailyModelTokens[len(m.prevStats.DailyModelTokens)-1]
			for _, v := range last.TokensByModel {
				prevTokens += v
			}
		}
		delta := ""
		if m.prevStats != nil {
			d := totalTokens - prevTokens
			if d > 0 {
				delta = styles.DeltaUp.Render(fmt.Sprintf(" (+%s)", format.Tokens(d)))
			}
		}
		lines = append(lines, fmt.Sprintf("    %-14s %s%s", "Tokens", styles.StatValue.Render(format.Tokens(totalTokens)), delta))
	}

	lines = append(lines, "")
	lines = append(lines, styles.Bold.Render("  Totals"))
	lines = append(lines, fmt.Sprintf("    Sessions: %s    Messages: %s",
		styles.StatValue.Render(format.Number(m.stats.TotalSessions)),
		styles.StatValue.Render(format.Number(m.stats.TotalMessages))))

	lines = append(lines, "")
	lines = append(lines, styles.Subtle.Render("  r → force refresh"))

	return strings.Join(lines, "\n")
}

func (m Model) statLine(label string, val int, prev *data.DailyActivity, getter func(*data.DailyActivity) int) string {
	delta := ""
	if prev != nil {
		d := val - getter(prev)
		if d > 0 {
			delta = styles.DeltaUp.Render(fmt.Sprintf(" (+%s)", format.Number(d)))
		} else if d < 0 {
			delta = styles.DeltaDown.Render(fmt.Sprintf(" (%s)", format.Number(d)))
		}
	}
	return fmt.Sprintf("    %-14s %s%s", label, styles.StatValue.Render(format.Number(val)), delta)
}
