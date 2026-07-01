// Command doopl-mcp is a Model Context Protocol (MCP) server that exposes
// DeepL translation capabilities to Claude, Claude Desktop, and other
// MCP-compatible AI clients. It enables seamless translation within AI workflows
// without requiring direct API integration.
//
// # Protocol
//
// Implements MCP v1 over stdio with JSON-RPC 2.0 transport. The server receives
// MCP requests as JSON lines on stdin and writes responses on stdout. Compliant
// with the Model Context Protocol specification:
// https://modelcontextprotocol.io/
//
// # Startup
//
// Start the server by setting DEEPL_AUTH_KEY and running:
//
//	export DEEPL_AUTH_KEY="your-api-key"
//	./doopl-mcp serve
//
// The server is designed to be long-lived; it reads requests from stdin and
// writes responses to stdout in a streaming fashion.
//
// # MCP Tools Exposed
//
// - translate: Translate text or content into a target language. Supports all
//   doopl translation options (formality, glossaries, context, etc).
// - languages: List languages supported for a given resource type
//   (translate, document, glossary, write).
// - usage: Check API quota and current usage (characters, documents, limits).
//
// # Configuration
//
// Authentication:
// - DEEPL_AUTH_KEY: Required. Your DeepL API key.
// - DEEPL_SERVER_URL: Optional. Override API endpoint (for testing or custom deployments).
//
// # Usage with Claude Desktop
//
// Configure in ~/.claude/settings.json:
//
//	{
//	  "mcpServers": {
//	    "doopl": {
//	      "command": "/path/to/doopl-mcp",
//	      "args": ["serve"],
//	      "env": { "DEEPL_AUTH_KEY": "your-key" }
//	    }
//	  }
//	}
//
// Restart Claude Desktop to load the server. New translation tools will appear in the
// assistant's tool menu.
//
// # Error Handling
//
// MCP protocol errors are returned as JSON-RPC error objects with:
// - code: JSON-RPC error code (e.g., -32600 for invalid request)
// - message: Human-readable error description
// - data: Optional error details (e.g., API response body)
//
// # Concurrency
//
// The server processes requests serially (one at a time), making it safe for
// single-threaded clients. If a request hangs, the entire server is blocked.
//
// For details on the MCP protocol, see:
// https://modelcontextprotocol.io/introduction
package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	deepl "github.com/arashackdev/doopl/pkg/deepl"
)

// Version mirrors the doopl library version.
const Version = "0.0.1"

func main() {
	// Read DEEPL_AUTH_KEY from environment
	authKey := os.Getenv("DEEPL_AUTH_KEY")
	if authKey == "" {
		fmt.Fprintf(os.Stderr, "Error: DEEPL_AUTH_KEY environment variable not set\n")
		os.Exit(1)
	}

	// Create doopl client
	client, err := deepl.New(authKey, deepl.WithAppInfo("doopl-mcp", Version))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to create DeepL client: %v\n", err)
		os.Exit(1)
	}

	server := &MCPServer{
		client: client,
		mu:     &sync.Mutex{},
	}

	if err := server.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// MCPServer implements the Model Context Protocol.
type MCPServer struct {
	client *deepl.Client
	mu     *sync.Mutex
	// Track request IDs to match responses with requests
	requestID int64
}

// Run starts the MCP server on stdio.
func (s *MCPServer) Run() error {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024) // 1MB max line size

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		go s.handleRequest(line)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %w", err)
	}

	return nil
}

// JSONRPCRequest is a JSON-RPC 2.0 request.
type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// JSONRPCResponse is a JSON-RPC 2.0 response.
type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   *ErrorObj   `json:"error,omitempty"`
}

// ErrorObj is a JSON-RPC error object.
type ErrorObj struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Predefined error codes (JSON-RPC 2.0)
const (
	ParseError      = -32700
	InvalidRequest  = -32600
	MethodNotFound  = -32601
	InvalidParams   = -32602
	InternalError   = -32603
	ServerErrorBase = -32099
)

func (s *MCPServer) handleRequest(data []byte) {
	var req JSONRPCRequest
	if err := json.Unmarshal(data, &req); err != nil {
		s.writeError(nil, ParseError, "Parse error")
		return
	}

	if req.JSONRPC != "2.0" {
		s.writeError(req.ID, InvalidRequest, "Invalid Request: jsonrpc must be 2.0")
		return
	}

	if req.Method == "" {
		s.writeError(req.ID, InvalidRequest, "Invalid Request: method required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*1000000000) // 30 seconds
	defer cancel()

	switch req.Method {
	case "initialize":
		s.handleInitialize(req.ID)
	case "tools/list":
		s.handleToolsList(req.ID)
	case "tools/call":
		s.handleToolsCall(ctx, req.ID, req.Params)
	default:
		s.writeError(req.ID, MethodNotFound, fmt.Sprintf("Method not found: %s", req.Method))
	}
}

func (s *MCPServer) handleInitialize(id interface{}) {
	result := map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities": map[string]interface{}{
			"tools": map[string]interface{}{},
		},
		"serverInfo": map[string]interface{}{
			"name":    "doopl",
			"version": Version,
		},
	}
	s.writeResult(id, result)
}

