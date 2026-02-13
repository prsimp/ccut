package format

import (
	"fmt"
	"testing"
	"time"
)

func TestNumber(t *testing.T) {
	tests := []struct {
		input int
		want  string
	}{
		{0, "0"},
		{1, "1"},
		{12, "12"},
		{123, "123"},
		{1234, "1,234"},
		{12345, "12,345"},
		{123456, "123,456"},
		{1234567, "1,234,567"},
		{1000000000, "1,000,000,000"},
		{-1, "-1"},
		{-1234, "-1,234"},
		{-1234567, "-1,234,567"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d", tt.input), func(t *testing.T) {
			got := Number(tt.input)
			if got != tt.want {
				t.Errorf("Number(%d) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestTokens(t *testing.T) {
	tests := []struct {
		input int
		want  string
	}{
		{0, "0"},
		{999, "999"},
		{1000, "1.0K"},
		{1500, "1.5K"},
		{10000, "10.0K"},
		{999999, "1000.0K"},
		{1000000, "1.0M"},
		{2500000, "2.5M"},
		{1000000000, "1.0B"},
		{2500000000, "2.5B"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d", tt.input), func(t *testing.T) {
			got := Tokens(tt.input)
			if got != tt.want {
				t.Errorf("Tokens(%d) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestDuration(t *testing.T) {
	tests := []struct {
		ms   int64
		want string
	}{
		{0, "0m"},
		{30000, "0m"},        // 30 seconds rounds to 0m
		{60000, "1m"},        // 1 minute
		{90000, "1m"},        // 1.5 min truncates to 1m
		{3600000, "1h 0m"},   // 1 hour
		{5400000, "1h 30m"},  // 1.5 hours
		{7200000, "2h 0m"},   // 2 hours
		{7260000, "2h 1m"},   // 2 hours 1 min
		{86400000, "24h 0m"}, // 24 hours
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%dms", tt.ms), func(t *testing.T) {
			got := Duration(tt.ms)
			if got != tt.want {
				t.Errorf("Duration(%d) = %q, want %q", tt.ms, got, tt.want)
			}
		})
	}
}

func TestDate(t *testing.T) {
	// Date converts to local time, so we construct expected values using Local timezone.
	tests := []struct {
		input string
		want  string
	}{
		// RFC3339
		{
			"2024-06-15T14:30:00Z",
			time.Date(2024, 6, 15, 14, 30, 0, 0, time.UTC).Local().Format("Jan 02 15:04"),
		},
		// Plain date
		{
			"2024-01-02",
			time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC).Local().Format("Jan 02 15:04"),
		},
		// Invalid returns as-is
		{"not-a-date", "not-a-date"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := Date(tt.input)
			if got != tt.want {
				t.Errorf("Date(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestDateShort(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"2024-06-15", "Jun 15"},
		{"2024-01-02", "Jan 02"},
		{"2024-12-25T10:00:00Z", "Dec 25"},
		{"garbage", "garbage"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := DateShort(tt.input)
			if got != tt.want {
				t.Errorf("DateShort(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestShortModel(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"claude-3-5-sonnet-20241022", "3-5-sonnet"},
		{"claude-sonnet-4-20250514", "sonnet-4"},
		{"claude-opus-4-20250514", "opus-4"},
		{"claude-3-5-haiku-20241022", "3-5-haiku"},
		// No date suffix
		{"claude-sonnet-4", "sonnet-4"},
		// Not a claude model
		{"gpt-4o-20240101", "gpt-4o"},
		// Short string (no stripping)
		{"short", "short"},
		// Exactly 9 chars with no date
		{"abcdefghi", "abcdefghi"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ShortModel(tt.input)
			if got != tt.want {
				t.Errorf("ShortModel(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestRelativeTime(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"seconds_ago", now.Add(-30 * time.Second).Format(time.RFC3339), "30s ago"},
		{"minutes_ago", now.Add(-5 * time.Minute).Format(time.RFC3339), "5m ago"},
		{"hours_ago", now.Add(-3 * time.Hour).Format(time.RFC3339), "3h ago"},
		{"days_ago", now.Add(-48 * time.Hour).Format(time.RFC3339), "2d ago"},
		{"invalid", "not-a-date", "not-a-date"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RelativeTime(tt.input)
			if got != tt.want {
				t.Errorf("RelativeTime(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
