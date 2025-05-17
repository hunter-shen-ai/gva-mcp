package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

var McpServerTool = new(mcpServerTool)

type mcpServerTool struct{}

// Init 初始化 mcpServerTool表 路由信息
func (r *mcpServerTool) Init(public *gin.RouterGroup, private *gin.RouterGroup) {
	{
		group := private.Group("mcpServerTool").Use(middleware.OperationRecord())
		group.POST("createMcpServerTool", apiMcpServerTool.CreateMcpServerTool)             // 新建mcpServerTool表
		group.DELETE("deleteMcpServerTool", apiMcpServerTool.DeleteMcpServerTool)           // 删除mcpServerTool表
		group.DELETE("deleteMcpServerToolByIds", apiMcpServerTool.DeleteMcpServerToolByIds) // 批量删除mcpServerTool表
		group.PUT("updateMcpServerTool", apiMcpServerTool.UpdateMcpServerTool)              // 更新mcpServerTool表
	}
	{
		group := private.Group("mcpServerTool")
		group.GET("findMcpServerTool", apiMcpServerTool.FindMcpServerTool)       // 根据ID获取mcpServerTool表
		group.GET("getMcpServerToolList", apiMcpServerTool.GetMcpServerToolList) // 获取mcpServerTool表列表
	}
	{
		group := public.Group("mcpServerTool")
		group.GET("getMcpServerToolPublic", apiMcpServerTool.GetMcpServerToolPublic) // mcpServerTool表开放接口
	}
}
