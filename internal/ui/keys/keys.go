package keys

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Tab      key.Binding
	ShiftTab key.Binding
	Tab1     key.Binding
	Tab2     key.Binding
	Tab3     key.Binding
	Up       key.Binding
	Down     key.Binding
	Enter    key.Binding
	Back     key.Binding
	Refresh  key.Binding
	Filter   key.Binding
	Group    key.Binding
	Quit     key.Binding
}

var Map = KeyMap{
	Tab:      key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next tab")),
	ShiftTab: key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "prev tab")),
	Tab1:     key.NewBinding(key.WithKeys("1"), key.WithHelp("1", "dashboard")),
	Tab2:     key.NewBinding(key.WithKeys("2"), key.WithHelp("2", "sessions")),
	Tab3:     key.NewBinding(key.WithKeys("3"), key.WithHelp("3", "live")),
	Up:       key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("k", "up")),
	Down:     key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("j", "down")),
	Enter:    key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
	Back:     key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
	Refresh:  key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "refresh")),
	Filter:   key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "filter")),
	Group:    key.NewBinding(key.WithKeys("g"), key.WithHelp("g", "group")),
	Quit:     key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
}
