<div align="center">
  <h1>🐦 Lumo - Your AI-Powered Terminal Assistant</h1>
</div>

<div align="center">
  <img src="https://img.shields.io/badge/version-1.0.1-blue.svg" alt="Version 1.0.1">
  <img src="https://img.shields.io/badge/go-%3E%3D1.22-00ADD8.svg" alt="Go Version">
  <img src="https://img.shields.io/badge/license-MIT-green.svg" alt="License MIT">
  <img src="https://github.com/agnath18K/lumo/workflows/CI/badge.svg" alt="CI Status">
  <img src="https://github.com/agnath18K/lumo/workflows/GoReleaser/badge.svg" alt="Release Status">
</div>

<p align="center">
  <b>Lumo is an intelligent CLI assistant that interprets natural language to help you navigate the terminal with ease.</b>
</p>

<div align="center">
  <a href="#-overview">Overview</a> •
  <a href="#-key-features">Key Features</a> •
  <a href="#-installation">Installation</a> •
  <a href="#-usage">Usage</a> •
  <a href="#%EF%B8%8F-configuration">Configuration</a> •
  <a href="#-technical-architecture">Architecture</a> •
  <a href="#%EF%B8%8F-development">Development</a> •
  <a href="#-license">License</a>
</div>

## 📖 Overview

Lumo is a powerful command-line interface (CLI) assistant designed to bridge the gap between natural language and terminal commands. By leveraging advanced AI models from Google Gemini and OpenAI, Lumo helps users find, understand, and execute terminal commands without memorizing complex syntax.

Whether you're a seasoned developer or a terminal novice, Lumo enhances your productivity by:

- Translating natural language queries into terminal commands
- Providing explanations for complex terminal concepts
- Executing sequences of commands to complete tasks
- Analyzing command outputs and system information
- Offering conversational assistance for general queries

Lumo is built with a focus on user experience, safety, and extensibility, making it a valuable addition to any developer's toolkit.


## 🌟 Key Features

### 🧠 Natural Language Command Processing
Lumo's core functionality is its ability to understand natural language queries and translate them into terminal commands.

- **Command Translation**: Ask questions in plain English and get the exact terminal commands you need
- **Contextual Understanding**: Lumo understands the intent behind your queries, even with ambiguous phrasing
- **Command Explanation**: Get detailed explanations of how commands work and what each parameter does
- **Intelligent Suggestions**: Receive command suggestions tailored to your specific needs and environment

### 🐦 Agent Mode (Auto Command Execution)
Agent Mode enables Lumo to function as an autonomous assistant that can plan and execute sequences of commands to complete complex tasks.

- **Task Planning**: Automatically generates a structured plan of shell commands to accomplish your task
- **Interactive REPL Interface**: Review, refine, and customize execution plans through a user-friendly interface
- **Step-by-Step Execution**: Execute commands sequentially with real-time feedback and progress tracking
- **Execution Safety**: Built-in safeguards to prevent potentially harmful operations, with user confirmation required
- **Performance Metrics**: Detailed timing and output reporting for each executed command

### 💬 Chat Mode
Engage in conversational interactions with Lumo for general assistance and information.

- **Contextual Conversations**: Maintain context throughout your conversation for more natural interactions
- **Specialized Knowledge**: Access AI knowledge about terminal commands, programming concepts, and general topics
- **Conversation Management**: Save, load, and manage multiple conversation threads
- **Custom System Instructions**: Separate system instructions optimize responses for general conversations

### 📊 System Health & Monitoring
Get insights into your system's performance and health with detailed reports.

- **Resource Monitoring**: Track CPU, memory, disk usage, and other critical system metrics
- **Performance Analysis**: Identify performance bottlenecks and resource-intensive processes
- **Health Recommendations**: Receive actionable recommendations for system optimization
- **Detailed Reports**: Generate comprehensive system reports for troubleshooting and documentation

### 🚀 Internet Speed Testing
Test and analyze your internet connection speed directly from the terminal.

- **Comprehensive Testing**: Measure download speed, upload speed, and latency
- **Beautiful Visualization**: View results with visual indicators and connection quality ratings
- **Targeted Testing**: Run specific tests for download or upload speeds
- **Connection Diagnostics**: Identify potential connectivity issues and bottlenecks

