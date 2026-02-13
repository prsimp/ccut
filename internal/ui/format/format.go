package format

import (
	"fmt"
	"strings"
	"time"
)

// Number formats a number with comma separators.
func Number(n int) string {
	if n < 0 {
		return "-" + Number(-n)
	}
	s := fmt.Sprintf("%d", n)
	if len(s) <= 3 {
		return s
	}
	var result []byte
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, byte(c))
	}
	return string(result)
}

// Tokens formats token counts in human-readable form.
func Tokens(n int) string {
	switch {
	case n >= 1_000_000_000:
		return fmt.Sprintf("%.1fB", float64(n)/1_000_000_000)
	case n >= 1_000_000:
		return fmt.Sprintf("%.1fM", float64(n)/1_000_000)
	case n >= 1_000:
		return fmt.Sprintf("%.1fK", float64(n)/1_000)
	default:
		return fmt.Sprintf("%d", n)
	}
}

// Duration formats milliseconds into a human-readable duration.
func Duration(ms int64) string {
	d := time.Duration(ms) * time.Millisecond
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	if h > 0 {
		return fmt.Sprintf("%dh %dm", h, m)
	}
	return fmt.Sprintf("%dm", m)
}

// Date formats an ISO date string into a short form.
func Date(isoDate string) string {
	t, err := time.Parse(time.RFC3339, isoDate)
	if err != nil {
		t, err = time.Parse("2006-01-02", isoDate)
		if err != nil {
			return isoDate
		}
	}
	return t.Local().Format("Jan 02 15:04")
}

// DateShort returns just the date portion.
func DateShort(isoDate string) string {
	t, err := time.Parse("2006-01-02", isoDate)
	if err != nil {
		t, err = time.Parse(time.RFC3339, isoDate)
		if err != nil {
			return isoDate
		}
	}
	return t.Format("Jan 02")
}

// ShortModel returns a readable model name.
func ShortModel(model string) string {
	m := model
	// Strip date suffix like -20250929 (dash + 8 digits)
	if len(m) > 9 {
		suffix := m[len(m)-8:]
		allDigits := true
		for _, c := range suffix {
			if c < '0' || c > '9' {
				allDigits = false
				break
			}
		}
		if allDigits && m[len(m)-9] == '-' {
			m = m[:len(m)-9]
		}
	}
	m = strings.TrimPrefix(m, "claude-")
	return m
}

// RelativeTime returns a relative time string.
func RelativeTime(isoDate string) string {
	t, err := time.Parse(time.RFC3339, isoDate)
	if err != nil {
		return isoDate
	}
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return fmt.Sprintf("%ds ago", int(d.Seconds()))
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	}
}
