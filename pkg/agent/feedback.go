package agent

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/agnath18K/lumo/pkg/config"
	"github.com/agnath18K/lumo/pkg/utils"
)

// Feedback handles user interaction and feedback
type Feedback struct {
	config *config.Config
	reader *bufio.Reader
}

// NewFeedback creates a new feedback instance
func NewFeedback(cfg *config.Config) *Feedback {
	return &Feedback{
		config: cfg,
		reader: bufio.NewReader(os.Stdin),
	}
}

// DisplayPlan shows the plan to the user
func (f *Feedback) DisplayPlan(plan *Plan) {
	fmt.Println("\n📋 " + plan.Task.Description)
	fmt.Println("───────────────────────────────────────────────")

	if plan.Description != "" {
		fmt.Printf("➤ %s\n\n", plan.Description)
	}

	for i, step := range plan.Steps {
		criticalMark := ""
		if step.IsCritical {
			criticalMark = " ⚠️"
		}

		// Add a separator between steps except for the first one
		if i > 0 {
			fmt.Println()
		}

		fmt.Printf("%d. %s%s\n", step.ID, step.Command, criticalMark)
		fmt.Printf("   %s\n", step.Description)
	}
}

// ConfirmExecution asks the user to confirm execution
func (f *Feedback) ConfirmExecution() (bool, error) {
	fmt.Println("\n🧐 I'm about to unleash these commands on your system...")
	fmt.Println("Don't worry, I've checked them twice, but you should too!")
	fmt.Println("Remember: with great commands comes great responsibility! 🦸")
	fmt.Print("\nDo you want to execute this plan? (y/n): ")
	response, err := f.reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("failed to read input: %w", err)
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes", nil
}

// DisplayStepStart shows that a step is starting
func (f *Feedback) DisplayStepStart(step *Step) {
	fmt.Printf("\n▶️ [%d] %s\n", step.ID, step.Command)
}

// DisplayStepResult shows the result of a step
func (f *Feedback) DisplayStepResult(step *Step) {
	result := step.Result

	if result.Success {
		fmt.Printf("✅ [%d] Completed in %s\n", step.ID, utils.FormatDuration(result.Duration))
	} else {
		fmt.Printf("❌ [%d] Failed in %s: %v\n", step.ID, utils.FormatDuration(result.Duration), result.Error)
	}

	// Display output if not empty, but limit it to avoid overwhelming the user
	if result.Output != "" {
		// Truncate long output
		output := result.Output
		maxLines := 5 // Reduced from 10 to 5 for more compact display
		lines := strings.Split(output, "\n")

		if len(lines) > maxLines {
			output = strings.Join(lines[:maxLines], "\n") + "\n... (" + fmt.Sprintf("%d", len(lines)-maxLines) + " more lines)"
		}

		// Add a subtle border around the output
		fmt.Println("┌─ Output ─────────────────────────────────")
		fmt.Printf("│ %s\n", strings.ReplaceAll(output, "\n", "\n│ "))
		fmt.Println("└──────────────────────────────────────────")
	}
}

// DisplaySummary shows a summary of the execution
func (f *Feedback) DisplaySummary(result *ExecutionResult) {
	// Count successful and failed steps
	successCount := 0
	failedCount := 0

	for _, step := range result.Plan.Steps {
		if step.Executed {
			if step.Result.Success {
				successCount++
			} else {
				failedCount++
			}
		}
	}

	fmt.Println("\n╭─────────────────────────────────────────╮")
	if result.Success {
		fmt.Printf("│ ✅ Task completed in %s              │\n", utils.FormatDuration(result.Duration))
		fmt.Printf("│ Steps: %d/%d successful                   │\n",
			successCount,
			successCount+failedCount)
	} else {
		fmt.Printf("│ ❌ Task failed in %s                 │\n", utils.FormatDuration(result.Duration))
		fmt.Printf("│ Error: %-35s │\n", result.Message)
		fmt.Printf("│ Steps: %d/%d successful                   │\n",
			successCount,
			successCount+failedCount)
	}
	fmt.Println("╰─────────────────────────────────────────╯")
}

