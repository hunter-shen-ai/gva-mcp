package api

import "github.com/flipped-aurora/gin-vue-admin/server/plugin/mcp/service"

var (
	Api                   = new(api)
	serviceMcpServerTool  = service.Service.McpServerTool
	serviceMcpServerParam = service.Service.McpServerParam
	serviceMcpClientTools = service.Service.McpClientTools
)

type api struct {
	McpServerTool  mcpServerTool
	McpServerParam mcpServerParam
	McpClientTools mcpClientTools
}
