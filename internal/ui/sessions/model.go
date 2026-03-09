package sessions

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/prs/ccut/internal/data"
	"github.com/prs/ccut/internal/ui/format"
	"github.com/prs/ccut/internal/ui/keys"
	"github.com/prs/ccut/internal/ui/styles"
)

type SessionsLoadedMsg struct {
	Sessions []data.SessionInfo
}

type Model struct {
	sessions []data.SessionInfo
	cursor   int
	detail   bool
	offset   int
	height   int
	width    int
	loading  bool
}

func New() Model {
	return Model{loading: true}
}

func LoadSessions() tea.Cmd {
	return func() tea.Msg {
		sessions, _ := data.DiscoverSessions(50)
		return SessionsLoadedMsg{Sessions: sessions}
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case SessionsLoadedMsg:
		m.sessions = msg.Sessions
		m.loading = false
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Map.Down):
			if !m.detail && m.cursor < len(m.sessions)-1 {
				m.cursor++
				m.ensureVisible()
			}
		case key.Matches(msg, keys.Map.Up):
			if !m.detail && m.cursor > 0 {
				m.cursor--
				m.ensureVisible()
			}
		case key.Matches(msg, keys.Map.Enter):
			if len(m.sessions) > 0 {
				m.detail = true
			}
		case key.Matches(msg, keys.Map.Back):
			m.detail = false
		}
	}
	return m, nil
}

func (m *Model) ensureVisible() {
	visible := m.visibleLines()
	if m.cursor < m.offset {
		m.offset = m.cursor
	}
	if m.cursor >= m.offset+visible {
		m.offset = m.cursor - visible + 1
	}
}

func (m Model) visibleLines() int {
	if m.height < 6 {
		return 20
	}
	return m.height - 4
}

func (m *Model) SetSize(w, h int) {
	m.width = w
	m.height = h
}

func (m Model) View() string {
	if m.loading {
		return styles.Subtle.Render("  Loading sessions...")
	}
	if len(m.sessions) == 0 {
		return styles.Subtle.Render("  No sessions found")
	}
	if m.detail {
		return m.renderDetail()
	}
	return m.renderList()
}

func (m Model) renderList() string {
	title := styles.Bold.Render("  Sessions")

	h := styles.TableHeader.Render
	header := "  " +
		h(pad("Project", 30)) +
		h(pad("Date", 14)) +
		h(padLeft("Msgs", 6)) + " " +
		h(pad("Model", 14)) +
		h(padLeft("Output", 8)) + " " +
		h(padLeft("Cost", 8))

	visible := m.visibleLines()
	end := m.offset + visible
	if end > len(m.sessions) {
		end = len(m.sessions)
	}

	var lines []string
	lines = append(lines, title, header)
	for i := m.offset; i < end; i++ {
		s := m.sessions[i]
		cursor := "  "
		nameStyle := styles.Highlight
		if i == m.cursor {
			cursor = styles.Special.Render("▸ ")
			nameStyle = styles.Special
		}

		project := truncate(s.ProjectName, 30)
		date := format.Date(s.LastTime)
		msgs := fmt.Sprintf("%d", s.MessageCount)
		model := format.ShortModel(s.Model)
		output := format.Tokens(s.OutputTokens)
		cost := format.Cost(data.CalculateCost(s.Model, s.TokenUsage))

		line := cursor +
			nameStyle.Render(pad(project, 30)) +
			styles.Subtle.Render(pad(date, 14)) +
			styles.StatValue.Render(padLeft(msgs, 6)) + " " +
			styles.Subtle.Render(pad(model, 14)) +
			styles.StatValue.Render(padLeft(output, 8)) + " " +
			styles.Special.Render(padLeft(cost, 8))

		lines = append(lines, line)
	}

	info := styles.Subtle.Render(fmt.Sprintf("  %d sessions | ↑↓ navigate | enter select", len(m.sessions)))
	lines = append(lines, "", info)
	return strings.Join(lines, "\n")
}

func pad(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}

func padLeft(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return strings.Repeat(" ", width-len(s)) + s
}

func (m Model) renderDetail() string {
	if m.cursor >= len(m.sessions) {
		return ""
	}
	s := m.sessions[m.cursor]

	var lines []string
	lines = append(lines, styles.Bold.Render("  Session Detail"))
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("  Project:    %s", styles.Highlight.Render(s.ProjectName)))
	lines = append(lines, fmt.Sprintf("  Session ID: %s", styles.Subtle.Render(s.SessionID)))
	lines = append(lines, fmt.Sprintf("  Branch:     %s", styles.Special.Render(s.GitBranch)))
	lines = append(lines, fmt.Sprintf("  Model:      %s", format.ShortModel(s.Model)))
	lines = append(lines, fmt.Sprintf("  Messages:   %s", styles.StatValue.Render(format.Number(s.MessageCount))))
	lines = append(lines, fmt.Sprintf("  Output:     %s tokens", styles.StatValue.Render(format.Tokens(s.OutputTokens))))
	lines = append(lines, fmt.Sprintf("  Est. Cost:  %s", styles.Special.Render(format.Cost(data.CalculateCost(s.Model, s.TokenUsage)))))
	lines = append(lines, fmt.Sprintf("  Started:    %s", format.Date(s.FirstTime)))
	lines = append(lines, fmt.Sprintf("  Last:       %s", format.Date(s.LastTime)))
	lines = append(lines, "")
	lines = append(lines, styles.Subtle.Render("  esc → back"))
	return strings.Join(lines, "\n")
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
