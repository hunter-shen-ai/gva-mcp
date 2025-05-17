package request

// McpToolCallPayload defines the structure for calling an MCP tool.
type McpToolCallPayload struct {
	ToolName string                 `json:"tool_name" binding:"required"` // The name of the tool to call (e.g., "get_weather")
	Params   map[string]interface{} `json:"params"`                       // Parameters for the tool
}
