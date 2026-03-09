package components

import (
	"fmt"
	"strings"

	"github.com/prs/ccut/internal/ui/format"
	"github.com/prs/ccut/internal/ui/styles"
)

// Sparkline renders a horizontal bar chart with labels.
func Sparkline(items []SparkItem, width int) string {
	if len(items) == 0 {
		return styles.Subtle.Render("  no data")
	}

	maxVal := 0
	maxLabelW := 0
	for _, item := range items {
		if item.Value > maxVal {
			maxVal = item.Value
		}
		if len(item.Label) > maxLabelW {
			maxLabelW = len(item.Label)
		}
	}
	if maxVal == 0 {
		maxVal = 1
	}

	barWidth := width - maxLabelW - 12
	if barWidth < 10 {
		barWidth = 10
	}

	var lines []string
	for _, item := range items {
		label := fmt.Sprintf("%*s", maxLabelW, item.Label)
		filled := item.Value * barWidth / maxVal
		bar := styles.BarFull.Render(strings.Repeat("█", filled)) +
			styles.BarEmpty.Render(strings.Repeat("░", barWidth-filled))
		val := styles.Subtle.Render(fmt.Sprintf(" %s", format.Number(item.Value)))
		suffix := ""
		if item.Suffix != "" {
			suffix = " " + styles.Special.Render(item.Suffix)
		}
		lines = append(lines, fmt.Sprintf("  %s %s%s%s", styles.Subtle.Render(label), bar, val, suffix))
	}
	return strings.Join(lines, "\n")
}

type SparkItem struct {
	Label  string
	Value  int
	Suffix string
}
