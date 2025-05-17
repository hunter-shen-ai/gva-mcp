package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/mcp/http_func/model"
	"go.uber.org/zap"
)

// ListHttpApiTools 获取数据库中所有HTTP API类型的工具
func ListHttpApiTools() ([]model.McpServerTool, error) {
	var tools []model.McpServerTool

	// 从数据库中查询所有HTTP API类型的工具
	err := global.GVA_DB.Preload("Parameters").Where("func_type = ?", "http_api").Find(&tools).Error
	if err != nil {
		global.GVA_LOG.Error("查询HTTP API工具失败", zap.Error(err))
		return nil, err
	}

	return tools, nil
}
