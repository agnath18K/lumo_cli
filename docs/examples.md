# Lumo Examples

This document provides comprehensive examples for all Lumo features and commands.

## Basic Usage

### Natural Language Queries

```bash
# Get help with finding files
lumo "How do I find files by content in Linux?"

# Get help with a specific command
lumo "How to use grep with regular expressions?"

# Get explanations for technical concepts
lumo "What's the difference between TCP and UDP?"

# Get help with system administration
lumo "How to check disk space usage?"

# Get help with Git commands
lumo "How to undo the last Git commit?"
```

## Agent Mode

### Basic Agent Tasks

```bash
# Create a backup of documents
lumo auto:create a backup of my documents folder

# Find and analyze large files
lumo agent:find all large files in the current directory and show their types

# Set up a development environment
lumo agent:set up a python virtual environment with flask and sqlalchemy

# Monitor system resources
lumo agent:monitor CPU and memory usage every 5 seconds

# Clean up temporary files
lumo agent:find and remove all temporary files older than 7 days
```

### Agent Mode REPL Commands

When in the Agent Mode REPL interface:

```
# Execute the current plan
run

# Modify the plan using natural language
refine add a step to compress the backup file

# Add a new step to the plan
add tar -czf backup.tar.gz backup/

# Edit a specific step in the plan
edit 3

# Remove a step from the plan
delete 2

# Reorder steps in the plan
move 4 2

# Show available commands
help

# Exit without executing
exit
```

## Chat Mode

### Starting Chat Mode

```bash
# Start interactive chat mode
lumo chat

# Ask a direct question in chat format
lumo chat:Tell me about Linux file permissions

# Ask a specific question
lumo ask:What is the capital of France?
```

### Chat Mode Commands

In the interactive Chat Mode:

```
# Send a message to the AI
How does public key encryption work?

# Show help information
help

# Clear conversation history
clear

# Display conversation history
history

# Start a new conversation
new

# List all conversations
list

# Switch to another conversation
switch 2

# Delete a conversation
delete 1

# Exit chat mode
exit
```

## System Commands

### Shell Commands

```bash
# Execute a shell command directly
lumo lumo:ls -la

# Execute another shell command
lumo shell:find . -name "*.go" -type f
```

### System Health

```bash
# Get a basic health report
lumo health

# Get a detailed system report
lumo system

# Check specific system components
lumo health:memory
lumo health:disk
lumo health:cpu
lumo health:network

# Generate a report with specific focus
lumo report:performance
lumo report:security
lumo report:storage
```

## Internet Speed Testing

```bash
# Run a complete speed test
lumo speed

# Test only download speed
lumo speed:download

# Test only upload speed
lumo speed:upload

# Test with natural language
lumo "check my internet speed"
lumo "how fast is my internet connection"
```

## Clipboard Operations

```bash
# Show current clipboard contents
lumo clipboard

# Copy text to clipboard
lumo clipboard "Hello World"

# Append text to existing clipboard content
lumo clipboard append "More text"

# Clear clipboard contents
lumo clipboard clear

# Copy piped content to clipboard
echo "This is some text" | lumo clipboard

# Append piped content to clipboard
cat file.txt | lumo clipboard append
```

## Project Creation

```bash
# Create a Flutter project with BLoC architecture
lumo create:"Flutter app with bloc architecture"

# Create a Flutter project with Provider state management
lumo create:"Flutter app with provider state management"

# Create a Flutter project with Riverpod
lumo create:"Flutter app with riverpod state management"

# Show help for the create command
lumo create
```

## Magic Commands

```bash
# Show a fun dance animation
lumo magic:dance

# List available magic commands
lumo magic:help
```

## Configuration Commands

```bash
# List available AI providers
lumo config:provider list

# Show current AI provider
lumo config:provider show

# Set AI provider
lumo config:provider set gemini
lumo config:provider set openai
lumo config:provider set ollama

# List available models for the current provider
lumo config:model list

# Show current model
lumo config:model show

# Set model for current provider
lumo config:model set gemini-2.0-flash-lite

# Show API key status
lumo config:key show

# Set API key for a provider
lumo config:key set gemini YOUR_API_KEY
lumo config:key set openai YOUR_API_KEY

# Remove API key for a provider
lumo config:key remove gemini

# Show current Ollama URL
lumo config:ollama show

# Set Ollama URL
lumo config:ollama set http://localhost:11434

# Test connection to Ollama server
lumo config:ollama test
```

## Pipe Support

```bash
# Analyze directory contents
ls -la | lumo

# Explain error logs
cat error.log | lumo

# Understand complex command output
ps aux | grep python | lumo

# Analyze system information
dmesg | lumo

# Get help with command output
ifconfig | lumo

# Analyze JSON data
cat data.json | lumo

# Analyze CSV data
cat data.csv | lumo
```

## Command-Line Options

```bash
# Display help information
lumo --help
lumo -h
lumo help

# Show version information
lumo --version
lumo -v
lumo version
```
