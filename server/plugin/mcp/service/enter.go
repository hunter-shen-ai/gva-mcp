package service

import "github.com/flipped-aurora/gin-vue-admin/server/plugin/mcp/service/mcpClient"

var Service = new(service)

type service struct {
	McpServerTool  mcpServerTool
	McpServerParam mcpServerParam
	McpClientTools mcpClient.McpClientTools
}
