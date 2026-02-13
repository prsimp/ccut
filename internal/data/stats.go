package data

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func ClaudeDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".claude")
}

func LoadStats() (*StatsCache, error) {
	path := filepath.Join(ClaudeDir(), "stats-cache.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var stats StatsCache
	if err := json.Unmarshal(data, &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}
