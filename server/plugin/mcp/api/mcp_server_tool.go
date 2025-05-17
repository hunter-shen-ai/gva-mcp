package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/mcp/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/mcp/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/mcp/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var McpServerTool = new(mcpServerTool)

type mcpServerTool struct{}

// CreateMcpServerTool 创建mcpServerTool表
// @Tags McpServerTool
// @Summary 创建mcpServerTool表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.McpServerTool true "创建mcpServerTool表"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /mcpServerTool/createMcpServerTool [post]
func (a *mcpServerTool) CreateMcpServerTool(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var info model.McpServerTool
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = serviceMcpServerTool.CreateMcpServerTool(ctx, &info)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	service.RestartSSE()
	response.OkWithMessage("创建成功", c)
}

// DeleteMcpServerTool 删除mcpServerTool表
// @Tags McpServerTool
// @Summary 删除mcpServerTool表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.McpServerTool true "删除mcpServerTool表"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /mcpServerTool/deleteMcpServerTool [delete]
func (a *mcpServerTool) DeleteMcpServerTool(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	id := c.Query("ID")
	err := serviceMcpServerTool.DeleteMcpServerTool(ctx, id)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	service.RestartSSE()
	response.OkWithMessage("删除成功", c)
}

// DeleteMcpServerToolByIds 批量删除mcpServerTool表
// @Tags McpServerTool
// @Summary 批量删除mcpServerTool表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /mcpServerTool/deleteMcpServerToolByIds [delete]
func (a *mcpServerTool) DeleteMcpServerToolByIds(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ids := c.QueryArray("ids[]")
	err := serviceMcpServerTool.DeleteMcpServerToolByIds(ctx, ids)
	if err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:"+err.Error(), c)
		return
	}
	service.RestartSSE()
	response.OkWithMessage("批量删除成功", c)
}

// UpdateMcpServerTool 更新mcpServerTool表
// @Tags McpServerTool
// @Summary 更新mcpServerTool表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.McpServerTool true "更新mcpServerTool表"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /mcpServerTool/updateMcpServerTool [put]
func (a *mcpServerTool) UpdateMcpServerTool(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var info model.McpServerTool
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = serviceMcpServerTool.UpdateMcpServerTool(ctx, info)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	service.RestartSSE()
	response.OkWithMessage("更新成功", c)
}

// FindMcpServerTool 用id查询mcpServerTool表
// @Tags McpServerTool
// @Summary 用id查询mcpServerTool表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param id query int true "用id查询mcpServerTool表"
// @Success 200 {object} response.Response{data=model.McpServerTool,msg=string} "查询成功"
// @Router /mcpServerTool/findMcpServerTool [get]
func (a *mcpServerTool) FindMcpServerTool(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	id := c.Query("id")
	remcpServerTool, err := serviceMcpServerTool.GetMcpServerTool(ctx, id)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(remcpServerTool, c)
}

// GetMcpServerToolList 分页获取mcpServerTool表列表
// @Tags McpServerTool
// @Summary 分页获取mcpServerTool表列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.McpServerToolSearch true "分页获取mcpServerTool表列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /mcpServerTool/getMcpServerToolList [get]
func (a *mcpServerTool) GetMcpServerToolList(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var pageInfo request.McpServerToolSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := serviceMcpServerTool.GetMcpServerToolInfoList(ctx, pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// GetMcpServerToolPublic 不需要鉴权的mcpServerTool表接口
// @Tags McpServerTool
// @Summary 不需要鉴权的mcpServerTool表接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /mcpServerTool/getMcpServerToolPublic [get]
func (a *mcpServerTool) GetMcpServerToolPublic(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	// 此接口不需要鉴权 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
	serviceMcpServerTool.GetMcpServerToolPublic(ctx)
	response.OkWithDetailed(gin.H{"info": "不需要鉴权的mcpServerTool表接口信息"}, "获取成功", c)
}
