package ui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/prs/ccut/internal/data"
	"github.com/prs/ccut/internal/ui/components"
	"github.com/prs/ccut/internal/ui/dashboard"
	"github.com/prs/ccut/internal/ui/keys"
	"github.com/prs/ccut/internal/ui/live"
	"github.com/prs/ccut/internal/ui/sessions"
)

type RootModel struct {
	activeTab int
	dashboard dashboard.Model
	sessions  sessions.Model
	live      live.Model
	width     int
	height    int
	ready     bool
}

func NewRoot() RootModel {
	stats, _ := data.LoadStats()
	creds, _ := data.LoadCredentials()

	return RootModel{
		dashboard: dashboard.New(stats, creds),
		sessions:  sessions.New(),
		live:      live.New(stats),
	}
}

func (m RootModel) Init() tea.Cmd {
	return tea.Batch(
		sessions.LoadSessions(),
		data.WatchStats(),
		m.live.Init(),
	)
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		m.dashboard.Width = msg.Width
		m.sessions.SetSize(msg.Width, msg.Height-3)
		m.live.SetSize(msg.Width, msg.Height-3)
		return m, nil

	case tea.KeyMsg:
		if key.Matches(msg, keys.Map.Quit) {
			return m, tea.Quit
		}
		if key.Matches(msg, keys.Map.Tab) {
			m.activeTab = (m.activeTab + 1) % 3
			return m, nil
		}
		if key.Matches(msg, keys.Map.ShiftTab) {
			m.activeTab = (m.activeTab + 2) % 3
			return m, nil
		}
		if key.Matches(msg, keys.Map.Tab1) {
			m.activeTab = 0
			return m, nil
		}
		if key.Matches(msg, keys.Map.Tab2) {
			m.activeTab = 1
			return m, nil
		}
		if key.Matches(msg, keys.Map.Tab3) {
			m.activeTab = 2
			return m, nil
		}

	case data.StatsUpdatedMsg:
		m.dashboard.Stats = msg.Stats
		var cmd tea.Cmd
		m.live, cmd = m.live.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)

	case sessions.SessionsLoadedMsg:
		var cmd tea.Cmd
		m.sessions, cmd = m.sessions.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}

	switch m.activeTab {
	case 1:
		var cmd tea.Cmd
		m.sessions, cmd = m.sessions.Update(msg)
		cmds = append(cmds, cmd)
	case 2:
		var cmd tea.Cmd
		m.live, cmd = m.live.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m RootModel) View() string {
	if !m.ready {
		return "  Loading..."
	}

	header := components.Header(m.activeTab, m.width)

	var body string
	switch m.activeTab {
	case 0:
		body = m.dashboard.View()
	case 1:
		body = m.sessions.View()
	case 2:
		body = m.live.View()
	}

	content := lipgloss.JoinVertical(lipgloss.Left, header, "", body)
	return lipgloss.Place(m.width, m.height, lipgloss.Left, lipgloss.Top, content)
}
