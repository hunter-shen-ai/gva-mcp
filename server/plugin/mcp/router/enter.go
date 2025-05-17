package router

import "github.com/flipped-aurora/gin-vue-admin/server/plugin/mcp/api"

var (
	Router            = new(router)
	apiMcpServerTool  = api.Api.McpServerTool
	apiMcpServerParam = api.Api.McpServerParam
	apiMcpClientTools = api.Api.McpClientTools
)

type router struct {
	McpServerTool  mcpServerTool
	McpServerParam mcpServerParam
	McpClientTools mcpClientToolsRouter
}
