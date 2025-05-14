package agent

import (
	"time"
)

// Status represents the current status of the agent
type Status string

const (
	// StatusIdle represents an idle agent
	StatusIdle Status = "idle"
	// StatusPlanning represents an agent in planning phase
	StatusPlanning Status = "planning"
	// StatusExecuting represents an agent executing commands
	StatusExecuting Status = "executing"
	// StatusCompleted represents an agent that has completed its task
	StatusCompleted Status = "completed"
	// StatusFailed represents an agent that has failed its task
	StatusFailed Status = "failed"
)

// Task represents a user task to be executed
type Task struct {
	// Description is the natural language description of the task
	Description string
	// CreatedAt is the time when the task was created
	CreatedAt time.Time
}

// Plan represents a sequence of steps to accomplish a task
type Plan struct {
	// Task is the task to be accomplished
	Task *Task
	// Steps is the sequence of steps to accomplish the task
	Steps []*Step
	// CreatedAt is the time when the plan was created
	CreatedAt time.Time
	// Description is a brief description of the overall approach
	Description string
}

// Step represents a single command to be executed
type Step struct {
	// ID is the step number
	ID int
	// Command is the shell command to execute
	Command string
	// Description is a brief description of what the command does
	Description string
	// IsCritical indicates whether the step is critical for the task
	IsCritical bool
	// Executed indicates whether the step has been executed
	Executed bool
	// Result is the result of executing the step
	Result *StepResult
}

// StepResult represents the result of executing a step
type StepResult struct {
	// Success indicates whether the step was successful
	Success bool
	// Output is the command output
	Output string
	// Error is any error that occurred
	Error error
	// StartTime is when the step execution started
	StartTime time.Time
	// EndTime is when the step execution ended
	EndTime time.Time
	// Duration is how long the step took to execute
	Duration time.Duration
}

// ExecutionResult represents the overall result of executing a plan
type ExecutionResult struct {
	// Success indicates whether the execution was successful
	Success bool
	// Message is a message describing the result
	Message string
	// Plan is the executed plan
	Plan *Plan
	// StartTime is when the execution started
	StartTime time.Time
	// EndTime is when the execution ended
	EndTime time.Time
	// Duration is how long the execution took
	Duration time.Duration
}

// AgentState represents the current state of the agent
type AgentState struct {
	// Status is the current status of the agent
	Status Status
	// CurrentTask is the task being executed
	CurrentTask *Task
	// CurrentPlan is the plan being executed
	CurrentPlan *Plan
	// CurrentStep is the step being executed
	CurrentStep *Step
}
