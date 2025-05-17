package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/mcp/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/mcp/service/mcpClient"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// McpToolsApi provides API handlers for MCP tools.
type mcpClientTools struct{}

// ListMcpTools godoc
// @Tags McpTools
// @Summary 获取所有MCP工具列表
// @Description 获取所有MCP工具及其参数的列表，用于前端展示
// @Produce json
// @Success 200 {object} response.Response{data=[]model.McpToolView,msg=string} "获取成功"
// @Router /mcp/tools/list [get]
func (a *mcpClientTools) ListMcpTools(c *gin.Context) {
	mcpToolsListService := mcpClient.McpClientTools{}
	tools, err := mcpToolsListService.ListMcpTools()
	if err != nil {
		global.GVA_LOG.Error("获取MCP工具列表失败!", zap.Error(err))
		response.FailWithMessage("获取MCP工具列表失败", c)
		return
	}
	response.OkWithDetailed(tools, "获取成功", c)
}

// ExecuteMcpTool godoc
// @Tags McpTools
// @Summary 执行指定的MCP工具
// @Description 通过MCP客户端执行一个工具，并返回结果
// @Accept  json
// @Produce json
// @Param   toolPayload body request.McpToolCallPayload true "工具调用参数 (ToolName 和 Params)"
// @Success 200 {object} response.Response{data=mcp.CallToolResult,msg=string} "工具执行成功"
// @Router  /mcpClientTools/executeTool [post]
func (a *mcpClientTools) ExecuteMcpTool(c *gin.Context) {
	var payload request.McpToolCallPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		global.GVA_LOG.Error("参数绑定失败!", zap.Error(err))
		response.FailWithMessage("参数绑定失败: "+err.Error(), c)
		return
	}

	// Basic validation (can be expanded with GVA validators if needed)
	if payload.ToolName == "" {
		response.FailWithMessage("工具名称 (ToolName) 不能为空", c)
		return
	}

	toolCallSvc := mcpClient.ToolCallService{} // Use the correct service for tool execution
	result, err := toolCallSvc.ExecuteTool(c.Request.Context(), payload)
	if err != nil {
		global.GVA_LOG.Error("执行MCP工具失败!", zap.Error(err), zap.String("toolName", payload.ToolName))
		response.FailWithMessage("执行MCP工具失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(result, "工具执行成功", c)
}
