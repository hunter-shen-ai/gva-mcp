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

var McpServerParam = new(mcpServerParam)

type mcpServerParam struct{}

// CreateMcpServerParam 创建mcpServerParam表
// @Tags McpServerParam
// @Summary 创建mcpServerParam表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.McpServerParam true "创建mcpServerParam表"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /mcpServerParam/createMcpServerParam [post]
func (a *mcpServerParam) CreateMcpServerParam(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var info model.McpServerParam
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = serviceMcpServerParam.CreateMcpServerParam(ctx, &info)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	service.RestartSSE()
	response.OkWithMessage("创建成功", c)
}

// DeleteMcpServerParam 删除mcpServerParam表
// @Tags McpServerParam
// @Summary 删除mcpServerParam表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.McpServerParam true "删除mcpServerParam表"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /mcpServerParam/deleteMcpServerParam [delete]
func (a *mcpServerParam) DeleteMcpServerParam(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	id := c.Query("ID")
	err := serviceMcpServerParam.DeleteMcpServerParam(ctx, id)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	service.RestartSSE()
	response.OkWithMessage("删除成功", c)
}

// DeleteMcpServerParamByIds 批量删除mcpServerParam表
// @Tags McpServerParam
// @Summary 批量删除mcpServerParam表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /mcpServerParam/deleteMcpServerParamByIds [delete]
func (a *mcpServerParam) DeleteMcpServerParamByIds(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ids := c.QueryArray("ids[]")
	err := serviceMcpServerParam.DeleteMcpServerParamByIds(ctx, ids)
	if err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:"+err.Error(), c)
		return
	}
	service.RestartSSE()
	response.OkWithMessage("批量删除成功", c)
}

// UpdateMcpServerParam 更新mcpServerParam表
// @Tags McpServerParam
// @Summary 更新mcpServerParam表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.McpServerParam true "更新mcpServerParam表"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /mcpServerParam/updateMcpServerParam [put]
func (a *mcpServerParam) UpdateMcpServerParam(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var info model.McpServerParam
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = serviceMcpServerParam.UpdateMcpServerParam(ctx, info)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	service.RestartSSE()
	response.OkWithMessage("更新成功", c)
}

// FindMcpServerParam 用id查询mcpServerParam表
// @Tags McpServerParam
// @Summary 用id查询mcpServerParam表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param id query int true "用id查询mcpServerParam表"
// @Success 200 {object} response.Response{data=model.McpServerParam,msg=string} "查询成功"
// @Router /mcpServerParam/findMcpServerParam [get]
func (a *mcpServerParam) FindMcpServerParam(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	id := c.Query("id")
	remcpServerParam, err := serviceMcpServerParam.GetMcpServerParam(ctx, id)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(remcpServerParam, c)
}

// GetMcpServerParamList 分页获取mcpServerParam表列表
// @Tags McpServerParam
// @Summary 分页获取mcpServerParam表列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.McpServerParamSearch true "分页获取mcpServerParam表列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /mcpServerParam/getMcpServerParamList [get]
func (a *mcpServerParam) GetMcpServerParamList(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var pageInfo request.McpServerParamSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := serviceMcpServerParam.GetMcpServerParamInfoList(ctx, pageInfo)
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

// GetMcpServerParamPublic 不需要鉴权的mcpServerParam表接口
// @Tags McpServerParam
// @Summary 不需要鉴权的mcpServerParam表接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /mcpServerParam/getMcpServerParamPublic [get]
func (a *mcpServerParam) GetMcpServerParamPublic(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	// 此接口不需要鉴权 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
	serviceMcpServerParam.GetMcpServerParamPublic(ctx)
	response.OkWithDetailed(gin.H{"info": "不需要鉴权的mcpServerParam表接口信息"}, "获取成功", c)
}
