package mcpClient

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/mcp/model/request"
	mcpclient "github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"
)

// ToolCallService provides functionality to call MCP tools.
type ToolCallService struct{}

// ExecuteTool executes a specified MCP tool with the given parameters.
// It uses the mcp-go client to interact with an MCP server.
func (s *ToolCallService) ExecuteTool(ctx context.Context, payload request.McpToolCallPayload) (*mcp.CallToolResult, error) {
	// Get the MCP server URL from configuration. This is typically the SSE URL where the MCP server listens.
	mcpServerBaseURL := global.GVA_VP.GetString("mcp-config.mcp-server.sse-url") + "/sse"
	if mcpServerBaseURL == "" {
		global.GVA_LOG.Error("MCP server SSE URL is not configured (mcp-config.mcp-server.sse-url)")
		return nil, fmt.Errorf("MCP server URL not configured. Please set 'mcp-config.mcp-server.sse-url' in your config")
	}

	global.GVA_LOG.Info("使用MCP服务器URL",
		zap.String("url", mcpServerBaseURL),
		zap.String("config_key", "mcp-config.mcp-server.sse-url"))

	// 创建MCP客户端
	global.GVA_LOG.Info("开始创建MCP客户端")
	client, err := mcpclient.NewSSEMCPClient(mcpServerBaseURL)
	if err != nil {
		global.GVA_LOG.Error("创建MCP客户端失败", zap.Error(err))
		return nil, fmt.Errorf("failed to create MCP client: %w", err)
	}

	// 创建上下文，带超时
	timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// 启动客户端
	global.GVA_LOG.Info("正在启动MCP客户端")
	if err := client.Start(timeoutCtx); err != nil {
		global.GVA_LOG.Error("启动MCP客户端失败", zap.Error(err))
		return nil, fmt.Errorf("failed to start MCP client: %w", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			global.GVA_LOG.Warn("关闭MCP客户端失败", zap.Error(err))
		}
	}()

	// 初始化客户端（关键步骤，之前缺少）
	global.GVA_LOG.Info("开始初始化MCP客户端")
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "sre-aiops-platform",
		Version: "1.0.0",
	}

	initResult, err := client.Initialize(timeoutCtx, initRequest)
	if err != nil {
		global.GVA_LOG.Error("MCP客户端初始化失败", zap.Error(err))
		return nil, fmt.Errorf("failed to initialize MCP client: %w", err)
	}

	global.GVA_LOG.Info("MCP客户端初始化成功",
		zap.String("serverName", initResult.ServerInfo.Name),
		zap.String("serverVersion", initResult.ServerInfo.Version))

	// 创建工具调用请求
	req := mcp.CallToolRequest{}
	req.Params.Name = payload.ToolName
	req.Params.Arguments = payload.Params

	global.GVA_LOG.Info("准备调用MCP工具",
		zap.String("toolName", payload.ToolName),
		zap.Any("params", payload.Params))

	// 实现重试逻辑
	const maxRetries = 1
	var result *mcp.CallToolResult
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		// 如果不是第一次尝试，等待一段时间再重试
		if attempt > 0 {
			waitTime := time.Duration(500*attempt) * time.Millisecond
			global.GVA_LOG.Info("重试调用MCP工具",
				zap.Int("attempt", attempt+1),
				zap.Int("maxRetries", maxRetries),
				zap.Duration("waitTime", waitTime))
			time.Sleep(waitTime)
		}

		// 调用工具
		result, lastErr = client.CallTool(timeoutCtx, req)

		// 调用成功
		if lastErr == nil {
			global.GVA_LOG.Info("MCP工具调用成功",
				zap.String("toolName", payload.ToolName),
				zap.Int("contentLength", len(result.Content)))
			return result, nil
		}

		// 如果错误不是初始化相关的，直接返回错误
		if !strings.Contains(strings.ToLower(lastErr.Error()), "not init") &&
			!strings.Contains(strings.ToLower(lastErr.Error()), "client not") {
			global.GVA_LOG.Error("MCP工具调用失败",
				zap.String("toolName", payload.ToolName),
				zap.Error(lastErr))
			return nil, fmt.Errorf("error calling tool '%s': %w", payload.ToolName, lastErr)
		}

		// 如果是客户端未初始化错误，将在下一次循环中重试
		global.GVA_LOG.Warn("客户端未初始化，将重试",
			zap.Int("attempt", attempt+1),
			zap.Int("maxRetries", maxRetries),
			zap.Error(lastErr))
	}

	// 如果重试失败，返回最后一个错误
	global.GVA_LOG.Error("MCP工具调用重试次数已用尽",
		zap.String("toolName", payload.ToolName),
		zap.Int("maxRetries", maxRetries),
		zap.Error(lastErr))
	return nil, fmt.Errorf("error calling tool '%s' after %d attempts: %w",
		payload.ToolName, maxRetries, lastErr)
}

// Note: You might need to define or find the definition for mcp.ToolCallParamsData
// if it's not directly available in the mcp package or if it's named differently.
