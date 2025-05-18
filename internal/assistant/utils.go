package assistant

import (
	"fmt"
	"strings"
)

// extractTarget extracts the target from the input
func extractTarget(input string, keywords []string) string {
	// Remove the keywords from the input
	target := input
	for _, keyword := range keywords {
		target = strings.ReplaceAll(target, keyword, "")
	}

	// Clean up the target
	target = strings.TrimSpace(target)
	target = strings.Trim(target, "\"'")

	// If the target is empty, return "current"
	if target == "" {
		return "current"
	}

	return target
}

// extractApplicationAndArgs extracts the application name and arguments from the input
func extractApplicationAndArgs(input string) (string, string) {
	fmt.Printf("DEBUG: Extracting application and args from: %s\n", input)

	// Create a copy of the original input for debugging
	originalInput := input

	// Remove keywords from the input
	cleaned := input

	// Replace keywords with spaces to avoid word concatenation
	for _, keyword := range []string{"launch", "open", "start", "app", "application", "program"} {
		cleaned = strings.ReplaceAll(cleaned, keyword, " ")
	}

	// Clean up the input
	cleaned = strings.TrimSpace(cleaned)
	cleaned = strings.Trim(cleaned, "\"'")

	// Remove multiple spaces
	for strings.Contains(cleaned, "  ") {
		cleaned = strings.ReplaceAll(cleaned, "  ", " ")
	}

	fmt.Printf("DEBUG: Cleaned input: %s\n", cleaned)

	// Split the input into application name and arguments
	parts := strings.SplitN(cleaned, " with ", 2)
	if len(parts) == 2 {
		appName := strings.TrimSpace(parts[0])
		args := strings.TrimSpace(parts[1])
		fmt.Printf("DEBUG: Found app with args (with): app=%s, args=%s\n", appName, args)
		return appName, args
	}

	parts = strings.SplitN(cleaned, " using ", 2)
	if len(parts) == 2 {
		appName := strings.TrimSpace(parts[0])
		args := strings.TrimSpace(parts[1])
		fmt.Printf("DEBUG: Found app with args (using): app=%s, args=%s\n", appName, args)
		return appName, args
	}

	parts = strings.SplitN(cleaned, " and ", 2)
	if len(parts) == 2 {
		appName := strings.TrimSpace(parts[0])
		args := strings.TrimSpace(parts[1])
		fmt.Printf("DEBUG: Found app with args (and): app=%s, args=%s\n", appName, args)
		return appName, args
	}

	// If no arguments are found, return the cleaned input as the application name
	fmt.Printf("DEBUG: No arguments found, app=%s\n", cleaned)

	// Special case for common applications
	if strings.Contains(originalInput, "terminal") {
		fmt.Printf("DEBUG: Special case: terminal application detected\n")
		return "gnome-terminal", ""
	}
	if strings.Contains(originalInput, "firefox") {
		fmt.Printf("DEBUG: Special case: firefox application detected\n")
		return "firefox", ""
	}
	if strings.Contains(originalInput, "chrome") {
		fmt.Printf("DEBUG: Special case: chrome application detected\n")
		return "google-chrome", ""
	}

	return cleaned, ""
}

// extractNotificationContent extracts the notification summary and body from the input
func extractNotificationContent(input string) (string, string) {
	// Remove keywords from the input
	cleaned := input
	for _, keyword := range []string{"send", "notification", "with", "message"} {
		cleaned = strings.ReplaceAll(cleaned, keyword, "")
	}

	// Clean up the input
	cleaned = strings.TrimSpace(cleaned)
	cleaned = strings.Trim(cleaned, "\"'")

	// Split the input into summary and body
	parts := strings.SplitN(cleaned, " and ", 2)
	if len(parts) == 2 {
		return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	}

	parts = strings.SplitN(cleaned, " with ", 2)
	if len(parts) == 2 {
		return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	}

	parts = strings.SplitN(cleaned, ": ", 2)
	if len(parts) == 2 {
		return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	}

	// If no body is found, return the cleaned input as the summary
	return cleaned, ""
}
