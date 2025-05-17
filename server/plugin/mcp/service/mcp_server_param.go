package service

import (
	"context"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/mcp/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/mcp/model/request"
)

var McpServerParam = new(mcpServerParam)

type mcpServerParam struct{}

// CreateMcpServerParam 创建mcpServerParam表记录
// Author [yourname](https://github.com/yourname)
func (s *mcpServerParam) CreateMcpServerParam(ctx context.Context, mcpServerParam *model.McpServerParam) (err error) {
	err = global.GVA_DB.Create(mcpServerParam).Error
	return err
}

// DeleteMcpServerParam 删除mcpServerParam表记录
// Author [yourname](https://github.com/yourname)
func (s *mcpServerParam) DeleteMcpServerParam(ctx context.Context, id string) (err error) {
	err = global.GVA_DB.Debug().Delete(&model.McpServerParam{}, "id = ?", id).Error
	return err
}

// DeleteMcpServerParamByIds 批量删除mcpServerParam表记录
// Author [yourname](https://github.com/yourname)
func (s *mcpServerParam) DeleteMcpServerParamByIds(ctx context.Context, ids []string) (err error) {
	err = global.GVA_DB.Delete(&[]model.McpServerParam{}, "id in ?", ids).Error
	return err
}

// UpdateMcpServerParam 更新mcpServerParam表记录
// Author [yourname](https://github.com/yourname)
func (s *mcpServerParam) UpdateMcpServerParam(ctx context.Context, mcpServerParam model.McpServerParam) (err error) {
	err = global.GVA_DB.Model(&model.McpServerParam{}).Where("id = ?", mcpServerParam.ID).Updates(&mcpServerParam).Error
	return err
}

// GetMcpServerParam 根据id获取mcpServerParam表记录
// Author [yourname](https://github.com/yourname)
func (s *mcpServerParam) GetMcpServerParam(ctx context.Context, id string) (mcpServerParam model.McpServerParam, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&mcpServerParam).Error
	return
}

// GetMcpServerParamInfoList 分页获取mcpServerParam表记录
// Author [yourname](https://github.com/yourname)
func (s *mcpServerParam) GetMcpServerParamInfoList(ctx context.Context, info request.McpServerParamSearch) (list []model.McpServerParam, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&model.McpServerParam{})
	var mcpServerParams []model.McpServerParam
	// 如果有条件搜索 下方会自动创建搜索语句

	// 根据 ToolID 筛选参数，这是非常常见的需求
	// 假设 info.ToolID 是 *uint 类型，在 request.McpServerParamSearch 中定义
	if info.ToolID != nil && *info.ToolID > 0 {
		db = db.Where("tool_id = ?", *info.ToolID)
	}

	if info.ParamName != nil && *info.ParamName != "" {
		db = db.Where("param_name LIKE ?", "%"+*info.ParamName+"%")
	}
	if info.ParamDescription != nil && *info.ParamDescription != "" {
		db = db.Where("param_description LIKE ?", "%"+*info.ParamDescription+"%")
	}
	if info.ParamRequired != nil {
		db = db.Where("param_required = ?", *info.ParamRequired)
	}
	// 可以根据需要为 ParamDataType 和 DefaultValue 添加搜索条件
	// if info.ParamDataType != nil && *info.ParamDataType != "" {
	// 	db = db.Where("param_data_type = ?", *info.ParamDataType)
	// }

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}
	err = db.Find(&mcpServerParams).Error
	return mcpServerParams, total, err
}

func (s *mcpServerParam) GetMcpServerParamPublic(ctx context.Context) {

}
