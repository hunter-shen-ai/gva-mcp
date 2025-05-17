package router

import (
	"github.com/gin-gonic/gin"
)

// InitHttpApiRouter 初始化HTTP API工具相关的路由
func InitHttpApiRouter(public *gin.RouterGroup, private *gin.RouterGroup) {
	HttpApiToolsRouter.Init(public, private)
}
