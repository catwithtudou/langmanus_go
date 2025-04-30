# LangManus Go

A Go implementation of [LangManus](https://github.com/Darwin-lfl/langmanus), built on the [Eino](https://github.com/cloudwego/eino) ecosystem.

## Overview

LangManus Go is a multi-agent system that combines language models with specialized tools for tasks like web search, crawling, and code execution. This implementation leverages the Eino ecosystem to provide a robust and scalable solution.

## Features

### Core Capabilities

- 🤖 **LLM Integration**
  - Support for multiple LLM models (basic, reasoning, vision)
  - Configurable API endpoints and keys
  - Multi-tier LLM system for different task complexities

### Tools and Integrations

- 🔍 **Search and Retrieval**
  - Web search via DuckDuckGo
  - HTML content extraction
  - Advanced content parsing

- 💻 **Code Execution**
  - Go REPL for Go code execution
  - Python REPL for Python code execution
  - Bash command execution

- 🌐 **Web Interaction**
  - Browser automation
  - Web page interaction
  - Content extraction

### Agent System

The system consists of the following specialized agents:

1. **Coordinator** - Entry point for task handling and routing
2. **Planner** - Task analysis and execution strategy creation
3. **Supervisor** - Execution oversight and agent management
4. **Researcher** - Information gathering and analysis
5. **Coder** - Code generation and execution
6. **Browser** - Web interaction and content extraction
7. **Reporter** - Results summarization and reporting

## Getting Started

### Prerequisites

- Go 1.24 or later
- Chrome browser (for browser automation)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/langmanus_go.git
cd langmanus_go
```

2. Install dependencies:
```bash
go mod download
```

3. Configure environment variables:
Create a `.env` file in the project root with the following variables:
```env
# Reasoning LLM Configuration
REASONING_MODEL=your_reasoning_model
REASONING_BASE_URL=your_reasoning_base_url
REASONING_API_KEY=your_reasoning_api_key

# Basic LLM Configuration
BASIC_MODEL=your_basic_model
BASIC_BASE_URL=your_basic_base_url
BASIC_API_KEY=your_basic_api_key

# Vision LLM Configuration
VL_MODEL=your_vision_model
VL_BASE_URL=your_vision_base_url
VL_API_KEY=your_vision_api_key

# Browser Configuration
CHROME_INSTANCE_PATH=/path/to/chrome
```

### Running the Server

```bash
go run cmd/server/main.go
```

## Project Structure

```
.
├── cmd/                 # Application entry points
│   ├── server/         # Server implementation
│   └── local/          # Local execution
├── config/             # Configuration management
├── internal/           # Internal packages
│   ├── agent/         # Agent implementations
│   ├── graph/         # Workflow graph
│   ├── llm/           # LLM client
│   ├── prompts/       # Agent prompts
│   ├── service/       # Service layer
│   └── tools/         # Tool implementations
└── log/               # Logging configuration
```

## Development

### Testing

```bash
go test ./...
```

### Building

```bash
go build -o langmanus cmd/server/main.go
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Original [LangManus](https://github.com/Darwin-lfl/langmanus) project
- [Eino](https://github.com/cloudwego/eino) ecosystem
- All contributors and open source projects that made this possible