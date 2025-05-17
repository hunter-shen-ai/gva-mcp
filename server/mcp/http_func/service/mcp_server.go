package service

import (
	"context"
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/mcp/http_func/model"
	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"
)

// HandleToolEntry 处理HTTP API工具调用入口
func HandleToolEntry(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	toolName := request.Params.Name
	var mcpToolData model.McpServerTool

	// 1. 从数据库获取工具定义, 并预加载参数
	// GORM的Preload会自动填充McpServerTool结构体中的Parameters字段
	if err := global.GVA_DB.Preload("Parameters").Where("func_name = ?", toolName).First(&mcpToolData).Error; err != nil {
		global.GVA_LOG.Error("未找到工具配置", zap.String("func_name", toolName), zap.Error(err))
		return nil, fmt.Errorf("未找到工具 '%s': %v", toolName, err)
	}

	// 2. 确认工具类型为 HTTP API
	// if mcpToolData.FuncType == nil || *mcpToolData.FuncType != "http_api" {
	// 	global.GVA_LOG.Error("工具类型错误",
	// 		zap.String("func_name", toolName),
	// 		zap.Stringp("func_type", mcpToolData.FuncType))
	// 	return nil, fmt.Errorf("工具 '%s' 不是 HTTP API 类型", toolName)
	// }

	// 3. 调用 HTTP API 处理器
	return HandleHttpApiTool(ctx, request, &mcpToolData)
}
