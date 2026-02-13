package data

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fsnotify/fsnotify"
)

type StatsUpdatedMsg struct {
	Stats *StatsCache
}

// WatchStats returns a tea.Cmd that watches stats-cache.json for changes.
func WatchStats() tea.Cmd {
	return func() tea.Msg {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			return nil
		}
		defer watcher.Close()

		dir := ClaudeDir()
		if err := watcher.Add(dir); err != nil {
			return nil
		}

		debounce := time.NewTimer(time.Hour)
		debounce.Stop()

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return nil
				}
				// Check Write, Create, and Rename — atomic writes (temp+rename)
				// produce Create/Rename events, not Write.
				if isStatsFile(event.Name) &&
					event.Has(fsnotify.Write|fsnotify.Create|fsnotify.Rename) {
					debounce.Reset(500 * time.Millisecond)
				}
			case <-debounce.C:
				stats, err := LoadStats()
				if err == nil {
					return StatsUpdatedMsg{Stats: stats}
				}
			case _, ok := <-watcher.Errors:
				if !ok {
					return nil
				}
			}
		}
	}
}

func isStatsFile(name string) bool {
	return len(name) >= 16 && name[len(name)-16:] == "stats-cache.json"
}
