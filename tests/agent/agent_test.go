package agent_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/agnath18/lumo/pkg/agent"
	"github.com/agnath18/lumo/pkg/ai"
	"github.com/agnath18/lumo/pkg/config"
	"github.com/agnath18/lumo/pkg/executor"
)

// MockAIClient is a mock implementation of the ai.Client interface for testing
type MockAIClient struct {
	response string
	err      error
}

func (m *MockAIClient) Query(prompt string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.response, nil
}

func (m *MockAIClient) GetCompletion(ctx context.Context, prompt string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.response, nil
}

// MockExecutor is a mock implementation of the executor.Executor for testing
type MockExecutor struct {
	config   *config.Config
	aiClient ai.Client
	agent    executor.AgentInterface
}

func NewMockExecutor(cfg *config.Config, aiClient ai.Client) *executor.Executor {
	// Use the actual NewExecutor function to create a properly initialized executor
	return executor.NewExecutor(cfg)
}

// TestAgentInitialization tests that the agent is initialized correctly
func TestAgentInitialization(t *testing.T) {
	// Create a config with agent mode enabled
	cfg := config.DefaultConfig()
	cfg.EnableAgentMode = true

	// Create a mock AI client
	mockClient := &MockAIClient{
		response: `{
			"plan": {
				"description": "Test plan",
				"steps": [
					{
						"command": "echo 'Hello, world!'",
						"description": "Print a greeting",
						"critical": true
					}
				]
			}
		}`,
	}

	// Create a mock executor
	mockExec := NewMockExecutor(cfg, mockClient)

	// Initialize the agent
	a := agent.Initialize(cfg, mockExec)

	// Check that the agent was created
	if a == nil {
		t.Fatal("Expected agent to be created, got nil")
	}
}

// TestAgentDisabled tests that the agent returns an error when disabled
func TestAgentDisabled(t *testing.T) {
	// Create a config with agent mode disabled
	cfg := config.DefaultConfig()
	cfg.EnableAgentMode = false

	// Create a mock AI client
	mockClient := &MockAIClient{
		response: "This response doesn't matter for this test",
	}

	// Create a mock executor
	mockExec := NewMockExecutor(cfg, mockClient)

	// Initialize the agent
	a := agent.Initialize(cfg, mockExec)

	// Check that the agent was created
	if a == nil {
		t.Fatal("Expected agent to be created, got nil")
	}

	// Create a context
	ctx := context.Background()

	// Execute a task
	result, err := a.Execute(ctx, "Test task")

	// Check that there was no error
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check that the result indicates an error
	if !result.IsError {
		t.Error("Expected IsError to be true when agent mode is disabled")
	}

	// Check that the error message mentions agent mode being disabled
	if !strings.Contains(result.Output, "Agent mode is disabled") {
		t.Errorf("Expected error message to mention agent mode being disabled, got: %s", result.Output)
	}
}

// TestPlanCreation tests the creation of a plan
func TestPlanCreation(t *testing.T) {
	// Create a config with agent mode enabled
	cfg := config.DefaultConfig()
	cfg.EnableAgentMode = true
	cfg.AgentMaxSteps = 5

	// Create a mock AI client with a valid JSON response
	mockClient := &MockAIClient{
		response: `{
			"description": "Test plan description",
			"steps": [
				{
					"id": 1,
					"command": "echo 'Step 1'",
					"description": "Print Step 1",
					"isCritical": true
				},
				{
					"id": 2,
					"command": "echo 'Step 2'",
					"description": "Print Step 2",
					"isCritical": false
				}
			]
		}`,
	}

	// Create a planner with the mock client
	planner := agent.NewPlanner(cfg, mockClient)

	// Create a task
	task := &agent.Task{
		Description: "Test task",
		CreatedAt:   time.Now(),
	}

	// Create a context
	ctx := context.Background()

	// Generate a plan
	plan, err := planner.CreatePlan(ctx, task)

	// Check that there was no error
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check that the plan was created
	if plan == nil {
		t.Fatal("Expected plan to be created, got nil")
	}

	// Check the plan description
	if plan.Description != "Test plan description" {
		t.Errorf("Expected plan description to be 'Test plan description', got '%s'", plan.Description)
	}

	// Check the number of steps
	if len(plan.Steps) != 2 {
		t.Errorf("Expected 2 steps, got %d", len(plan.Steps))
	}

	// Check the first step
	if plan.Steps[0].Command != "echo 'Step 1'" {
		t.Errorf("Expected first step command to be 'echo 'Step 1'', got '%s'", plan.Steps[0].Command)
	}
	if plan.Steps[0].Description != "Print Step 1" {
		t.Errorf("Expected first step description to be 'Print Step 1', got '%s'", plan.Steps[0].Description)
	}
	if !plan.Steps[0].IsCritical {
		t.Error("Expected first step to be critical")
	}

	// Check the second step
	if plan.Steps[1].Command != "echo 'Step 2'" {
		t.Errorf("Expected second step command to be 'echo 'Step 2'', got '%s'", plan.Steps[1].Command)
	}
	if plan.Steps[1].Description != "Print Step 2" {
		t.Errorf("Expected second step description to be 'Print Step 2', got '%s'", plan.Steps[1].Description)
	}
	if plan.Steps[1].IsCritical {
		t.Error("Expected second step to not be critical")
	}
}
