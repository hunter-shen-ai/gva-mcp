package model

// ToolParameterView represents the structure of a tool's parameter for API response.
type ToolParameterView struct {
	ParamName        string `json:"param_name"`
	ParamDescription string `json:"param_description"`
	ParamRequired    bool   `json:"param_required"`
	ParamType        string `json:"param_type"` // Defaulting to "string" as per image
}

// McpToolView represents the structure of an MCP tool for API response.
type McpToolView struct {
	ID              uint                `json:"id"`
	FuncName        string              `json:"func_name"`
	FuncDescription string              `json:"func_description"`
	Enabled         bool                `json:"enabled"`
	Parameters      []ToolParameterView `json:"parameters"`
}