// InteractiveREPL provides an interactive REPL for plan customization and execution
func (f *Feedback) InteractiveREPL(ctx context.Context, plan *Plan, executor *Executor) (*ExecutionResult, error) {
	var result *ExecutionResult

	// Main REPL loop
	for {
		// Display the current plan
		f.DisplayPlan(plan)

		// Display REPL options in a more compact and beautiful format
		fmt.Println()
		fmt.Println("╭─ Commands ──────────────────────────────────╮")
		fmt.Println("│ run                refine		           │")
		fmt.Println("│ add <cmd>          edit <num>               │")
		fmt.Println("│ delete <num>       move <num> <pos>         │")
		fmt.Println("│ exit               help                     │")
		fmt.Println("╰─────────────────────────────────────────────╯")

		// Get user input with a simple prompt
		fmt.Print("\nlumo> ")
		input, err := f.reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("failed to read input: %w", err)
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		// Parse the command
		parts := strings.SplitN(input, " ", 2)
		cmd := strings.ToLower(parts[0])
		args := ""
		if len(parts) > 1 {
			args = strings.TrimSpace(parts[1])
		}

		// Process the command
		switch cmd {
		case "run":
			// Execute the plan
			result, err = executor.ExecutePlan(ctx, plan, f)
			if err != nil {
				return nil, err
			}

			// Display the summary
			f.DisplaySummary(result)

			// Return the result without asking to edit
			return result, nil

		case "refine":
			// Refine the plan using natural language
			var modificationRequest string
			if args != "" {
				// Use the provided prompt directly
				modificationRequest = args
			} else {
				// Ask for input if no args provided
				fmt.Println("\n💬 Enter your refinement request in natural language:")
				fmt.Print("> ")
				var err error
				modificationRequest, err = f.reader.ReadString('\n')
				if err != nil {
					fmt.Printf("❌ Error reading input: %v\n", err)
					continue
				}
				modificationRequest = strings.TrimSpace(modificationRequest)
			}

			if modificationRequest == "" {
				fmt.Println("❌ Error: Empty refinement request")
				continue
			}

			// Get the AI client from the executor
			aiClient := executor.GetAIClient()
			if aiClient == nil {
				fmt.Println("❌ Error: AI client not available")
				continue
			}

			fmt.Println("🔄 Processing your request...")

			// Create a prompt for the AI to modify the plan
			var planText strings.Builder
			planText.WriteString(fmt.Sprintf("Current plan for task: %s\n\n", plan.Task.Description))
			planText.WriteString(fmt.Sprintf("Approach: %s\n\n", plan.Description))
			planText.WriteString("Steps:\n")

			for _, step := range plan.Steps {
				criticalMark := ""
				if step.IsCritical {
					criticalMark = " (critical)"
				}
				planText.WriteString(fmt.Sprintf("%d. %s%s\n", step.ID, step.Command, criticalMark))
				planText.WriteString(fmt.Sprintf("   %s\n\n", step.Description))
			}

			prompt := fmt.Sprintf(`You are an AI assistant helping to modify a shell command execution plan.

Current Plan:
%s

User's modification request: "%s"

Please modify the plan according to the user's request. Your response must be a valid JSON object with the following structure:
{
  "description": "Overall approach description",
  "steps": [
    {
      "id": 1,
      "command": "exact shell command",
      "description": "what this command does",
      "isCritical": true/false
    },
    ...
  ]
}

Do not include any text before or after the JSON object. The response must be parseable as JSON.
Ensure all commands are safe to execute and won't cause data loss or system damage.
Use relative paths when possible and avoid commands that require sudo.
Limit the plan to at most %d steps.
`, planText.String(), modificationRequest, executor.GetConfig().AgentMaxSteps)

			// Get response from AI
			response, err := aiClient.GetCompletion(ctx, prompt)
			if err != nil {
				fmt.Printf("❌ Error getting AI completion: %v\n", err)
				continue
			}

			// Extract JSON from the response
			jsonStart := -1
			jsonEnd := -1

			// Find the start of the JSON object
			for i := 0; i < len(response); i++ {
				if response[i] == '{' {
					jsonStart = i
					break
				}
			}

			// Find the end of the JSON object
			if jsonStart >= 0 {
				braceCount := 1
				for i := jsonStart + 1; i < len(response); i++ {
					if response[i] == '{' {
						braceCount++
					} else if response[i] == '}' {
						braceCount--
						if braceCount == 0 {
							jsonEnd = i + 1
							break
						}
					}
				}
			}

			if jsonStart < 0 || jsonEnd < 0 {
				fmt.Println("❌ Error: Could not extract valid JSON from AI response")
				continue
			}

			// Parse the JSON
			jsonStr := response[jsonStart:jsonEnd]
			var planData struct {
				Description string `json:"description"`
				Steps       []struct {
					ID          int    `json:"id"`
					Command     string `json:"command"`
					Description string `json:"description"`
					IsCritical  bool   `json:"isCritical"`
				} `json:"steps"`
			}

			if err := json.Unmarshal([]byte(jsonStr), &planData); err != nil {
				fmt.Printf("❌ Error parsing JSON: %v\n", err)
				continue
			}

			// Update the plan
			plan.Description = planData.Description

			// Create new steps
			newSteps := make([]*Step, 0, len(planData.Steps))
			for _, stepData := range planData.Steps {
				newSteps = append(newSteps, &Step{
					ID:          stepData.ID,
					Command:     stepData.Command,
					Description: stepData.Description,
					IsCritical:  stepData.IsCritical,
				})
			}

			// Replace the steps
			plan.Steps = newSteps

			fmt.Println("✅ Plan modified successfully!")

		case "add":
			if args == "" {
				fmt.Println("❌ Error: Command required")
				continue
			}

			// Add a new step
			f.addStep(plan, args)

		case "edit":
			if args == "" {
				fmt.Println("❌ Error: Step number required")
				continue
			}

			// Parse the step number
			stepNum, err := strconv.Atoi(args)
			if err != nil {
				fmt.Println("❌ Error: Invalid step number")
				continue
			}

			// Edit the step
			f.editStep(plan, stepNum)

		case "delete":
			if args == "" {
				fmt.Println("❌ Error: Step number required")
				continue
			}

			// Parse the step number
			stepNum, err := strconv.Atoi(args)
			if err != nil {
				fmt.Println("❌ Error: Invalid step number")
				continue
			}

			// Delete the step
			f.deleteStep(plan, stepNum)

		case "move":
			// Parse the arguments
			moveParts := strings.SplitN(args, " ", 2)
			if len(moveParts) != 2 {
				fmt.Println("❌ Error: Both source and destination positions required")
				continue
			}

			// Parse the step numbers
			srcNum, err := strconv.Atoi(strings.TrimSpace(moveParts[0]))
			if err != nil {
				fmt.Println("❌ Error: Invalid source step number")
				continue
			}

			destNum, err := strconv.Atoi(strings.TrimSpace(moveParts[1]))
			if err != nil {
				fmt.Println("❌ Error: Invalid destination position")
				continue
			}

			// Move the step
			f.moveStep(plan, srcNum, destNum)

		case "exit":
			// Exit without executing
			return nil, nil

		case "help":
			// Display detailed help
			fmt.Println("\n📚 REPL Mode Help:")
			fmt.Println("  run                  - Execute the current plan")
			fmt.Println("  refine <prompt>      - Refine the plan using natural language")
			fmt.Println("  add <command>        - Add a new step to the plan")
			fmt.Println("  edit <num>           - Edit a step in the plan")
			fmt.Println("  delete <num>         - Delete a step from the plan")
			fmt.Println("  move <num> <pos>     - Move a step to a new position")
			fmt.Println("  exit                 - Exit without executing")
			fmt.Println("  help                 - Show this help message")
			continue

		default:
			fmt.Println("❌ Error: Unknown command. Type 'help' for available commands.")
		}
	}
}

