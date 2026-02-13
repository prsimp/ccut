package data

import "testing"

func TestIsStatsFile(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"exact match", "stats-cache.json", true},
		{"with directory prefix", "/home/user/.claude/stats-cache.json", true},
		{"with relative prefix", "some/dir/stats-cache.json", true},
		{"wrong name", "other-file.json", false},
		{"empty string", "", false},
		{"partial match", "cache.json", false},
		{"prefix only", "stats-cache", false},
		{"too short", "stats-cache.jso", false},
		{"substring but not suffix", "stats-cache.json.bak", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isStatsFile(tt.input)
			if got != tt.want {
				t.Errorf("isStatsFile(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
