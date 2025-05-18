<div align="center">
  <h1>🐦 Lumo</h1>
  <p><b>Your AI-Powered Terminal Assistant</b></p>

  <a href="https://getlumo.dev"><b>getlumo.dev</b></a>

  <p>
    <img src="https://img.shields.io/badge/version-1.0.1-blue.svg" alt="Version 1.0.1">
    <img src="https://img.shields.io/badge/go-%3E%3D1.22-00ADD8.svg" alt="Go Version">
    <img src="https://img.shields.io/badge/license-MIT-green.svg" alt="License MIT">
  </p>
</div>

<div align="center">
  <a href="#-overview">Overview</a> •
  <a href="#-key-features">Key Features</a> •
  <a href="#-installation">Installation</a> •
  <a href="#-usage">Usage</a> •
  <a href="#-license">License</a>
</div>

## 📖 Overview

Lumo is an intelligent CLI assistant that interprets natural language to help you navigate the terminal with ease. It bridges the gap between natural language and terminal commands using advanced AI models.

**For comprehensive documentation, visit [getlumo.dev](https://getlumo.dev)**


## 🌟 Key Features

- **Natural Language Command Processing**: Translate plain English into terminal commands
- **Agent Mode**: Autonomous planning and execution of command sequences
- **Chat Mode**: Conversational assistance for terminal and general queries
- **Desktop Assistant**: Control your desktop environment with natural language
- **System Monitoring**: Track system health and performance
- **Pipe Support**: Analyze and explain command outputs
- **Web Interface**: Access Lumo through a browser-based interface
- **Secure Authentication**: JWT-based authentication for the REST API
- **Multiple AI Providers**: Support for Google Gemini, OpenAI, and Ollama

## 🚀 Installation

### Quick Install

```bash

# Using pre-built binary (Linux)
curl -L https://github.com/agnath18K/lumo_cli/releases/download/v1.0.1/lumo_1.0.1_linux_amd64.tar.gz -o lumo.tar.gz
tar -xzf lumo.tar.gz
sudo mv lumo /usr/local/bin/

# Using Debian package
curl -L https://github.com/agnath18K/lumo_cli/releases/download/v1.0.1/lumo_1.0.1_amd64.deb -o lumo.deb
sudo dpkg -i lumo.deb
```

**For detailed installation instructions and system requirements, visit [getlumo.dev/installation](https://getlumo.dev/installation)**

## 🔍 Usage

```bash
# Basic usage - ask in natural language
lumo "How do I find large files in Linux?"

# Agent mode - execute sequences of commands
lumo auto:create a backup of my documents folder

# Chat mode - conversational assistance
lumo chat

# Pipe support - analyze command output
ls -la | lumo

# System health check
lumo health

# Internet speed test
lumo speed

# Desktop assistant
lumo desktop:"close firefox window"
lumo desktop:"launch terminal"

# Web interface - start the server and access via browser
lumo server:start
# Then open http://localhost:7531 in your browser
# Default credentials: admin / lumo

# Enable authentication for the REST API
lumo config:server auth enable
```

**For complete usage documentation and examples, visit [getlumo.dev/documentation](https://getlumo.dev/documentation)**

**For information about the authentication system, see [Authentication Documentation](docs/authentication.md)**


## 🛠️ Development

**For development documentation, visit [getlumo.dev/documentation](https://getlumo.dev/documentation)**

Contributions to Lumo are welcome! Please fork the repository and submit a pull request.

## 📜 License

Lumo is released under the [MIT License](LICENSE).
---

<div align="center">
  <p>
    <a href="https://getlumo.dev">Website</a> •
    <a href="https://github.com/agnath18K/lumo_cli">GitHub</a> •
    <a href="https://github.com/agnath18K/lumo_cli/issues">Issues</a>
  </p>

  <p>Designed by <a href="https://github.com/agnath18K">agnath18</a></p>
</div>
