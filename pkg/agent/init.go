package agent

import (
	"github.com/agnath18K/lumo/pkg/config"
	"github.com/agnath18K/lumo/pkg/executor"
)

// Initialize initializes the agent and registers it with the executor
func Initialize(cfg *config.Config, exec *executor.Executor) *Agent {
	// Get the AI client from the executor
	aiClient := exec.GetAIClient()

	// Create a new agent
	agent := &Agent{
		config:   cfg,
		planner:  NewPlanner(cfg, aiClient),
		executor: NewExecutor(cfg, aiClient),
		feedback: NewFeedback(cfg),
		state: &AgentState{
			Status: StatusIdle,
		},
		aiClient: aiClient,
	}

	// Register the agent with the executor
	exec.SetAgent(agent)

	return agent
}
