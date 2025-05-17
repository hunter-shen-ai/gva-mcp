package service

import (
	"context"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/mcp/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/mcp/model/request"
	"gorm.io/gorm"
)

var McpServerTool = new(mcpServerTool)

type mcpServerTool struct{}

// CreateMcpServerTool 创建mcpServerTool表记录
// Author [yourname](https://github.com/yourname)
func (s *mcpServerTool) CreateMcpServerTool(ctx context.Context, mcpServerTool *model.McpServerTool) (err error) {
	err = global.GVA_DB.Create(mcpServerTool).Error
	return err
}

// DeleteMcpServerTool 删除mcpServerTool表记录
// Author [yourname](https://github.com/yourname)
func (s *mcpServerTool) DeleteMcpServerTool(ctx context.Context, id string) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 首先尝试加载工具，以确保它存在并且可以获取其 ID (虽然这里 id 是 string，但 GVA_MODEL.ID 是 uint)
		var tool model.McpServerTool
		if err := tx.First(&tool, "id = ?", id).Error; err != nil {
			return err // 如果工具不存在，则返回错误
		}

		// 删除关联的 Parameters
		if err := tx.Where("tool_id = ?", tool.ID).Delete(&model.McpServerParam{}).Error; err != nil {
			return err
		}

		// 删除工具本身
		if err := tx.Delete(&tool).Error; err != nil { // 直接删除加载到的 tool 对象
			return err
		}
		return nil
	})
	return err
}

// DeleteMcpServerToolByIds 批量删除mcpServerTool表记录
// Author [yourname](https://github.com/yourname)
func (s *mcpServerTool) DeleteMcpServerToolByIds(ctx context.Context, ids []string) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 注意：ids 是 []string。如果 McpServerTool 的 ID 在数据库中是数字类型，
		// GORM 在处理 `id in ?` 时通常能够处理字符串到数字的转换（取决于数据库驱动和小版本）。
		// 但更安全的方式是先查询出这些工具的实际数字 ID。
		// 为了简化，这里假设 GORM 的 `id in ?` 对字符串 ID 列表有效，或者 ids 列表中的字符串确实是数字。

		// 1. 删除与这些工具ID关联的所有参数
		// 首先需要获取这些 string IDs 对应的 uint IDs，因为 McpServerParam.ToolID 是 uint
		var tools []model.McpServerTool
		if err := tx.Where("id IN ?", ids).Find(&tools).Error; err != nil {
			return err
		}
		if len(tools) == 0 {
			return nil // 没有找到匹配的工具
		}

		var toolIDs []uint
		for _, tool := range tools {
			toolIDs = append(toolIDs, tool.ID)
		}

		if len(toolIDs) > 0 {
			if err := tx.Where("tool_id IN ?", toolIDs).Delete(&model.McpServerParam{}).Error; err != nil {
				return err
			}
		}

		// 2. 删除这些工具本身
		if err := tx.Delete(&model.McpServerTool{}, "id IN ?", ids).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// UpdateMcpServerTool 更新mcpServerTool表记录
// Author [yourname](https://github.com/yourname)
func (s *mcpServerTool) UpdateMcpServerTool(ctx context.Context, mcpServerTool model.McpServerTool) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 更新工具主表信息 (不包括 Parameters)
		if err := tx.Model(&mcpServerTool).Updates(model.McpServerTool{
			FuncName:        mcpServerTool.FuncName,
			FuncDescription: mcpServerTool.FuncDescription,
			FuncType:        mcpServerTool.FuncType,
			ApiUrl:          mcpServerTool.ApiUrl,
			RequestMethod:   mcpServerTool.RequestMethod,
			// GVA_MODEL fields like UpdatedAt will be handled by GORM hooks if present or automatically
		}).Error; err != nil {
			return err
		}

		// 更新关联的 Parameters
		// 对于传入的 Parameters，如果其 ID 为 0，GORM 会尝试创建它们；如果 ID 非 0，会尝试更新。
		// Replace 会删除 mcpServerTool 当前所有关联的 Parameters，然后用 mcpServerTool.Parameters 中的记录替换它们。
		if err := tx.Model(&mcpServerTool).Association("Parameters").Replace(mcpServerTool.Parameters); err != nil {
			return err
		}
		return nil
	})
	return err
}

// GetMcpServerTool 根据id获取mcpServerTool表记录
// Author [yourname](https://github.com/yourname)
func (s *mcpServerTool) GetMcpServerTool(ctx context.Context, id string) (mcpServerTool model.McpServerTool, err error) {
	err = global.GVA_DB.Preload("Parameters").Where("id = ?", id).First(&mcpServerTool).Error
	return
}

// GetMcpServerToolInfoList 分页获取mcpServerTool表记录
// Author [yourname](https://github.com/yourname)
func (s *mcpServerTool) GetMcpServerToolInfoList(ctx context.Context, info request.McpServerToolSearch) (list []model.McpServerTool, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&model.McpServerTool{})
	var mcpServerTools []model.McpServerTool
	// 如果有条件搜索 下方会自动创建搜索语句

	if info.FuncName != nil && *info.FuncName != "" {
		db = db.Where("func_name LIKE ?", "%"+*info.FuncName+"%")
	}
	if info.FuncDescription != nil && *info.FuncDescription != "" {
		db = db.Where("func_description LIKE ?", "%"+*info.FuncDescription+"%")
	}
	if info.FuncType != nil && *info.FuncType != "" {
		db = db.Where("func_type = ?", *info.FuncType)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}
	err = db.Find(&mcpServerTools).Error
	return mcpServerTools, total, err
}

func (s *mcpServerTool) GetMcpServerToolPublic(ctx context.Context) {

}
