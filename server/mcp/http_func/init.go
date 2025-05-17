package http_func

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/mcp/http_func/router"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 初始化HTTP API工具
func InitHttpApiTools() {
	global.GVA_LOG.Info("初始化HTTP API MCP工具")

	// 注册HTTP API工具
	// HttpApiTool在init函数中已经自动注册了

	// 注册HTTP API代理工具
	proxy := CreateHttpApiProxyTool()
	if proxy != nil {
		global.GVA_LOG.Info("成功注册HTTP API代理工具", zap.String("name", "httpApiProxy"))
	}
}

// InitHttpApiRouter 初始化HTTP API相关路由
func InitHttpApiRouter(public *gin.RouterGroup, private *gin.RouterGroup) {
	// 调用路由模块的初始化函数
	router.InitHttpApiRouter(public, private)
}