// addStep adds a new step to the plan
func (f *Feedback) addStep(plan *Plan, command string) {
	// Get the description
	fmt.Print("Enter description for this step: ")
	description, err := f.reader.ReadString('\n')
	if err != nil {
		fmt.Printf("❌ Error reading description: %v\n", err)
		return
	}
	description = strings.TrimSpace(description)

	// Ask if the step is critical
	fmt.Print("Is this step critical? (y/n): ")
	criticalInput, err := f.reader.ReadString('\n')
	if err != nil {
		fmt.Printf("❌ Error reading input: %v\n", err)
		return
	}
	criticalInput = strings.TrimSpace(strings.ToLower(criticalInput))
	isCritical := criticalInput == "y" || criticalInput == "yes"

	// Create the new step
	newStep := &Step{
		ID:          len(plan.Steps) + 1,
		Command:     command,
		Description: description,
		IsCritical:  isCritical,
	}

	// Add the step to the plan
	plan.Steps = append(plan.Steps, newStep)

	// Renumber steps
	f.renumberSteps(plan)

	fmt.Println("✅ Step added successfully")
}

// editStep edits a step in the plan
func (f *Feedback) editStep(plan *Plan, stepNum int) {
	// Check if the step exists
	if stepNum < 1 || stepNum > len(plan.Steps) {
		fmt.Println("❌ Error: Step number out of range")
		return
	}

	// Get the step
	step := plan.Steps[stepNum-1]

	// Get the new command
	fmt.Printf("Current command: %s\n", step.Command)
	fmt.Print("Enter new command (leave empty to keep current): ")
	command, err := f.reader.ReadString('\n')
	if err != nil {
		fmt.Printf("❌ Error reading command: %v\n", err)
		return
	}
	command = strings.TrimSpace(command)
	if command != "" {
		step.Command = command
	}

	// Get the new description
	fmt.Printf("Current description: %s\n", step.Description)
	fmt.Print("Enter new description (leave empty to keep current): ")
	description, err := f.reader.ReadString('\n')
	if err != nil {
		fmt.Printf("❌ Error reading description: %v\n", err)
		return
	}
	description = strings.TrimSpace(description)
	if description != "" {
		step.Description = description
	}

	// Ask if the step is critical
	fmt.Printf("Current critical status: %v\n", step.IsCritical)
	fmt.Print("Is this step critical? (y/n/leave empty to keep current): ")
	criticalInput, err := f.reader.ReadString('\n')
	if err != nil {
		fmt.Printf("❌ Error reading input: %v\n", err)
		return
	}
	criticalInput = strings.TrimSpace(strings.ToLower(criticalInput))
	if criticalInput != "" {
		step.IsCritical = criticalInput == "y" || criticalInput == "yes"
	}

	fmt.Println("✅ Step updated successfully")
}

