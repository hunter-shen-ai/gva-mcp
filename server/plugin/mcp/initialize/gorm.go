package initialize

import (
	"context"
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/mcp/model"
)

func Gorm(ctx context.Context) {
	err := global.GVA_DB.WithContext(ctx).AutoMigrate(
		model.McpServerTool{},
		model.McpServerParam{},
	)
	if err != nil {
		global.GVA_LOG.Error(fmt.Sprintf("注册表失败: %v", err))
		return
	}
	global.GVA_LOG.Info("register table success")
}