### 📋 Pipe Support & Output Analysis
Process and analyze command outputs directly through Lumo's pipe functionality.

- **Output Interpretation**: Get plain-language explanations of complex command outputs
- **Pattern Recognition**: Identify patterns, anomalies, and important information in large outputs
- **Data Extraction**: Extract meaningful information from logs, configuration files, and other text data
- **Format Conversion**: Transform outputs into more readable and understandable formats

### 🔄 Multiple AI Provider Integration
Flexibility to choose between different AI providers based on your preferences and requirements.

- **Provider Options**: Seamless support for Google Gemini, OpenAI, and Ollama models
- **Model Selection**: Choose different models based on your performance and accuracy needs
- **API Key Management**: Secure storage and management of API keys
- **Local AI Support**: Use Ollama for local AI inference without internet connectivity
- **Fallback Mechanisms**: Automatic fallback options if one provider is unavailable

## 🚀 Installation

### System Requirements

Before installing Lumo, ensure your system meets the following requirements:

| Requirement | Details |
|-------------|---------|
| **Operating System** | Linux (Ubuntu, Debian, CentOS, etc.) or macOS. Windows users can use WSL (Windows Subsystem for Linux) |
| **Go Version** | Go 1.22 or higher |
| **Disk Space** | Approximately 20MB for the binary and dependencies |
| **Internet Connection** | Required for API access to AI providers |
| **API Keys** | Google Gemini or OpenAI API key (obtained during setup) |

### Installation Methods

#### Method 1: Quick Install from Source

The recommended way to install Lumo is by building from source:

```bash
# Clone the Lumo repository from GitHub
git clone https://github.com/agnath18K/lumo.git

# Navigate into the cloned repository
cd lumo

# Build the binary
make build

# Install the binary to /usr/bin for global access
sudo cp build/lumo /usr/bin/lumo

# Verify the installation
lumo --version
```

#### Method 2: Using Pre-built Binaries

