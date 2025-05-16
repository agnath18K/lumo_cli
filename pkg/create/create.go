package create

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/agnath18/lumo/pkg/ai"
)

// Generator handles project creation
type Generator struct {
	aiClient ai.Client
}

// NewGenerator creates a new project generator
func NewGenerator(aiClient ai.Client) *Generator {
	return &Generator{
		aiClient: aiClient,
	}
}

// Execute processes a project creation command
func (g *Generator) Execute(query string) (string, error) {
	// If no query is provided, show help
	if query == "" {
		return g.showHelp(), nil
	}

	// Parse the query to determine project type
	projectType, framework, options, err := g.parseQuery(query)
	if err != nil {
		return "", err
	}

	// Generate the project
	return g.generateProject(projectType, framework, options)
}

// parseQuery analyzes the natural language query to determine project details
func (g *Generator) parseQuery(query string) (string, string, map[string]string, error) {
	// Create a prompt for the AI to analyze the query
	prompt := fmt.Sprintf(`
You are a project creation assistant. Analyze the following query and extract the following information:
1. Project type/framework (e.g., Flutter, React, Next.js)
2. State management approach (e.g., Bloc, Provider, Riverpod for Flutter)
3. Any other specific requirements or options

Query: %s

Respond in the following JSON format:
{
  "projectType": "flutter|react|nextjs|etc",
  "framework": "bloc|provider|riverpod|redux|etc",
  "options": {
    "name": "project_name",
    "additionalFeatures": ["feature1", "feature2"]
  }
}

Only include fields that you can confidently determine from the query. Use snake_case for project names.
`, query)

	// Get response from AI
	response, err := g.aiClient.Query(prompt)
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to analyze query: %w", err)
	}

	// Extract JSON from response (this is a simplified approach)
	// In a real implementation, you would parse the JSON properly
	projectType := extractValue(response, "projectType")
	framework := extractValue(response, "framework")

	// Extract options
	options := make(map[string]string)

	// Extract project name
	name := extractValue(response, "name")
	if name != "" {
		options["name"] = name
	} else {
		// Default project name based on project type
		switch strings.ToLower(projectType) {
		case "flutter":
			options["name"] = "my_flutter_app"
		case "react":
			options["name"] = "my-react-app"
		case "nextjs":
			options["name"] = "my-nextjs-app"
		default:
			options["name"] = "my-app"
		}
	}

	// If project type is not specified, default to Flutter
	if projectType == "" {
		projectType = "flutter"
	}

	return projectType, framework, options, nil
}

// generateProject creates a project based on the specified type and framework
func (g *Generator) generateProject(projectType, framework string, options map[string]string) (string, error) {
	// Convert project type to lowercase for case-insensitive comparison
	projectType = strings.ToLower(projectType)

	// Generate the project based on type
	switch projectType {
	case "flutter":
		return generateFlutterProject(framework, options)
	case "nextjs":
		return generateNextJSProject(framework, options)
	case "react":
		return generateReactProject(framework, options)
	case "fastapi", "flask", "python":
		return generatePythonProject(framework, options)
	// Add more project types here as needed
	default:
		return "", fmt.Errorf("unsupported project type: %s", projectType)
	}
}

// showHelp returns help information for the create command
func (g *Generator) showHelp() string {
	return `
╭─────────────────── Lumo Project Creator ───────────────────╮
│                                                            │
│  Create new projects with natural language descriptions.   │
│                                                            │
│  Usage:                                                    │
│    lumo create:"<project description>"                     │
│                                                            │
│  Examples:                                                 │
│    lumo create:"Flutter app with bloc architecture"        │
│    lumo create:"Flutter app with provider state management"│
│    lumo create:"Next.js app with Redux"                    │
│    lumo create:"Next.js project with Context API"          │
│    lumo create:"React app with MobX state management"      │
│    lumo create:"React project with Recoil"                 │
│    lumo create:"FastAPI project with SQLAlchemy"           │
│    lumo create:"Flask web application"                     │
│                                                            │
│  Supported Frameworks:                                     │
│    • Flutter (with Bloc, Provider, Riverpod)               │
│    • Next.js (with Redux, Context API, Zustand)            │
│    • React (with Redux, Context API, MobX, Recoil)         │
│    • Python (FastAPI, Flask)                               │
│                                                            │
╰────────────────────────────────────────────────────────────╯
`
}

// extractValue is a simple helper to extract values from the AI response
// In a real implementation, you would use proper JSON parsing
func extractValue(response, key string) string {
	// This is a very simplified approach - in production code, use proper JSON parsing
	pattern := fmt.Sprintf(`"%s":\s*"([^"]+)"`, key)
	matches := regexp.MustCompile(pattern).FindStringSubmatch(response)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}
