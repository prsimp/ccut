package dashboard

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/prs/ccut/internal/data"
	"github.com/prs/ccut/internal/ui/components"
	"github.com/prs/ccut/internal/ui/format"
	"github.com/prs/ccut/internal/ui/styles"
)

type Model struct {
	Stats      *data.StatsCache
	Creds      *data.Credentials
	DailyCosts map[string]float64
	Width      int
}

func New(stats *data.StatsCache, creds *data.Credentials, dailyCosts map[string]float64) Model {
	return Model{Stats: stats, Creds: creds, DailyCosts: dailyCosts}
}

func (m Model) View() string {
	if m.Stats == nil {
		return styles.Subtle.Render("  Loading stats...")
	}

	var sections []string
	sections = append(sections, m.renderBadge())
	sections = append(sections, m.renderToday())
	sections = append(sections, m.renderLifetime())
	sections = append(sections, m.renderModelTable())
	sections = append(sections, m.renderDailyChart())
	sections = append(sections, m.renderHeatmap())
	return strings.Join(sections, "\n\n")
}

func (m Model) renderBadge() string {
	if m.Creds == nil {
		return ""
	}
	sub := strings.ToUpper(m.Creds.ClaudeAiOauth.SubscriptionType)
	tier := m.Creds.ClaudeAiOauth.RateLimitTier
	badge := styles.BadgeMax.Render(sub)
	if sub == "PRO" {
		badge = styles.BadgePro.Render(sub)
	}
	return "  " + badge + " " + styles.Subtle.Render(tier)
}

func (m Model) renderToday() string {
	title := styles.Bold.Render("  Today")
	if len(m.Stats.DailyActivity) == 0 {
		return title + "\n" + styles.Subtle.Render("  no data")
	}

	today := m.Stats.DailyActivity[len(m.Stats.DailyActivity)-1]

	totalTokens := 0
	if len(m.Stats.DailyModelTokens) > 0 {
		last := m.Stats.DailyModelTokens[len(m.Stats.DailyModelTokens)-1]
		for _, v := range last.TokensByModel {
			totalTokens += v
		}
	}

	boxes := []string{
		m.statBox("Messages", format.Number(today.MessageCount)),
		m.statBox("Sessions", format.Number(today.SessionCount)),
		m.statBox("Tool Calls", format.Number(today.ToolCallCount)),
		m.statBox("Tokens", format.Tokens(totalTokens)),
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, boxes...)
	return title + "\n" + row
}

func (m Model) statBox(label, value string) string {
	content := styles.StatValue.Render(value) + "\n" + styles.StatLabel.Render(label)
	return styles.Box.Width(16).Render(content)
}

func (m Model) renderLifetime() string {
	title := styles.Bold.Render("  Lifetime")
	s := m.Stats
	dur := format.Duration(s.LongestSession.Duration)

	totalCost := data.TotalCost(s.ModelUsage)
	lines := []string{
		fmt.Sprintf("  Sessions: %s    Messages: %s    Est. Cost: %s",
			styles.StatValue.Render(format.Number(s.TotalSessions)),
			styles.StatValue.Render(format.Number(s.TotalMessages)),
			styles.Special.Render(format.Cost(totalCost))),
		fmt.Sprintf("  Since: %s    Longest: %s (%d msgs)",
			styles.Special.Render(format.DateShort(s.FirstSessionDate)),
			styles.Special.Render(dur),
			s.LongestSession.MessageCount),
	}
	return title + "\n" + strings.Join(lines, "\n")
}

func (m Model) renderModelTable() string {
	title := styles.Bold.Render("  Model Usage")

	type modelRow struct {
		name  string
		usage data.ModelUsage
	}
	var rows []modelRow
	for name, usage := range m.Stats.ModelUsage {
		rows = append(rows, modelRow{name, usage})
	}
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].usage.OutputTokens > rows[j].usage.OutputTokens
	})

	hdr := styles.TableHeader.Render
	header := "  " +
		hdr(pad("Model", 16)) +
		hdr(padL("Input", 10)) +
		hdr(padL("Output", 10)) +
		hdr(padL("Cache Read", 12)) +
		hdr(padL("Cache Write", 12)) +
		hdr(padL("Cost", 10))

	var lines []string
	lines = append(lines, header)
	var totalCost float64
	for _, r := range rows {
		cost := data.CalculateCost(r.name, r.usage)
		totalCost += cost
		line := "  " +
			styles.Highlight.Render(pad(format.ShortModel(r.name), 16)) +
			padL(format.Tokens(r.usage.InputTokens), 10) +
			padL(format.Tokens(r.usage.OutputTokens), 10) +
			padL(format.Tokens(r.usage.CacheReadInputTokens), 12) +
			padL(format.Tokens(r.usage.CacheCreationInputTokens), 12) +
			styles.Special.Render(padL(format.Cost(cost), 10))
		lines = append(lines, line)
	}

	totalLine := "  " +
		pad("", 16) +
		padL("", 10) +
		padL("", 10) +
		padL("", 12) +
		styles.Bold.Render(padL("Total", 12)) +
		styles.Special.Render(padL(format.Cost(totalCost), 10))
	lines = append(lines, totalLine)

	return title + "\n" + strings.Join(lines, "\n")
}

func (m Model) renderDailyChart() string {
	title := styles.Bold.Render("  Daily Activity (last 14 days)")
	activities := m.Stats.DailyActivity
	if len(activities) > 14 {
		activities = activities[len(activities)-14:]
	}

	var items []components.SparkItem
	for _, a := range activities {
		suffix := ""
		if cost, ok := m.DailyCosts[a.Date]; ok && cost > 0 {
			suffix = format.Cost(cost)
		}
		items = append(items, components.SparkItem{
			Label:  format.DateShort(a.Date),
			Value:  a.MessageCount,
			Suffix: suffix,
		})
	}

	w := m.Width
	if w < 40 {
		w = 80
	}
	return title + "\n" + components.Sparkline(items, w-4)
}

func (m Model) renderHeatmap() string {
	title := styles.Bold.Render("  Hourly Activity (sessions started)")
	return title + "\n" + components.Heatmap(m.Stats.HourCounts)
}

func pad(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}

func padL(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return strings.Repeat(" ", width-len(s)) + s
}
