package utils

import (
	"fmt"
	"time"
)

func FormatDuration(d time.Duration) string {
	// Convert the duration to minutes and seconds
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60

	// Format the result as "Xm Ys"
	return fmt.Sprintf("%dm %ds", minutes, seconds)
}