// deleteStep deletes a step from the plan
func (f *Feedback) deleteStep(plan *Plan, stepNum int) {
	// Check if the step exists
	if stepNum < 1 || stepNum > len(plan.Steps) {
		fmt.Println("❌ Error: Step number out of range")
		return
	}

	// Remove the step
	plan.Steps = append(plan.Steps[:stepNum-1], plan.Steps[stepNum:]...)

	// Renumber steps
	f.renumberSteps(plan)

	fmt.Println("✅ Step deleted successfully")
}

// moveStep moves a step to a new position
func (f *Feedback) moveStep(plan *Plan, srcNum, destNum int) {
	// Check if the step exists
	if srcNum < 1 || srcNum > len(plan.Steps) {
		fmt.Println("❌ Error: Source step number out of range")
		return
	}

	// Check if the destination is valid
	if destNum < 1 || destNum > len(plan.Steps) {
		fmt.Println("❌ Error: Destination position out of range")
		return
	}

	// Get the step to move
	step := plan.Steps[srcNum-1]

	// Remove the step from its current position
	plan.Steps = append(plan.Steps[:srcNum-1], plan.Steps[srcNum:]...)

	// Insert the step at the new position
	newSteps := make([]*Step, 0, len(plan.Steps)+1)
	newSteps = append(newSteps, plan.Steps[:destNum-1]...)
	newSteps = append(newSteps, step)
	newSteps = append(newSteps, plan.Steps[destNum-1:]...)
	plan.Steps = newSteps

	// Renumber steps
	f.renumberSteps(plan)

	fmt.Println("✅ Step moved successfully")
}

// renumberSteps renumbers the steps in the plan
func (f *Feedback) renumberSteps(plan *Plan) {
	for i, step := range plan.Steps {
		step.ID = i + 1
	}
}
