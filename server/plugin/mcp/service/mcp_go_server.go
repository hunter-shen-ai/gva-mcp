package service

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/mcp/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/mcp/service/mcpTool"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.uber.org/zap"
)

type MCPServer struct {
	server *server.MCPServer
}

var sseServer *server.SSEServer

func RestartSSE() {
	s := NewMCPServer()
	go func() {
		defer func() {
			if err := recover(); err != nil {
				global.GVA_LOG.Error("Server panic", zap.Any("error", err))
			}
		}()
		if sseServer != nil {
			if err := sseServer.Shutdown(context.Background()); err != nil {
				global.GVA_LOG.Error("Shutdown existing SSE server failed", zap.Error(err))
			}
			//这里给一个短暂的延迟来给系统释放空间
			time.Sleep(500 * time.Millisecond)
		}
		// 这里先获取mcp_sse_url参数
		// 这里直接从配置文件里获取
		sseUrl := global.GVA_VP.GetString("mcp-config.mcp-server.sse-url")

		// param, err := system.SysParamsServiceApp.GetSysParam("mcp_sse_url")
		// if err != nil {
		// 	global.GVA_LOG.Error("Failed to get sys param", zap.String("key", "mcp_sse_url"), zap.Error(err))
		// 	return
		// }
		// sseUrl := param.Value
		sseServer = s.ServeSSE(sseUrl)
		global.GVA_LOG.Info(fmt.Sprintf("SSE server listening on %s", sseUrl[strings.LastIndex(sseUrl, ":"):]))
		if err := sseServer.Start(sseUrl[strings.LastIndex(sseUrl, ":"):]); err != nil {
			global.GVA_LOG.Error(fmt.Sprintf("Server error: %v", err))
		}
	}()

}
func NewMCPServer() *MCPServer {
	mcpServer := server.NewMCPServer(
		"example-server",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithToolCapabilities(true),
	)

	loadMcpTools(mcpServer)
	return &MCPServer{server: mcpServer}
}

func loadMcpTools(s *server.MCPServer) {
	var tools []model.McpServerTool
	if err := global.GVA_DB.Find(&tools).Error; err != nil {
		global.GVA_LOG.Error("Load MCP tools failed", zap.Error(err))
		return
	}

	for _, tool := range tools {
		funcName := safeString(tool.FuncName)
		desc := safeString(tool.FuncDescription)
		global.GVA_LOG.Info("Register MCP Tool", zap.String("name", funcName), zap.String("desc", desc))

		options := []mcp.ToolOption{mcp.WithDescription(desc)}
		var params []model.McpServerParam

		if err := global.GVA_DB.Where("tool_id = ?", tool.ID).Find(&params).Error; err != nil {
			global.GVA_LOG.Error("Load tool params failed", zap.Uint("tool_id", tool.ID), zap.Error(err))
			continue
		}

		for _, param := range params {
			paramName := safeString(param.ParamName)
			paramDesc := safeString(param.ParamDescription)
			global.GVA_LOG.Info("Tool Param", zap.String("name", paramName), zap.String("desc", paramDesc))

			// 动态构建传递给 mcp.WithString 的参数
			// 首先添加必须的参数描述
			currentParamOpts := []interface{}{mcp.Description(paramDesc)}
			// 根据数据库中的 ParamRequired 字段决定是否添加 mcp.Required()
			if param.ParamRequired != nil && *param.ParamRequired {
				currentParamOpts = append(currentParamOpts, mcp.Required())
			}

			// 将参数选项展开传递
			// 注意：这里需要确认 mcp.WithString 是否接受 []interface{} 作为可变参数的替代
			// 或者是否有更合适的类型，例如 mcp.Option，或者直接动态调用
			// 为简化，我们先尝试这样。如果 mcp-go 有特定的选项类型，应使用该类型。
			// 构建一个 mcp.ToolOption
			var paramSpecificOption mcp.ToolOption
			// 假设 mcp.WithString 接受名称和一系列选项
			// 我们需要将 currentParamOpts 转换为 mcp.WithString 期望的类型。
			// 从原始代码推断，mcp.Description() 和 mcp.Required() 的返回值可以直接用。

			if param.ParamRequired != nil && *param.ParamRequired {
				paramSpecificOption = mcp.WithString(paramName, mcp.Description(paramDesc), mcp.Required())
			} else {
				paramSpecificOption = mcp.WithString(paramName, mcp.Description(paramDesc))
			}
			options = append(options, paramSpecificOption)
		}

		s.AddTool(mcp.NewTool(funcName, options...), handleToolEntry)
	}
}

// safeString 处理 nil 指针
func safeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
func (s *MCPServer) ServeSSE(addr string) *server.SSEServer {
	return server.NewSSEServer(s.server,
		server.WithBaseURL(addr),
		server.WithSSEContextFunc(authFromRequest),
	)
}

// authFromRequest extracts the auth token from the request headers.
func authFromRequest(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, ctxKeyAuth{}, r.Header.Get("Authorization"))
}

type ctxKeyAuth struct{}

// handleToolEntry 此方法是所有函数的入口，便于后期扩展
func handleToolEntry(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	toolName := request.Params.Name
	var mcpToolData model.McpServerTool

	// 1. 从数据库获取工具定义, 并预加载参数
	// GORM的Preload会自动填充McpServerTool结构体中的Parameters字段
	if err := global.GVA_DB.Preload("Parameters").Where("func_name = ?", toolName).First(&mcpToolData).Error; err != nil {
		global.GVA_LOG.Error("Failed to find tool by func_name", zap.String("func_name", toolName), zap.Error(err))
		return nil, fmt.Errorf("tool '%s' not found: %v", toolName, err)
	}

	// 2. 根据 FuncType 路由到 HTTP API 或 MCP代理 处理器
	if mcpToolData.FuncType != nil {
		// 处理HTTP API类型工具
		if *mcpToolData.FuncType == "http_api" {
			return mcpTool.HandleHttpApiTool(ctx, request, &mcpToolData)
		}

	}

	// 3. 处理其他类型的工具 (原有的 switch case)
	switch toolName {
	case "get_nickname":

		return mcpTool.GetNickname(request)

	default:
		// 如果 FuncType 不是 "http_api" 且不匹配任何内建 case
		global.GVA_LOG.Warn("Unhandled tool or tool type",
			zap.String("func_name", toolName),
			zap.Stringp("func_type", mcpToolData.FuncType))
		return nil, fmt.Errorf("tool '%s' (type: %s) cannot be handled", toolName, safeString(mcpToolData.FuncType))
	}
}
