package system

import (
	"strings"
)

// padRight pads a string to the right to reach the specified width
func padRight(s string, width int) string {
	if len(s) > width {
		return s[:width-3] + "..."
	}
	return s + strings.Repeat(" ", width-len(s))
}

// padCenter centers a string within a field of the specified width
func padCenter(s string, width int, padChar string) string {
	if len(s) >= width {
		return s[:width]
	}
	
	leftPad := (width - len(s)) / 2
	rightPad := width - len(s) - leftPad
	
	return strings.Repeat(padChar, leftPad) + s + strings.Repeat(padChar, rightPad)
}

// truncateString truncates a string to the specified length
func truncateString(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen-3] + "..."
	}
	return s
}

// padOrTruncate ensures a string fits within the specified width
func padOrTruncate(s string, width int) string {
	if len(s) > width {
		return s[:width-3] + "..."
	}
	return s + strings.Repeat(" ", width-len(s))
}
