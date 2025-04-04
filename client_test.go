package main

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

// TestAptosMcpClient tests the interaction between Aptos-RAG-TOOL client and server
func TestAptosMcpClient(t *testing.T) {

	// Get server URL, can be overridden by environment variable
	serverURL := os.Getenv("MCP_SERVER_URL")
	if serverURL == "" {
		serverURL = "http://localhost:8282"
	}

	t.Logf("Connecting to MCP Server: %s", serverURL)

	// Create client
	c, err := client.NewSSEMCPClient(serverURL + "/sse")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer c.Close()

	t.Run("Can initialize and make requests", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Start the client
		if err := c.Start(ctx); err != nil {
			t.Fatalf("Failed to start client: %v", err)
		}

		// Initialize
		initRequest := mcp.InitializeRequest{}
		initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
		initRequest.Params.ClientInfo = mcp.Implementation{
			Name:    "test-client",
			Version: "1.0.0",
		}

		result, err := c.Initialize(ctx, initRequest)
		if err != nil {
			t.Fatalf("Failed to initialize: %v", err)
		}

		if result.ServerInfo.Name != "Aptos-MCP" {
			t.Errorf(
				"Expected server name 'Aptos-RAG-MCP', got '%s'",
				result.ServerInfo.Name,
			)
		}

		// Test Ping
		if err = c.Ping(ctx); err != nil {
			t.Errorf("Ping failed: %v", err)
		}

		// Test ListTools
		toolsRequest := mcp.ListToolsRequest{}
		toolListResult, err := c.ListTools(ctx, toolsRequest)
		if err != nil {
			t.Errorf("ListTools failed: %v", err)
		}
		t.Logf("Tool list: %v", toolListResult)

		request := mcp.CallToolRequest{}
		request.Params.Name = "Aptos-TOOL"
		request.Params.Arguments = map[string]any{
			"message": "I need to create an automated market maker (AMM) contract similar to Uniswap v2 for the Aptos blockchain. Please provide the main module structure and implementation approach, including factory contract, liquidity pool, swap function and fee mechanism.",
		}

		callToolResult, err := c.CallTool(ctx, request)
		if err != nil {
			t.Fatalf("CallTool failed: %v", err)
		}

		t.Logf("CallTool result: %v", callToolResult.Content)
	})

}
