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
lumo talk:Tell me about Linux file permissions

# Ask a specific question without entering chat mode
lumo ask:What is the capital of France?
lumo ai:What is the capital of France?
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
# Execute a shell command (MUST use shell: prefix)
lumo shell:ls -la

# Execute another shell command
lumo shell:find . -name "*.go" -type f

# Note: Shell commands are ONLY executed when explicitly prefixed with "shell:"
# The following will NOT execute as a shell command, but will be processed as an AI query:
lumo ls -la
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

# Alternative syntax for health commands
lumo syshealth:memory
lumo syshealth:disk

# Generate a report with specific focus
lumo report:performance
lumo report:security
lumo report:storage

# Alternative syntax for report commands
lumo sysreport:performance
lumo sysreport:security
```

## Internet Speed Testing

```bash
# Run a complete speed test
lumo speed

# Test only download speed
lumo speed:download

# Test only upload speed
lumo speed:upload

# Alternative syntax for speed tests
lumo speedtest:download
lumo speed-test:upload

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

# Create a Next.js project with Redux
lumo create:"Next.js app with Redux state management"

# Create a Next.js project with Context API
lumo create:"Next.js project with Context API"

# Create a Next.js project with Zustand
lumo create:"Next.js application using Zustand for state"

# Create a basic Next.js project without specific state management
lumo create:"Simple Next.js project"

# Create a React project with Redux
lumo create:"React app with Redux state management"

# Create a React project with Context API
lumo create:"React project with Context API"

# Create a React project with MobX
lumo create:"React application using MobX for state"

# Create a React project with Recoil
lumo create:"React app with Recoil state management"

# Create a basic React project without specific state management
lumo create:"Simple React project"

# Create a FastAPI project
lumo create:"FastAPI project with SQLAlchemy"

# Create a FastAPI project with specific options
lumo create:"Create a FastAPI REST API for a blog"

# Create a Flask project
lumo create:"Flask web application"

# Create a Flask project with specific options
lumo create:"Create a Flask app with SQLAlchemy and authentication"

# Show help for the create command
lumo create
```

## Desktop Assistant

The desktop assistant allows you to control your desktop environment using natural language commands. It uses AI to understand complex commands and execute them.

```bash
# Close a specific window
lumo desktop:"close firefox window"

# Minimize a window
lumo desktop:"minimize terminal window"

# Maximize a window
lumo desktop:"maximize chrome window"

# List all open windows
lumo desktop:"list windows"

# Launch an application
lumo desktop:"launch terminal"

# List running applications
lumo desktop:"list applications"

# Lock the screen
lumo desktop:"lock screen"

# Send a notification
lumo desktop:"send notification Hello World with body This is a test"

# Control media playback
lumo desktop:"play media"
lumo desktop:"pause media"
lumo desktop:"stop media"
lumo desktop:"next track"
lumo desktop:"previous track"

# Change appearance settings (GNOME)
lumo desktop:"set dark mode on"
lumo desktop:"set light mode"
lumo desktop:"change desktop background to /path/to/image.jpg"
lumo desktop:"set GTK theme to Adwaita-dark"
lumo desktop:"change icon theme to Papirus"
lumo desktop:"get current theme"
lumo desktop:"show desktop background"

# Control sound settings (GNOME)
lumo desktop:"set volume to 50 percent"
lumo desktop:"increase volume to 75 percent"
lumo desktop:"mute the sound"
lumo desktop:"unmute the sound"
lumo desktop:"set microphone volume to 80 percent"
lumo desktop:"mute the microphone"
lumo desktop:"show all sound devices"
lumo desktop:"get current volume level"
lumo desktop:"set default sound device to alsa_output.pci-0000_00_1f.3.analog-stereo"

# Control connectivity settings (GNOME)
lumo desktop:"show all network devices"
lumo desktop:"turn on WiFi"
lumo desktop:"turn off WiFi"
lumo desktop:"check WiFi status"
lumo desktop:"enable Bluetooth"
lumo desktop:"disable Bluetooth"
lumo desktop:"check Bluetooth status"
lumo desktop:"turn on airplane mode"
lumo desktop:"turn off airplane mode"
lumo desktop:"check airplane mode status"
lumo desktop:"create a WiFi hotspot with name 'MyHotspot'"
lumo desktop:"create a WiFi hotspot with name 'MyHotspot' and password 'securepass'"
lumo desktop:"turn off WiFi hotspot"
lumo desktop:"check hotspot status"

# AI-powered natural language commands
lumo desktop:"I want to close all Firefox windows and then open a new terminal"
lumo desktop:"Could you please minimize all my windows and then lock my screen?"
lumo desktop:"First open Firefox, then maximize it, and finally play some music"
lumo desktop:"Switch to dark mode and set my background to night-sky.jpg"
lumo desktop:"Increase the volume to 80 percent and then play some music"
lumo desktop:"Turn off WiFi and Bluetooth to save battery"
lumo desktop:"Enable airplane mode and then set volume to 0"
lumo desktop:"Create a hotspot named 'LumoShare' with password 'lumo1234'"
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

## File Transfer with Connect

```bash
# Start a server to receive files
lumo connect --receive

# Start a server on a specific port
lumo connect --receive --port 9000

# Specify a custom download directory
lumo connect --receive --path ~/Downloads/transfers

# Connect to a peer to send/receive files
lumo connect 192.168.1.5

# Connect to a peer on a specific port
lumo connect 192.168.1.5:9000

# Connect to a peer with a custom download directory
lumo connect 192.168.1.5 --path ~/Downloads/transfers

# Show connect command help
lumo connect --help
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
