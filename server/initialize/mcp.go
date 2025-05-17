package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	mcpTool "github.com/flipped-aurora/gin-vue-admin/server/mcp"
	"github.com/flipped-aurora/gin-vue-admin/server/mcp/http_func"
	"github.com/mark3labs/mcp-go/server"
)

func McpRun() *server.SSEServer {
	config := global.GVA_CONFIG.MCP

	s := server.NewMCPServer(
		config.Name,
		config.Version,
	)

	global.GVA_MCP_SERVER = s

	// 初始化HTTP API工具
	http_func.InitHttpApiTools()

	// 注册所有工具
	mcpTool.RegisterAllTools(s)

	return server.NewSSEServer(s,
		server.WithSSEEndpoint(config.SSEPath),
		server.WithMessageEndpoint(config.MessagePath),
		server.WithBaseURL(config.UrlPrefix))
}