func (s *MCPServer) handleToolsList(id interface{}) {
	tools := []map[string]interface{}{
		{
			"name":        "translate",
			"description": "Translate text to a target language using DeepL API",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"text": map[string]interface{}{
						"type":        "string",
						"description": "Text to translate",
					},
					"target_lang": map[string]interface{}{
						"type":        "string",
						"description": "Target language code (e.g., DE, FR, ES, EN)",
					},
					"source_lang": map[string]interface{}{
						"type":        "string",
						"description": "Source language code (optional, auto-detected if omitted)",
					},
					"formality": map[string]interface{}{
						"type":        "string",
						"enum":        []string{"default", "more", "less", "prefer_more", "prefer_less"},
						"description": "Formality level (optional)",
					},
					"glossary_id": map[string]interface{}{
						"type":        "string",
						"description": "Glossary ID to apply (optional)",
					},
				},
				"required": []string{"text", "target_lang"},
			},
		},
		{
			"name":        "languages",
			"description": "List supported languages for a given resource type",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"resource": map[string]interface{}{
						"type":        "string",
						"enum":        []string{"translate", "document", "glossary", "write"},
						"description": "Resource type (default: translate)",
					},
				},
			},
		},
		{
			"name":        "usage",
			"description": "Check API quota and usage for your DeepL account",
			"inputSchema": map[string]interface{}{
				"type": "object",
			},
		},
	}

	result := map[string]interface{}{
		"tools": tools,
	}
	s.writeResult(id, result)
}

type ToolCall struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

func (s *MCPServer) handleToolsCall(ctx context.Context, id interface{}, params json.RawMessage) {
	var call ToolCall
	if err := json.Unmarshal(params, &call); err != nil {
		s.writeError(id, InvalidParams, fmt.Sprintf("Invalid params: %v", err))
		return
	}

	switch call.Name {
	case "translate":
		s.callTranslate(ctx, id, call.Arguments)
	case "languages":
		s.callLanguages(ctx, id, call.Arguments)
	case "usage":
		s.callUsage(ctx, id, call.Arguments)
	default:
		s.writeError(id, InvalidRequest, fmt.Sprintf("Unknown tool: %s", call.Name))
	}
}

type TranslateArgs struct {
	Text       string  `json:"text"`
	TargetLang string  `json:"target_lang"`
	SourceLang *string `json:"source_lang,omitempty"`
	Formality  *string `json:"formality,omitempty"`
	GlossaryID *string `json:"glossary_id,omitempty"`
}

func (s *MCPServer) callTranslate(ctx context.Context, id interface{}, args json.RawMessage) {
	var params TranslateArgs
	if err := json.Unmarshal(args, &params); err != nil {
		s.writeError(id, InvalidParams, fmt.Sprintf("Invalid arguments: %v", err))
		return
	}

	if params.Text == "" {
		s.writeError(id, InvalidParams, "text is required")
		return
	}
	if params.TargetLang == "" {
		s.writeError(id, InvalidParams, "target_lang is required")
		return
	}

	opts := []deepl.TranslateTextOption{}

	if params.SourceLang != nil && *params.SourceLang != "" {
		opts = append(opts, deepl.WithSourceLang(*params.SourceLang))
	}

	if params.Formality != nil && *params.Formality != "" {
		formality := deepl.Formality(*params.Formality)
		opts = append(opts, deepl.WithFormality(formality))
	}

	if params.GlossaryID != nil && *params.GlossaryID != "" {
		opts = append(opts, deepl.WithGlossaryID(*params.GlossaryID))
	}

	results, err := s.client.TranslateText(ctx, []string{params.Text}, params.TargetLang, opts...)
	if err != nil {
		s.writeError(id, InternalError, fmt.Sprintf("Translation failed: %v", err))
		return
	}

	if len(results) == 0 {
		s.writeError(id, InternalError, "No translation results returned")
		return
	}

	result := map[string]interface{}{
		"text":          results[0].Text,
		"detected_lang": results[0].DetectedSourceLang,
	}

	s.writeResult(id, result)
}

type LanguagesArgs struct {
	Resource *string `json:"resource,omitempty"`
}

func (s *MCPServer) callLanguages(ctx context.Context, id interface{}, args json.RawMessage) {
	var params LanguagesArgs
	if err := json.Unmarshal(args, &params); err != nil {
		s.writeError(id, InvalidParams, fmt.Sprintf("Invalid arguments: %v", err))
		return
	}

	resource := "translate"
	if params.Resource != nil && *params.Resource != "" {
		resource = *params.Resource
	}

	langs, err := s.client.Languages(ctx, resource)
	if err != nil {
		s.writeError(id, InternalError, fmt.Sprintf("Failed to fetch languages: %v", err))
		return
	}

	languages := make([]map[string]interface{}, len(langs))
	for i, lang := range langs {
		languages[i] = map[string]interface{}{
			"code": string(lang.Code),
			"name": lang.Name,
		}
	}

	result := map[string]interface{}{
		"languages": languages,
	}

	s.writeResult(id, result)
}

func (s *MCPServer) callUsage(ctx context.Context, id interface{}, args json.RawMessage) {
	usage, err := s.client.Usage(ctx)
	if err != nil {
		s.writeError(id, InternalError, fmt.Sprintf("Failed to fetch usage: %v", err))
		return
	}

	result := map[string]interface{}{
		"character_count":     usage.CharacterCount,
		"character_limit":     usage.CharacterLimit,
		"document_count":      usage.DocumentCount,
		"document_limit":      usage.DocumentLimit,
		"team_document_count": usage.TeamDocumentCount,
		"team_document_limit": usage.TeamDocumentLimit,
	}

	s.writeResult(id, result)
}

func (s *MCPServer) writeResult(id interface{}, result interface{}) {
	response := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
	s.writeResponse(response)
}

func (s *MCPServer) writeError(id interface{}, code int, message string) {
	response := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &ErrorObj{
			Code:    code,
			Message: message,
		},
	}
	s.writeResponse(response)
}

func (s *MCPServer) writeResponse(response JSONRPCResponse) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshaling response: %v", err)
		return
	}

	fmt.Println(string(data))
}
