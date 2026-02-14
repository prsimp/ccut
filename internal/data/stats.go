package data

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

func ClaudeDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".claude")
}

func StatsPath() string {
	return filepath.Join(ClaudeDir(), "stats-cache.json")
}

func LoadStats() (*StatsCache, error) {
	data, err := os.ReadFile(StatsPath())
	if err != nil {
		return nil, err
	}
	var stats StatsCache
	if err := json.Unmarshal(data, &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}

// StatsMtime returns the modification time of stats-cache.json.
func StatsMtime() time.Time {
	info, err := os.Stat(StatsPath())
	if err != nil {
		return time.Time{}
	}
	return info.ModTime()
}
