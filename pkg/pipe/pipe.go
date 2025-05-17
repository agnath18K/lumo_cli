package pipe

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/agnath18K/lumo/pkg/ai"
)

// Processor handles processing of piped input
type Processor struct {
	aiClient ai.Client
}

// NewProcessor creates a new pipe processor
func NewProcessor(aiClient ai.Client) *Processor {
	return &Processor{
		aiClient: aiClient,
	}
}

// ProcessInput reads input from a reader and processes it
func (p *Processor) ProcessInput(reader io.Reader) (string, error) {
	// Read all input from the reader
	content, err := readAllInput(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read piped input: %w", err)
	}

	// If content is empty, return an error
	if strings.TrimSpace(content) == "" {
		return "", fmt.Errorf("empty input")
	}

	// Process the content using AI
	return p.analyzeContent(content)
}

// readAllInput reads all input from a reader
func readAllInput(reader io.Reader) (string, error) {
	scanner := bufio.NewScanner(reader)
	var builder strings.Builder

	// Read all lines
	for scanner.Scan() {
		line := scanner.Text()
		builder.WriteString(line)
		builder.WriteString("\n")
	}

	// Check for errors
	if err := scanner.Err(); err != nil {
		return "", err
	}

	return builder.String(), nil
}

// analyzeContent uses AI to analyze the content
func (p *Processor) analyzeContent(content string) (string, error) {
	// Create a prompt for the AI
	prompt := fmt.Sprintf(`
Analyze the following text and provide a clear explanation of its contents.
If it contains code or commands, explain what they do in a user-friendly manner.
If it's data or text, summarize the key points and structure.
Be concise but thorough in your explanation.

TEXT TO ANALYZE:
%s

Your analysis should include:
1. Type of content (code, commands, data, text, etc.)
2. Purpose or function of the content
3. Key components or structure
4. Any potential issues or considerations
`, content)

	// Get response from AI
	response, err := p.aiClient.Query(prompt)
	if err != nil {
		return "", fmt.Errorf("failed to analyze content: %w", err)
	}

	return response, nil
}
