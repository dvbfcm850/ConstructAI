package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type ToolName string

const (
	AptosTool ToolName = "Aptos-TOOL"
)

// logMCPRequest logs the MCP request and response to a timestamped file
func logMCPRequest(requestData any, responseData any, operation string) error {
	// Create logs directory if it doesn't exist
	logsDir := "logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %v", err)
	}

	// Create log file with timestamp
	timestamp := time.Now().Format("20060102-150405")
	logFilename := filepath.Join(logsDir, fmt.Sprintf("%s.log", timestamp))
	logFile, err := os.Create(logFilename)
	if err != nil {
		return fmt.Errorf("failed to create log file: %v", err)
	}
	defer logFile.Close()

	// Write operation type and timestamp
	fmt.Fprintf(logFile, "=== %s OPERATION [%s] ===\n\n", operation, time.Now().Format(time.RFC3339))

	// Write request data
	fmt.Fprintf(logFile, "REQUEST:\n")
	requestJSON, err := json.MarshalIndent(requestData, "", "  ")
	if err != nil {
		fmt.Fprintf(logFile, "Error marshaling request: %v\n", err)
	} else {
		fmt.Fprintf(logFile, "%s\n\n", requestJSON)
	}

	// Write response data
	fmt.Fprintf(logFile, "RESPONSE:\n")
	responseJSON, err := json.MarshalIndent(responseData, "", "  ")
	if err != nil {
		fmt.Fprintf(logFile, "Error marshaling response: %v\n", err)
	} else {
		fmt.Fprintf(logFile, "%s\n", responseJSON)
	}

	log.Printf("Log saved to %s", logFilename)
	return nil
}

func NewMCPServer() *server.MCPServer {
	log.Printf("Initializing Aptos-MCP server...")
	mcpServer := server.NewMCPServer(
		"Aptos-MCP",
		"0.1.0",
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithLogging(),
	)
	log.Printf("MCP server core components created")

	log.Printf("Registering Aptos-TOOL...")
	mcpServer.AddTool(mcp.NewTool(string(AptosTool),
		mcp.WithDescription("Retrieval-Augmented Generation tool for contextual responses"),
		mcp.WithString("message",
			mcp.Description("Input message to process"),
			mcp.Required(),
		),
	), handleAptosTool)
	log.Printf("Aptos-TOOL registration completed")

	mcpServer.AddNotificationHandler("notification", handleNotification)

	return mcpServer
}

func handleAptosTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	log.Printf("===> Aptos-TOOL has been called <==")
	log.Printf("Request parameters: %v", request.Params.Arguments)

	// Get message from parameters
	message, ok := request.Params.Arguments["message"].(string)
	if !ok {
		return nil, fmt.Errorf("message parameter is required and must be a string")
	}

	var agentId = "9a248a2d-89b7-402e-9e23-2112a3083c7f" // advanced
	// var agentId = "3bf2440b-9718-4f8d-b3cf-4d9ba0384219" // basic

	// Create API request
	apiURL := fmt.Sprintf("http://127.0.0.1:7860/api/v1/run/%s?stream=false", agentId)
	reqBody := fmt.Sprintf(`{
		"input_value": %q,
		"output_type": "chat",
		"input_type": "chat",
		"tweaks": {
			"ChatInput-a18M0": {},
			"ParseData-UWXBP": {},
			"Prompt-zalIe": {},
			"SplitText-9kwYE": {},
			"ChatOutput-oFtXw": {},
			"Directory-aGzT0": {},
			"NVIDIAEmbeddingsComponent-sQTue": {},
			"FAISS-SAdhf": {},
			"NVIDIAEmbeddingsComponent-GMjQV": {},
			"FAISS-wKFtX": {},
			"NVIDIAModelComponent-0s6HX": {},
			"OpenAIModel-gJpk6": {}
		}
	}`, message)

	log.Printf("Sending request to API: %s", apiURL)
	client := &http.Client{}
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(reqBody))
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	log.Printf("Sending HTTP request...")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to send request: %v", err)
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	log.Printf("Received HTTP response, status code: %d", resp.StatusCode)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned error: %s", body)
	}

	var responseData map[string]interface{}
	if err := json.Unmarshal(body, &responseData); err != nil {
		log.Printf("Failed to parse JSON response: %v", err)
		return nil, fmt.Errorf("Failed to parse response: %v", err)
	}

	var jsonOk bool
	var outputs []interface{}
	outputs, jsonOk = responseData["outputs"].([]interface{})
	if !jsonOk || len(outputs) == 0 {
		log.Printf("Invalid response format: missing outputs array")
		return nil, fmt.Errorf("Invalid response format: missing outputs array")
	}

	var output0 map[string]interface{}
	output0, jsonOk = outputs[0].(map[string]interface{})
	if !jsonOk {
		log.Printf("Invalid response format: outputs[0] is not an object")
		return nil, fmt.Errorf("Invalid response format: outputs[0] is not an object")
	}

	var output0Outputs []interface{}
	output0Outputs, jsonOk = output0["outputs"].([]interface{})
	if !jsonOk || len(output0Outputs) == 0 {
		log.Printf("Invalid response format: missing outputs[0].outputs array")
		return nil, fmt.Errorf("Invalid response format: missing outputs[0].outputs array")
	}

	var output0Output0 map[string]interface{}
	output0Output0, jsonOk = output0Outputs[0].(map[string]interface{})
	if !jsonOk {
		log.Printf("Invalid response format: outputs[0].outputs[0] is not an object")
		return nil, fmt.Errorf("Invalid response format: outputs[0].outputs[0] is not an object")
	}

	var results map[string]interface{}
	results, jsonOk = output0Output0["results"].(map[string]interface{})
	if !jsonOk {
		log.Printf("Invalid response format: missing outputs[0].outputs[0].results object")
		return nil, fmt.Errorf("Invalid response format: missing outputs[0].outputs[0].results object")
	}

	var messageObj map[string]interface{}
	messageObj, jsonOk = results["message"].(map[string]interface{})
	if !jsonOk {
		log.Printf("Invalid response format: missing outputs[0].outputs[0].results.message object")
		return nil, fmt.Errorf("Invalid response format: missing outputs[0].outputs[0].results.message object")
	}

	var responseText string
	responseText, jsonOk = messageObj["text"].(string)
	if !jsonOk {
		log.Printf("Invalid response format: missing outputs[0].outputs[0].results.message.text field")
		return nil, fmt.Errorf("Invalid response format: missing outputs[0].outputs[0].results.message.text field")
	}

	result := responseText
	// log.Printf("Processing completed, preparing response (response length: %d bytes)", len(body))
	log.Printf("Processing completed, preparing response: %s", result)

	// Prepare the response
	response := &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: result,
			},
		},
	}

	// Log the request and response
	logData := map[string]interface{}{
		"message": message,
		"url":     apiURL,
		"reqBody": reqBody,
	}

	logResponseData := map[string]interface{}{
		"status":  resp.StatusCode,
		"content": result,
	}

	if err := logMCPRequest(logData, logResponseData, "Aptos-TOOL"); err != nil {
		log.Printf("Warning: Failed to log request: %v", err)
	}

	return response, nil
}

func handleNotification(
	ctx context.Context,
	notification mcp.JSONRPCNotification,
) {
	log.Printf("Received notification: %s", notification.Method)
}

func main() {
	// Ensure logs directory exists
	logsDir := "logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		log.Fatalf("Failed to create logs directory: %v", err)
	}

	var transport string
	flag.StringVar(&transport, "t", "stdio", "Transport type (stdio or sse)")
	flag.StringVar(&transport, "transport", "stdio", "Transport type (stdio or sse)")
	flag.Parse()

	log.Printf("=== Aptos-MCP service starting... ===")
	mcpServer := NewMCPServer()
	log.Printf("=== Aptos-MCP service initialization completed ===")

	if transport == "sse" {
		sseServer := server.NewSSEServer(mcpServer, server.WithBaseURL("http://localhost:8282"))
		log.Printf("=== SSE server starting, listening on port :8282 ===")
		if err := sseServer.Start(":8283"); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	} else {
		log.Printf("=== Standard input/output (STDIO) server starting ===")
		if err := server.ServeStdio(mcpServer); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}

}
