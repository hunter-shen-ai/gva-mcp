package router

import (
	"github.com/gin-gonic/gin"
)

var McpClientToolsRouter = new(mcpClientToolsRouter)

type mcpClientToolsRouter struct{}

func (r *mcpClientToolsRouter) Init(public *gin.RouterGroup, private *gin.RouterGroup) {

	{
		group := private.Group("mcpClientTools")
		group.GET("getMcpClientToolsList", apiMcpClientTools.ListMcpTools) // 获取mcpClientTools列表
		group.POST("executeTool", apiMcpClientTools.ExecuteMcpTool)        // 执行mcpClientTool
	}
}
