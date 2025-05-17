package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

var McpServerParam = new(mcpServerParam)

type mcpServerParam struct{}

// Init 初始化 mcpServerParam表 路由信息
func (r *mcpServerParam) Init(public *gin.RouterGroup, private *gin.RouterGroup) {
	{
		group := private.Group("mcpServerParam").Use(middleware.OperationRecord())
		group.POST("createMcpServerParam", apiMcpServerParam.CreateMcpServerParam)             // 新建mcpServerParam表
		group.DELETE("deleteMcpServerParam", apiMcpServerParam.DeleteMcpServerParam)           // 删除mcpServerParam表
		group.DELETE("deleteMcpServerParamByIds", apiMcpServerParam.DeleteMcpServerParamByIds) // 批量删除mcpServerParam表
		group.PUT("updateMcpServerParam", apiMcpServerParam.UpdateMcpServerParam)              // 更新mcpServerParam表
	}
	{
		group := private.Group("mcpServerParam")
		group.GET("findMcpServerParam", apiMcpServerParam.FindMcpServerParam)       // 根据ID获取mcpServerParam表
		group.GET("getMcpServerParamList", apiMcpServerParam.GetMcpServerParamList) // 获取mcpServerParam表列表
	}
	{
		group := public.Group("mcpServerParam")
		group.GET("getMcpServerParamPublic", apiMcpServerParam.GetMcpServerParamPublic) // mcpServerParam表开放接口
	}
}
