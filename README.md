# ConstructAI

## Aptos Move Smart Contract Architecture Development Assistant

## Project Overview

ConstructAI is a smart contract architecture development assistant designed for Aptos Move ecosystem developers. It aims to address the challenges faced by Aptos Move developers by providing an intelligent architecture design workflow, helping developers systematically adopt ecosystem best practices, and offering smart assistance tools to accelerate development and improve code quality.

## Core Features

- **Intelligent Architecture Design Workflow**: Provides developers with a systematic design process to ensure best practices
- **Ecosystem Best Practices Integration**: Collects and analyzes excellent practices and patterns in the Aptos ecosystem
- **Smart Development Assistance**: Offers real-time coding suggestions, optimization solutions, and error detection
- **MCP Service Integration**: Connects to the backend Aptos Architecture Design Agent via MCP (Machine-aided Coding Protocol) service
- **Cross-platform Support**: Provides development support in any environment with MCP client compatibility

## Technical Architecture

ConstructAI implements the MCP server in Go language, connecting to the powerful Aptos Architecture Design Agent. This architecture enables:

1. Real-time intelligent suggestions in development tools
2. Seamless integration with existing workflows
3. Context-aware responses based on RAG (Retrieval-Augmented Generation) technology

## Quick Start

### Prerequisites

- Go 1.16 or higher
- Running Aptos Architecture Design Agent (port 7860)

### Installation

```bash
# Clone the repository
git clone https://github.com/dvbfcm850/ConstructAI.git
cd ConstructAI

# Install dependencies
go mod tidy
```

### Running

```bash
# Run in standard input/output mode
go run main.go

# Or run in SSE server mode
go run main.go -t sse
```

## Project Structure

```
constructai/
├── main.go            # MCP service main program
├── client_test.go     # Client test program
├── go.mod             # Go module definition
├── go.sum             # Dependency checksum
├── documents/         # Project documentation
├── logs/              # Log files directory
└── aptos-agent.json   # Aptos Agent configuration
```

## Usage Example

ConstructAI can be integrated into various development environments via MCP client:

```go
// Example: Interacting with MCP service
client := mcp.NewClient("http://localhost:8282")
response, err := client.CallTool("Aptos-TOOL", map[string]interface{}{
    "message": "How to design a storage structure for an Aptos token contract?",
})
```

## Advanced Features

- **Intelligent Architecture Validation**: Automatically detects potential security risks and optimization opportunities
- **Code Generation**: Generates template code based on design patterns
- **Interactive Learning**: Continuously optimizes recommendations through usage feedback

## License

This project is licensed under the [MIT License](LICENSE).
