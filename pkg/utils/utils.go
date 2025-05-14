package utils

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// FormatDuration formats a duration in a human-readable format
func FormatDuration(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%d µs", d.Microseconds())
	} else if d < time.Second {
		return fmt.Sprintf("%d ms", d.Milliseconds())
	} else if d < time.Minute {
		return fmt.Sprintf("%.2f s", d.Seconds())
	} else {
		return fmt.Sprintf("%d m %d s", int(d.Minutes()), int(d.Seconds())%60)
	}
}

// TruncateString truncates a string to the specified length
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// IsTerminal returns true if the file descriptor is a terminal
func IsTerminal(fd *os.File) bool {
	fileInfo, err := fd.Stat()
	if err != nil {
		return false
	}
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

// SplitCommandArgs splits a command string into command and arguments
func SplitCommandArgs(cmd string) (string, []string) {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return "", nil
	}
	return parts[0], parts[1:]
}

// ExpandPath expands a path with ~ to the user's home directory
func ExpandPath(path string) (string, error) {
	if !strings.HasPrefix(path, "~") {
		return path, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return strings.Replace(path, "~", homeDir, 1), nil
}

// FormatTimeAgo formats a time as a human-readable "time ago" string
func FormatTimeAgo(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	if diff < time.Minute {
		return "just now"
	} else if diff < time.Hour {
		minutes := int(diff.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	} else if diff < 24*time.Hour {
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	} else if diff < 48*time.Hour {
		return "yesterday"
	} else if diff < 7*24*time.Hour {
		days := int(diff.Hours() / 24)
		return fmt.Sprintf("%d days ago", days)
	} else if diff < 30*24*time.Hour {
		weeks := int(diff.Hours() / 24 / 7)
		if weeks == 1 {
			return "1 week ago"
		}
		return fmt.Sprintf("%d weeks ago", weeks)
	} else if diff < 365*24*time.Hour {
		months := int(diff.Hours() / 24 / 30)
		if months == 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", months)
	} else {
		years := int(diff.Hours() / 24 / 365)
		if years == 1 {
			return "1 year ago"
		}
		return fmt.Sprintf("%d years ago", years)
	}
}

// CleanMarkdown removes markdown formatting from a string for cleaner terminal output
func CleanMarkdown(text string) string {
	// Get terminal width for proper code block formatting
	termWidth := GetTerminalWidth()
	maxCodeWidth := termWidth - 6 // Account for borders and padding
	if maxCodeWidth < 40 {
		maxCodeWidth = 40
	}
	if maxCodeWidth > 100 {
		maxCodeWidth = 100
	}

	// Replace code blocks with a cleaner format
	codeBlockRegex := regexp.MustCompile("```(?:bash|sh)?\n((?s).*?)\n```")
	text = codeBlockRegex.ReplaceAllStringFunc(text, func(match string) string {
		// Extract the code content
		submatch := codeBlockRegex.FindStringSubmatch(match)
		if len(submatch) < 2 {
			return match
		}
		code := submatch[1]

		// Format the code block with a simple border
		lines := strings.Split(code, "\n")

		// Calculate the maximum line length, but cap it
		maxLength := 0
		for _, line := range lines {
			if len(line) > maxLength {
				maxLength = len(line)
			}
		}

		// Cap the maximum length to prevent overflow
		if maxLength > maxCodeWidth {
			maxLength = maxCodeWidth
		}

		// Use the same box style as the main response for consistency
		result := "\n╭" + strings.Repeat("─", maxLength+2) + "╮\n"
		for _, line := range lines {
			// Truncate line if it's too long
			displayLine := line
			if len(displayLine) > maxLength {
				displayLine = displayLine[:maxLength-3] + "..."
			}
			result += "│ " + displayLine + strings.Repeat(" ", maxLength-len(displayLine)) + " │\n"
		}
		result += "╰" + strings.Repeat("─", maxLength+2) + "╯\n"

		return result
	})

	// Remove bold/italic formatting
	text = regexp.MustCompile(`\*\*(.*?)\*\*`).ReplaceAllString(text, "$1")
	text = regexp.MustCompile(`\*(.*?)\*`).ReplaceAllString(text, "$1")

	// Remove inline code formatting
	text = regexp.MustCompile("`([^`]*)`").ReplaceAllString(text, "$1")

	// Remove bullet points and replace with dashes
	text = regexp.MustCompile(`(?m)^\s*\*\s+(.*)$`).ReplaceAllString(text, "- $1")

	return text
}

// FormatWithBox formats text with a box around it
func FormatWithBox(text string, title string) string {
	// Get terminal width
	termWidth := GetTerminalWidth()

	// Set maximum width for the box content (accounting for borders and padding)
	// Use a more conservative width to prevent overflow issues
	maxBoxWidth := termWidth - 6 // 2 chars for borders, 4 for padding and safety margin

	// Set reasonable minimum and maximum widths
	if maxBoxWidth < 40 {
		maxBoxWidth = 40 // Minimum reasonable width
	}
	if maxBoxWidth > 100 {
		maxBoxWidth = 100 // Maximum reasonable width to prevent overflow
	}

	// Process the text to wrap long lines
	var processedLines []string
	for _, line := range strings.Split(text, "\n") {
		// If line is empty, add it as is
		if strings.TrimSpace(line) == "" {
			processedLines = append(processedLines, "")
			continue
		}

		// If line is shorter than max width, add it as is
		if getDisplayWidth(line) <= maxBoxWidth {
			processedLines = append(processedLines, line)
			continue
		}

		// Wrap long lines
		words := strings.Fields(line)
		var currentLine string

		for _, word := range words {
			// If adding this word would make the line too long, start a new line
			if currentLine != "" && getDisplayWidth(currentLine)+getDisplayWidth(word)+1 > maxBoxWidth {
				processedLines = append(processedLines, currentLine)
				currentLine = word
			} else if currentLine == "" {
				currentLine = word
			} else {
				currentLine += " " + word
			}
		}

		// Add the last line if not empty
		if currentLine != "" {
			processedLines = append(processedLines, currentLine)
		}
	}

	// Find the maximum line display width after processing
	contentWidth := 0
	for _, line := range processedLines {
		lineWidth := getDisplayWidth(line)
		if lineWidth > contentWidth {
			contentWidth = lineWidth
		}
	}

	// Ensure minimum width and cap at maximum width
	if contentWidth < 40 {
		contentWidth = 40 // Minimum width
	}
	if contentWidth > maxBoxWidth {
		contentWidth = maxBoxWidth
	}

	// Build the box
	var sb strings.Builder

	// Add the top border with title if provided
	if title != "" {
		// Ensure the title has consistent spacing
		formattedTitle := " " + title + " "
		sb.WriteString("╭" + PadCenter(formattedTitle, contentWidth+2, "─") + "╮\n")
	} else {
		sb.WriteString("╭" + strings.Repeat("─", contentWidth+2) + "╮\n")
	}

	// Add the content
	for _, line := range processedLines {
		if strings.TrimSpace(line) == "" {
			sb.WriteString("│" + strings.Repeat(" ", contentWidth+2) + "│\n")
		} else {
			// Ensure line doesn't exceed contentWidth by truncating if necessary
			displayLine := line
			if getDisplayWidth(displayLine) > contentWidth {
				// Truncate by display width, not byte length
				truncated := ""
				currentWidth := 0
				for _, r := range displayLine {
					charWidth := 1
					if r > 0x1F00 {
						charWidth = 2
					}
					if currentWidth+charWidth > contentWidth-3 {
						break
					}
					truncated += string(r)
					currentWidth += charWidth
				}
				displayLine = truncated + "..."
			}
			sb.WriteString("│ " + PadRight(displayLine, contentWidth) + " │\n")
		}
	}

	// Add the bottom border
	sb.WriteString("╰" + strings.Repeat("─", contentWidth+2) + "╯\n")

	return sb.String()
}

// GetTerminalWidth returns the width of the terminal
func GetTerminalWidth() int {
	// Default width if we can't determine the actual width
	defaultWidth := 80

	// Try to get the terminal width using stty
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		// Fallback to environment variable
		if colsStr := os.Getenv("COLUMNS"); colsStr != "" {
			if cols, err := strconv.Atoi(colsStr); err == nil && cols > 0 {
				return cols
			}
		}
		return defaultWidth
	}

	// Parse the output (format: "rows cols")
	parts := strings.Fields(string(out))
	if len(parts) != 2 {
		return defaultWidth
	}

	// Convert the second part (cols) to an integer
	width, err := strconv.Atoi(parts[1])
	if err != nil {
		return defaultWidth
	}

	// Sanity check - if width is unreasonably small or large, use default
	if width < 20 || width > 500 {
		return defaultWidth
	}

	return width
}

// PadRight pads a string to the right to reach the specified length
func PadRight(s string, length int) string {
	displayWidth := getDisplayWidth(s)
	if displayWidth >= length {
		return s
	}
	return s + strings.Repeat(" ", length-displayWidth)
}

// getDisplayWidth returns the display width of a string, accounting for emojis and other wide characters
func getDisplayWidth(s string) int {
	width := 0
	for _, r := range s {
		// Check if the rune is an emoji or other wide character
		if r > 0x1F00 { // Most emojis and wide characters are above this range
			width += 2 // Emojis and wide characters typically have a display width of 2
		} else {
			width += 1 // Regular ASCII characters have a display width of 1
		}
	}
	return width
}

// CheckInternetConnectivity checks if there is an active internet connection
// by attempting to connect to a reliable host (Google's DNS server)
func CheckInternetConnectivity() bool {
	// Try to connect to Google's DNS server with a short timeout
	client := &http.Client{
		Timeout: 3 * time.Second,
	}
	_, err := client.Get("https://8.8.8.8:443")
	if err != nil {
		// Try another reliable host (Cloudflare's DNS)
		_, err = client.Get("https://1.1.1.1:443")
		if err != nil {
			return false
		}
	}
	return true
}

// FormatOfflineWarning formats the offline warning message in a humorous way without a box
func FormatOfflineWarning(provider string, ollamaAvailable bool, isAgentMode bool) string {
	var sb strings.Builder

	sb.WriteString("⚠️  Oops! Looks like you're offline in the digital wilderness! 🏕️\n\n")
	sb.WriteString(fmt.Sprintf("Lumo is currently set to use %s, which needs internet to work.\n", provider))

	if ollamaAvailable {
		sb.WriteString("\nGood news though! You have Ollama installed locally. 🎉\n")
		sb.WriteString("Switch to it with this magic spell:\n")
		sb.WriteString("  config:provider set ollama\n\n")

		if isAgentMode {
			sb.WriteString("This will let you use Lumo's agent mode offline with local models.\n")
		} else {
			sb.WriteString("This will let you use Lumo offline with local models.\n")
		}
		sb.WriteString("No internet? No problem! 💪\n")
	} else {
		sb.WriteString("\nWant to use Lumo offline? Install Ollama:\n")
		sb.WriteString("  https://ollama.com/download\n\n")
		sb.WriteString("Then switch to it with:\n")
		sb.WriteString("  config:provider set ollama\n\n")
		sb.WriteString("Your terminal assistant will work even when the internet doesn't! 🧙‍♂️\n")
	}

	return sb.String()
}

// PadCenter centers a string within a field of the specified length
func PadCenter(s string, length int, padChar string) string {
	displayWidth := getDisplayWidth(s)
	if displayWidth >= length {
		return s
	}

	// Calculate padding to ensure perfect centering
	totalPadding := length - displayWidth
	leftPad := totalPadding / 2
	rightPad := totalPadding - leftPad

	return strings.Repeat(padChar, leftPad) + s + strings.Repeat(padChar, rightPad)
}
