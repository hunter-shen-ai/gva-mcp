package mcpClient

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/mcp/model"
	"go.uber.org/zap"
)

// McpToolsListService provides a service for listing MCP tools.
type McpClientTools struct{}

// ListMcpTools retrieves all MCP tools with their parameters, formatted for frontend display.
func (s *McpClientTools) ListMcpTools() ([]model.McpToolView, error) {
	var tools []model.McpServerTool
	// Preload Parameters for each tool
	if err := global.GVA_DB.Preload("Parameters").Find(&tools).Error; err != nil {
		global.GVA_LOG.Error("Failed to load MCP tools from database", zap.Error(err))
		return nil, err
	}

	var toolViews []model.McpToolView
	for _, tool := range tools {
		var paramViews []model.ToolParameterView
		for _, param := range tool.Parameters {
			paramView := model.ToolParameterView{
				ParamName:        safeString(param.ParamName),
				ParamDescription: safeString(param.ParamDescription),
				ParamRequired:    param.ParamRequired != nil && *param.ParamRequired,
				ParamType:        "string", // Defaulting to string, adjust if type info is available
			}
			paramViews = append(paramViews, paramView)
		}

		// Retrieve the Enabled status from the database model.McpServerTool
		// Assuming McpServerTool has an 'Enabled' field of type *bool.
		enabled := false // Default to false if not specified or nil
		if tool.Enabled != nil {
			enabled = *tool.Enabled
		}

		toolView := model.McpToolView{
			ID:              tool.ID,
			FuncName:        safeString(tool.FuncName),
			FuncDescription: safeString(tool.FuncDescription),
			Enabled:         enabled,
			Parameters:      paramViews,
		}
		toolViews = append(toolViews, toolView)
	}

	return toolViews, nil
}

// safeString handles nil string pointers.
func safeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
