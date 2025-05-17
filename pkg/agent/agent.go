package agent

import (
	"context"
	"fmt"
	"time"

	"github.com/agnath18K/lumo/pkg/ai"
	"github.com/agnath18K/lumo/pkg/config"
	"github.com/agnath18K/lumo/pkg/executor"
)

// Agent represents the auto command executor
type Agent struct {
	config   *config.Config
	planner  *Planner
	executor *Executor
	feedback *Feedback
	state    *AgentState
	aiClient ai.Client
}

// Execute processes a task and executes the necessary commands
func (a *Agent) Execute(ctx context.Context, taskDescription string) (*executor.Result, error) {
	// Check if agent mode is enabled
	if !a.config.EnableAgentMode {
		return &executor.Result{
			IsError: true,
			Output:  "Agent mode is disabled. Enable it in the configuration file.",
		}, nil
	}

	// Create a new task
	task := &Task{
		Description: taskDescription,
		CreatedAt:   time.Now(),
	}

	// Update agent state
	a.state.Status = StatusPlanning
	a.state.CurrentTask = task

	// Generate a plan
	plan, err := a.planner.CreatePlan(ctx, task)
	if err != nil {
		return &executor.Result{
			IsError: true,
			Output:  fmt.Sprintf("Failed to create plan: %v", err),
		}, nil
	}

	// Update agent state
	a.state.CurrentPlan = plan

	// Display warning about agent mode
	fmt.Println("\nAGENT MODE WARNING:")
	fmt.Println("Agent mode will execute shell commands on your behalf.")
	fmt.Println("Always review the plan carefully before confirming execution!")
	fmt.Println("Commands may have unintended consequences if not properly reviewed.")

	// Check if we should use interactive REPL mode
	var result *ExecutionResult
	var executionErr error

	if a.config.EnableAgentREPL {
		// Use interactive REPL mode
		result, executionErr = a.feedback.InteractiveREPL(ctx, plan, a.executor)
		if executionErr != nil {
			return &executor.Result{
				IsError: true,
				Output:  fmt.Sprintf("Failed during interactive REPL: %v", executionErr),
			}, nil
		}

		// Check if the user exited without executing
		if result == nil {
			return &executor.Result{
				IsError: false,
				Output:  "Execution cancelled by user.",
			}, nil
		}
	} else {
		// Use traditional confirmation mode
		a.feedback.DisplayPlan(plan)

		// Confirm execution with the user if required
		if a.config.AgentConfirmBeforeExecution {
			confirmed, err := a.feedback.ConfirmExecution()
			if err != nil {
				return &executor.Result{
					IsError: true,
					Output:  fmt.Sprintf("Failed to confirm execution: %v", err),
				}, nil
			}

			if !confirmed {
				return &executor.Result{
					IsError: false,
					Output:  "Execution cancelled by user.",
				}, nil
			}
		}

		// Update agent state
		a.state.Status = StatusExecuting

		// Execute the plan
		result, executionErr = a.executor.ExecutePlan(ctx, plan, a.feedback)
		if executionErr != nil {
			return &executor.Result{
				IsError: true,
				Output:  fmt.Sprintf("Failed to execute plan: %v", executionErr),
			}, nil
		}
	}

	// Update agent state
	if result.Success {
		a.state.Status = StatusCompleted
	} else {
		a.state.Status = StatusFailed
	}

	// Provide final summary
	a.feedback.DisplaySummary(result)

	// Return the result
	return &executor.Result{
		IsError: !result.Success,
		Output:  result.Message,
	}, nil
}
