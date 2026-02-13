package components

import (
	"fmt"
	"strings"

	"github.com/prs/ccut/internal/ui/styles"
)

// Heatmap renders a 24-hour activity heatmap.
func Heatmap(hourCounts map[string]int) string {
	maxVal := 0
	for _, v := range hourCounts {
		if v > maxVal {
			maxVal = v
		}
	}
	if maxVal == 0 {
		maxVal = 1
	}

	var labels []string
	var blocks []string
	for h := 0; h < 24; h++ {
		key := fmt.Sprintf("%d", h)
		count := hourCounts[key]

		labels = append(labels, styles.Subtle.Render(fmt.Sprintf("%2d", h)))

		ratio := float64(count) / float64(maxVal)
		var block string
		switch {
		case count == 0:
			block = styles.HeatLow.Render("██")
		case ratio < 0.33:
			block = styles.HeatMedLow.Render("██")
		case ratio < 0.66:
			block = styles.HeatMed.Render("██")
		default:
			block = styles.HeatHigh.Render("██")
		}
		blocks = append(blocks, block)
	}

	return "  " + strings.Join(labels, " ") + "\n" +
		"  " + strings.Join(blocks, " ")
}
