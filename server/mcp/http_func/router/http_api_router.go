package router

import (
	"github.com/gin-gonic/gin"
)

var HttpApiToolsRouter = new(httpApiToolsRouter)

type httpApiToolsRouter struct{}

func (r *httpApiToolsRouter) Init(public *gin.RouterGroup, private *gin.RouterGroup) {
	// 配置HTTP API工具相关的路由
	{
		group := private.Group("mcpHttpFuncApi")
		group.GET("getHttpToolslist", apiHttpTools.ListHttpApiTools)   // 获取HTTP API工具列表
		group.POST("executeHttpTool", apiHttpTools.ExecuteHttpApiTool) // 执行HTTP API工具
	}
}
