package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/mcp/http_func/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/mcp/http_func/service"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"
)

// httpApiTools HTTP API工具的API处理
type httpApiTools struct{}

// ListHttpApiTools godoc
// @Tags HttpApiTools
// @Summary 获取所有HTTP API工具列表
// @Description 获取所有配置的HTTP API工具及其参数的列表
// @Produce json
// @Success 200 {object} response.Response{data=[]model.McpServerTool,msg=string} "获取成功"
// @Router /mcp/httpApi/list [get]
func (a *httpApiTools) ListHttpApiTools(c *gin.Context) {
	// 调用服务获取工具列表
	tools, err := service.ListHttpApiTools()
	if err != nil {
		global.GVA_LOG.Error("获取HTTP API工具列表失败!", zap.Error(err))
		response.FailWithMessage("获取工具列表失败", c)
		return
	}
	response.OkWithDetailed(tools, "获取成功", c)
}

// ExecuteHttpApiTool godoc
// @Tags HttpApiTools
// @Summary 执行指定的HTTP API工具
// @Description 通过MCP执行一个HTTP API工具，并返回结果
// @Accept json
// @Produce json
// @Param toolPayload body request.McpToolCallPayload true "工具调用参数 (ToolName 和 Params)"
// @Success 200 {object} response.Response{data=mcp.CallToolResult,msg=string} "工具执行成功"
// @Router /mcp/httpApi/execute [post]
func (a *httpApiTools) ExecuteHttpApiTool(c *gin.Context) {
	var payload request.McpToolCallPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		global.GVA_LOG.Error("参数绑定失败!", zap.Error(err))
		response.FailWithMessage("参数绑定失败: "+err.Error(), c)
		return
	}

	// 基本验证
	if payload.ToolName == "" {
		response.FailWithMessage("工具名称 (ToolName) 不能为空", c)
		return
	}

	// 创建请求
	request := mcp.CallToolRequest{}
	request.Params.Name = payload.ToolName
	request.Params.Arguments = payload.Params

	// 调用服务执行工具
	result, err := service.HandleToolEntry(c.Request.Context(), request)
	if err != nil {
		global.GVA_LOG.Error("执行HTTP API工具失败!", zap.Error(err), zap.String("toolName", payload.ToolName))
		response.FailWithMessage("执行工具失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(result, "工具执行成功", c)
}
