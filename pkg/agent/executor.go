package agent

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"

	"github.com/agnath18K/lumo/pkg/ai"
	"github.com/agnath18K/lumo/pkg/config"
)

// Executor handles the execution of plans
type Executor struct {
	config   *config.Config
	aiClient ai.Client
}

// NewExecutor creates a new executor instance
func NewExecutor(cfg *config.Config, aiClient ai.Client) *Executor {
	return &Executor{
		config:   cfg,
		aiClient: aiClient,
	}
}

// ExecutePlan executes all steps in a plan using a single inline terminal session
func (e *Executor) ExecutePlan(ctx context.Context, plan *Plan, feedback *Feedback) (*ExecutionResult, error) {
	result := &ExecutionResult{
		Plan:      plan,
		StartTime: time.Now(),
		Success:   true,
	}

	// Start a single bash session for the entire plan
	cmd := exec.CommandContext(ctx, "bash")

	// Create pipes for stdin, stdout, and stderr
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the bash process
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start bash process: %w", err)
	}

	// Create a combined reader for stdout and stderr
	outputReader := io.MultiReader(stdout, stderr)
	outputScanner := bufio.NewScanner(outputReader)

	// Execute each step in the plan
	for _, step := range plan.Steps {
		// Update the current step
		feedback.DisplayStepStart(step)

		// Execute the step in the inline terminal
		stepResult, err := e.ExecuteStepInline(ctx, step, stdin, outputScanner)
		if err != nil {
			// Try to terminate the bash process
			cmd.Process.Kill()
			return nil, fmt.Errorf("failed to execute step %d: %w", step.ID, err)
		}

		// Update the step with the result
		step.Result = stepResult
		step.Executed = true

		// Display the step result
		feedback.DisplayStepResult(step)

		// Check if the step failed
		if !stepResult.Success {
			// If the step is critical, stop execution
			if step.IsCritical {
				result.Success = false
				result.Message = fmt.Sprintf("Critical step %d failed: %v", step.ID, stepResult.Error)
				break
			}
			// For non-critical steps, mark the overall result as failed but continue execution
			result.Success = false
			result.Message = fmt.Sprintf("Step %d failed: %v", step.ID, stepResult.Error)
		}
	}

	// Send exit command to bash
	fmt.Fprintln(stdin, "exit")
	stdin.Close()

	// Wait for the bash process to complete
	cmd.Wait()

	// Set the end time and duration
	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	// Set success message if successful
	if result.Success {
		result.Message = "All steps completed successfully"
	}

	return result, nil
}

// ExecuteStepInline executes a single step in the inline terminal
func (e *Executor) ExecuteStepInline(ctx context.Context, step *Step, stdin io.Writer, scanner *bufio.Scanner) (*StepResult, error) {
	result := &StepResult{
		StartTime: time.Now(),
	}

	// Check if the command is empty
	if strings.TrimSpace(step.Command) == "" {
		result.Success = false
		result.Error = fmt.Errorf("empty command")
		result.EndTime = time.Now()
		result.Duration = result.EndTime.Sub(result.StartTime)
		return result, nil
	}

	// Add a unique marker to identify the end of command output
	marker := fmt.Sprintf("LUMO_CMD_COMPLETE_%d", time.Now().UnixNano())

	// Send the command followed by an echo of the marker
	fmt.Fprintf(stdin, "%s\necho $? > /tmp/lumo_exit_code\necho %s\n", step.Command, marker)

	// Collect output until we see the marker
	var outputBuilder strings.Builder
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == marker {
			break
		}
		outputBuilder.WriteString(line)
		outputBuilder.WriteString("\n")
	}

	// Send command to get the exit code
	fmt.Fprintf(stdin, "cat /tmp/lumo_exit_code\necho %s\n", marker)

	// Read the exit code
	var exitCode string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == marker {
			break
		}
		exitCode = strings.TrimSpace(line)
	}

	// Set the end time and duration
	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	// Set the output
	result.Output = outputBuilder.String()

	// Check for errors based on exit code
	if exitCode != "0" {
		result.Success = false
		result.Error = fmt.Errorf("exit status %s", exitCode)
		return result, nil
	}

	// Set success
	result.Success = true

	return result, nil
}

// ExecuteStep executes a single step in the plan (legacy method, kept for compatibility)
func (e *Executor) ExecuteStep(ctx context.Context, step *Step) (*StepResult, error) {
	result := &StepResult{
		StartTime: time.Now(),
	}

	// Check if the command is empty
	if strings.TrimSpace(step.Command) == "" {
		result.Success = false
		result.Error = fmt.Errorf("empty command")
		result.EndTime = time.Now()
		result.Duration = result.EndTime.Sub(result.StartTime)
		return result, nil
	}

	// Create the command using bash to handle pipes, redirects, etc.
	cmd := exec.CommandContext(ctx, "bash", "-c", step.Command)

	// Capture the output
	output, err := cmd.CombinedOutput()

	// Set the end time and duration
	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	// Set the output
	result.Output = string(output)

	// Check for errors
	if err != nil {
		result.Success = false
		result.Error = err
		return result, nil
	}

	// Set success
	result.Success = true

	return result, nil
}

// GetAIClient returns the AI client
func (e *Executor) GetAIClient() ai.Client {
	return e.aiClient
}

// GetConfig returns the configuration
func (e *Executor) GetConfig() *config.Config {
	return e.config
}