If you prefer not to build from source, you can download pre-built binaries from the [releases page](https://github.com/agnath18K/lumo/releases):

```bash
# Download the latest release for your platform (example for Linux amd64)
curl -L https://github.com/agnath18K/lumo/releases/latest/download/lumo-linux-amd64 -o lumo

# Make the binary executable
chmod +x lumo

# Move to a directory in your PATH
sudo mv lumo /usr/bin/lumo

# Verify the installation
lumo --version
```

### First-Time Setup and Configuration

When you first run Lumo, it will guide you through an interactive setup process:

1. **Configuration Initialization**:
   - Lumo automatically creates the necessary configuration directory at `~/.config/lumo/`
   - Default configuration files are generated if they don't exist

2. **AI Provider Selection**:
   - You'll be prompted to choose your preferred AI provider:
     - **Google Gemini**: Generally faster with good performance for terminal tasks
     - **OpenAI**: May provide more detailed explanations for complex queries
     - **Ollama**: Local AI inference without requiring internet connectivity

3. **API Key Configuration**:
   - You'll need to provide API keys for your chosen provider(s):
     - For **Gemini**: Get your API key from [Google AI Studio](https://aistudio.google.com/apikey)
     - For **OpenAI**: Get your API key from [OpenAI Platform](https://platform.openai.com/api-keys)
     - For **Ollama**: Provide the URL where Ollama is running (default: http://localhost:11434)

4. **Key Verification and Storage**:
   - Lumo will verify your API keys to ensure they're valid
   - Keys are securely stored in your configuration file
   - You can update keys later by editing `~/.config/lumo/config.json`

5. **Feature Enablement**:
   - By default, most features are enabled
   - You can customize which features are active in the configuration file

## 🔍 Usage Guide

Lumo offers multiple ways to interact with your terminal through natural language. This section provides detailed examples and explanations for each mode of operation.

### Command-Line Options

| Option | Alternative | Description |
|--------|-------------|-------------|
| `--help`, `-h` | `help` | Display help information |
| `--version`, `-v` | `version` | Show version information |

### Command Prefixes

Lumo supports various command prefixes for different functionalities:

| Prefix | Description | Example |
|--------|-------------|---------|
| `auto:`, `agent:` | Execute a sequence of commands as an agent | `lumo auto:create a backup` |
| `chat:`, `talk:` | Start or continue a conversation | `lumo chat:tell me about Linux` |
| `lumo:`, `shell:` | Execute a shell command directly | `lumo lumo:ls -la` |
| `health:`, `syshealth:` | Check system health | `lumo health:memory` |
| `report:`, `sysreport:` | Generate system reports | `lumo report:performance` |
| `speed:` | Run internet speed tests | `lumo speed:download` |
| `config:` | Configure Lumo settings | `lumo config:provider list` |
| `magic:` | Run fun magic commands | `lumo magic:dance` |
| `clipboard` | Clipboard operations | `lumo clipboard "text"` |

### Basic Usage: Natural Language Queries

The simplest way to use Lumo is by asking questions in natural language:

```bash
# Get help with command-line options
lumo --help

# Ask how to perform a specific task
lumo "How do I find large files in Linux?"

# Get help with a specific command
lumo "How to use grep with multiple patterns?"

# Get explanations for technical concepts
lumo "What's the difference between TCP and UDP?"

# Show version information
lumo --version
```

Lumo will process your query and provide:
- The exact command(s) you need
- An explanation of how the command works
- Additional context and examples
- Warnings about potential issues

### Agent Mode: Automated Task Execution

Agent Mode allows Lumo to execute sequences of commands to complete complex tasks. This mode is activated using the `auto:` or `agent:` prefix:

```bash
# Create a backup of documents
lumo auto:create a backup of my documents folder

# Find and analyze large files
lumo agent:find all large files in the current directory and show their types

# Set up a development environment
lumo agent:set up a python virtual environment with flask and sqlalchemy

# Monitor system resources
lumo agent:monitor CPU and memory usage every 5 seconds
```

#### Agent Mode Workflow

When using Agent Mode, Lumo follows a structured workflow:

1. **Plan Generation**: Lumo analyzes your request and generates a detailed plan of shell commands
2. **Plan Review**: The plan is displayed for your review in an interactive REPL interface
3. **Plan Refinement**: You can modify the plan using natural language with the `refine` command
4. **Execution**: After your approval, Lumo executes the commands sequentially
5. **Real-time Feedback**: Each command's output and status is displayed during execution
6. **Summary**: A comprehensive summary of the execution is provided upon completion

#### Agent Mode REPL Commands

In the Agent Mode REPL interface, you can use the following commands:

| Command | Description |
|---------|-------------|
| `run` | Execute the current plan |
| `refine <prompt>` | Modify the plan using natural language |
| `add <command>` | Add a new step to the plan |
| `edit <num>` | Edit a specific step in the plan |
| `delete <num>` | Remove a step from the plan |
| `move <num> <pos>` | Reorder steps in the plan |
| `help` | Show available commands |
| `exit` | Exit without executing |

### Chat Mode: Conversational Assistance

Chat Mode provides a conversational interface for more interactive assistance:

```bash
# Start interactive chat mode
lumo chat

# Ask a direct question in chat format
lumo chat:Tell me about Linux file permissions

# Ask a specific question
lumo ask:What is the capital of France?
```

#### Chat Mode Commands

In the interactive Chat Mode, you can use these commands:

| Command | Description |
|---------|-------------|
| `<message>` | Send a message to the AI |
| `help` | Show help information |
| `clear` | Clear conversation history |
| `history` | Display conversation history |
| `new` | Start a new conversation |
| `list` | List all conversations |
| `switch <id>` | Switch to another conversation |
| `delete <id>` | Delete a conversation |
| `exit`, `quit` | Exit chat mode |

### Pipe Support: Analyzing Command Output

Lumo can analyze the output of other commands when used with pipes:

```bash
# Analyze directory contents
ls -la | lumo

# Explain error logs
cat error.log | lumo

# Understand complex command output
ps aux | grep python | lumo

# Analyze system information
dmesg | lumo
```

When processing piped input, Lumo will:
- Identify the type of content (logs, directory listings, etc.)
- Provide a summary of the key information
- Highlight important patterns or anomalies
- Offer explanations for technical terms or error messages

### System Health and Monitoring

Lumo provides built-in commands for checking system health and generating reports:

```bash
# Get a basic health report
lumo health

# Get a detailed system report
lumo system

# Check specific system components
lumo health:memory
lumo health:disk

# Generate a report with specific focus
lumo report:performance
lumo report:security
```

System health reports include information about:
- CPU usage and performance
- Memory utilization
- Disk space and I/O statistics
- Network connectivity
- Running processes
- System load and uptime

### Internet Speed Testing

Lumo can test your internet connection speed directly from the terminal:

```bash
# Run a complete speed test (download, upload, and latency)
lumo speed

# Test only download speed
lumo speed:download

# Test only upload speed
lumo speed:upload

# Use natural language to request a speed test
lumo "check my internet speed"
lumo "how fast is my internet connection"
```

Speed test results include:
- Download speed in Mbps
- Upload speed in Mbps
- Latency/ping in milliseconds
- ISP information
- Connection quality rating
- Visual indicators for easy interpretation

### Clipboard Operations

Lumo provides built-in clipboard functionality to view, copy, append, and clear clipboard content:

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

The clipboard feature is useful for:
- Saving command outputs for later use
- Building up complex commands piece by piece
- Transferring text between terminal sessions
- Quickly accessing frequently used text snippets

### Magic Commands

Lumo includes fun "magic" commands for entertainment and utility:

```bash
# Show a fun dance animation
lumo magic:dance

# List available magic commands
lumo magic:help
```

The `magic:dance` command displays a randomly selected dance animation in your terminal, providing a moment of fun during your terminal work.

## ⚙️ Configuration

Lumo provides extensive configuration options to customize its behavior according to your preferences and requirements.

### Configuration File

Lumo stores its configuration in `~/.config/lumo/config.json`. This file is automatically created during the first-time setup but can be manually edited at any time:

```json
{
  "ai_provider": "gemini",
  "gemini_api_key": "your-gemini-api-key",
  "gemini_model": "gemini-2.0-flash-lite",
  "openai_api_key": "your-openai-api-key",
  "openai_model": "gpt-4o",
  "ollama_url": "http://localhost:11434",
  "ollama_model": "llama3",
  "enable_agent_mode": true,
  "agent_confirm_before_execution": true,
  "agent_max_steps": 10,
  "agent_safety_level": "medium",
  "enable_agent_repl": true,
  "enable_shell_in_interactive": false,
  "enable_chat_repl": true,
  "enable_pipe_processing": true,
  "enable_system_health": true,
  "enable_system_report": true,
  "enable_speed_test": true,
  "speed_test_timeout": 30,
  "max_history_size": 1000,
  "enable_logging": true,
  "debug": false
}
```

### Configuration Options

| Option | Type | Description | Default |
|--------|------|-------------|---------|
| **AI Provider Settings** |
| `ai_provider` | string | The AI provider to use ("gemini", "openai", or "ollama") | "gemini" |
| `gemini_api_key` | string | Your Google Gemini API key | "" |
| `gemini_model` | string | The Gemini model to use | "gemini-2.0-flash-lite" |
| `openai_api_key` | string | Your OpenAI API key | "" |
| `openai_model` | string | The OpenAI model to use | "gpt-4o" |
| `ollama_url` | string | URL for Ollama server | "http://localhost:11434" |
| `ollama_model` | string | The Ollama model to use | "llama3" |
| **Agent Mode Settings** |
| `enable_agent_mode` | boolean | Enable or disable Agent Mode | true |
| `agent_confirm_before_execution` | boolean | Require confirmation before executing commands | true |
| `agent_max_steps` | integer | Maximum number of steps in an agent plan | 10 |
| `agent_safety_level` | string | Safety level for command execution ("low", "medium", "high") | "medium" |
| `enable_agent_repl` | boolean | Enable interactive REPL interface for Agent Mode | true |
| **Feature Toggles** |
| `enable_shell_in_interactive` | boolean | Allow shell commands in interactive mode | false |
| `enable_chat_repl` | boolean | Enable interactive chat REPL | true |
| `enable_pipe_processing` | boolean | Enable processing of piped input | true |
| `enable_system_health` | boolean | Enable system health checks | true |
| `enable_system_report` | boolean | Enable system reports | true |
| `enable_speed_test` | boolean | Enable internet speed testing | true |
| `speed_test_timeout` | integer | Timeout in seconds for speed tests | 30 |
| **General Settings** |
| `max_history_size` | integer | Maximum number of commands to store in history | 1000 |
| `enable_logging` | boolean | Enable command logging | true |
| `debug` | boolean | Enable debug mode with additional output | false |

### Configuration Commands

Lumo provides command-line options to view and modify configuration settings without editing the config file directly:

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

### Environment Variables

You can also configure Lumo using environment variables, which take precedence over the configuration file:

| Environment Variable | Corresponding Config Option |
|----------------------|----------------------------|
| `LUMO_AI_PROVIDER` | `ai_provider` |
| `LUMO_GEMINI_API_KEY` | `gemini_api_key` |
| `LUMO_OPENAI_API_KEY` | `openai_api_key` |
| `LUMO_OLLAMA_URL` | `ollama_url` |
| `LUMO_ENABLE_AGENT_MODE` | `enable_agent_mode` |
| `LUMO_DEBUG` | `debug` |

## 🔧 Technical Architecture

Lumo is built with a modular architecture in Go, designed for extensibility, maintainability, and performance.

### Core Components

<div align="center">
  <img src="https://github.com/agnath18K/lumo/raw/main/docs/images/architecture.png" alt="Lumo Architecture" width="700">
</div>

#### Command Processing Pipeline

1. **Input Handling**: Processes user input from command line arguments or pipes
2. **NLP Parser**: Interprets natural language queries and determines command intent
3. **Command Executor**: Routes commands to appropriate handlers based on type
4. **Output Formatter**: Formats results for terminal display

#### Key Modules

| Module | Description |
|--------|-------------|
| **NLP Parser** | Interprets natural language queries and extracts intent and parameters |
| **AI Client** | Provides a unified interface to Gemini, OpenAI, and Ollama APIs with proper error handling |
| **Agent System** | Plans and executes command sequences with safety checks and user feedback |
| **Terminal Interface** | Handles user interaction with ANSI color support and progress indicators |
| **System Monitoring** | Collects and analyzes system information using gopsutil |
| **Pipe Processor** | Analyzes piped input and generates insights |
| **Configuration Manager** | Handles loading, validation, and saving of user configuration |
| **Chat Manager** | Manages conversational context and history |
| **Clipboard** | Handles clipboard operations for viewing, copying, appending, and clearing content |
| **Magic** | Provides fun and entertaining terminal animations and utilities |
| **Speed Test** | Measures and reports internet connection speed and quality |

### Data Flow

1. User input is received via command line or pipe
2. The NLP parser determines the command type and intent
3. The command is routed to the appropriate handler
4. For AI queries, the request is sent to the selected AI provider
5. Results are processed and formatted for display
6. Output is presented to the user with appropriate formatting

## 🛠️ Development

This section provides information for developers who want to contribute to Lumo or customize it for their own needs.

### Development Environment Setup

To set up a development environment for Lumo:

```bash
# Clone the repository
git clone https://github.com/agnath18/lumo.git
cd lumo

# Install Go (if not already installed)
# For Ubuntu/Debian:
# sudo apt-get install golang-go

# For macOS with Homebrew:
# brew install go

# Verify Go installation
go version  # Should be 1.22 or higher

# Install dependencies
go mod download

# Build the development version
make build
```

### Project Structure

```
lumo/
├── cmd/
│   └── lumo/           # Main application entry point
├── pkg/
│   ├── agent/          # Agent mode implementation
│   ├── ai/             # AI provider clients
│   ├── chat/           # Chat functionality
│   ├── clipboard/      # Clipboard operations
│   ├── config/         # Configuration management
│   ├── executor/       # Command execution
│   ├── magic/          # Fun terminal animations
│   ├── nlp/            # Natural language processing
│   ├── pipe/           # Pipe processing
│   ├── setup/          # First-time setup
│   ├── speedtest/      # Internet speed testing
│   ├── system/         # System health and monitoring
│   ├── terminal/       # Terminal UI
│   ├── utils/          # Utility functions
│   └── version/        # Version information
├── build/              # Build artifacts
├── docs/               # Documentation
├── tests/              # Test files
├── go.mod              # Go module definition
├── go.sum              # Go module checksums
└── Makefile            # Build automation
```

### Build and Test

Lumo includes a comprehensive Makefile for common development tasks:

```bash
# Build the binary
make build

# Run all tests
make test

# Run specific tests
go test ./pkg/agent/...

# Install the binary
make install

# Clean build artifacts
make clean

# Show version information
make version

# Get help on available commands
make help
```

### Testing

Lumo uses Go's built-in testing framework. Tests are organized by package and functionality:

```bash
# Run all tests
make test

# Run tests with coverage report
go test -cover ./...

# Run tests for a specific package
go test ./pkg/agent/...

# Run a specific test
go test ./pkg/agent/... -run TestCreatePlan
```

### Contributing

Contributions to Lumo are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests to ensure everything works (`make test`)
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

### Continuous Integration and Deployment

Lumo uses GitHub Actions for continuous integration and deployment:

- **CI Workflow**: Automatically runs tests and builds the binary for each push to the main branch and pull requests
- **Release Workflow**: Automatically builds binaries for multiple platforms and creates a GitHub release when a new tag is pushed
- **GoReleaser**: Used for streamlined release management, packaging, and distribution

To create a new release:

```bash
# Update version in pkg/version/version.go
# Commit the changes
git commit -am "Release version X.Y.Z"

# Tag the release
git tag -a vX.Y.Z -m "Version X.Y.Z"

# Push the changes and tags
git push && git push --tags
```

The GitHub Actions workflow will automatically:
1. Build binaries for Linux, macOS (Intel and Apple Silicon)
2. Create Debian packages
3. Generate checksums
4. Create a GitHub release with all artifacts attached

## 📜 License

Lumo is released under the MIT License.

```
MIT License

Copyright (c) 2024 agnath18

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

## 🙏 Acknowledgements

Lumo stands on the shoulders of these amazing technologies and projects:

### Core Technologies
- [Go Programming Language](https://golang.org/) - The foundation of Lumo's performance and reliability
- [Google Gemini API](https://ai.google.dev/) - Provides advanced natural language understanding capabilities
- [OpenAI API](https://openai.com/api/) - Powers conversational AI features

### Libraries and Dependencies
- [gopsutil](https://github.com/shirou/gopsutil) - Cross-platform system monitoring library
- [ANSI Color](https://github.com/fatih/color) - Terminal color formatting
- [Go Prompt](https://github.com/c-bata/go-prompt) - Interactive prompt functionality
- [Clipboard](https://github.com/atotto/clipboard) - Cross-platform clipboard operations


### Community
- The Go community for their excellent documentation and support
- All open-source contributors who have helped improve Lumo

## 🔗 Additional Resources

- [Examples](docs/examples.md) - Comprehensive examples for all features
- [Project Wiki](https://github.com/agnath18K/lumo/wiki) - Detailed documentation and guides
- [Issue Tracker](https://github.com/agnath18K/lumo/issues) - Report bugs or request features
- [Discussions](https://github.com/agnath18K/lumo/discussions) - Community discussions and Q&A

---

<div align="center">
  <p>
    <a href="https://github.com/agnath18K/lumo/stargazers">
      <img src="https://img.shields.io/github/stars/agnath18K/lumo?style=social" alt="GitHub stars">
    </a>
    <a href="https://github.com/agnath18K/lumo/network/members">
      <img src="https://img.shields.io/github/forks/agnath18K/lumo?style=social" alt="GitHub forks">
    </a>
  </p>

  <p>Made with ❤️ by <a href="https://github.com/agnath18K">agnath18</a></p>

  <p>
    <a href="https://github.com/agnath18K">GitHub</a> •
    <a href="https://twitter.com/agnath18">Twitter</a> •
    <a href="https://linkedin.com/in/agnath18">LinkedIn</a>
  </p>
</div>
