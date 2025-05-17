import service from '@/utils/request'
// @Tags McpServerParam
// @Summary 创建mcpServerParam表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.McpServerParam true "创建mcpServerParam表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /mcpServerParam/createMcpServerParam [post]
export const createMcpServerParam = (data) => {
  return service({
    url: '/mcpServerParam/createMcpServerParam',
    method: 'post',
    data
  })
}

// @Tags McpServerParam
// @Summary 删除mcpServerParam表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.McpServerParam true "删除mcpServerParam表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /mcpServerParam/deleteMcpServerParam [delete]
export const deleteMcpServerParam = (params) => {
  return service({
    url: '/mcpServerParam/deleteMcpServerParam',
    method: 'delete',
    params
  })
}

// @Tags McpServerParam
// @Summary 批量删除mcpServerParam表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除mcpServerParam表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /mcpServerParam/deleteMcpServerParam [delete]
export const deleteMcpServerParamByIds = (params) => {
  return service({
    url: '/mcpServerParam/deleteMcpServerParamByIds',
    method: 'delete',
    params
  })
}

// @Tags McpServerParam
// @Summary 更新mcpServerParam表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.McpServerParam true "更新mcpServerParam表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /mcpServerParam/updateMcpServerParam [put]
export const updateMcpServerParam = (data) => {
  return service({
    url: '/mcpServerParam/updateMcpServerParam',
    method: 'put',
    data
  })
}

// @Tags McpServerParam
// @Summary 用id查询mcpServerParam表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.McpServerParam true "用id查询mcpServerParam表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /mcpServerParam/findMcpServerParam [get]
export const findMcpServerParam = (params) => {
  return service({
    url: '/mcpServerParam/findMcpServerParam',
    method: 'get',
    params
  })
}

// @Tags McpServerParam
// @Summary 分页获取mcpServerParam表列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取mcpServerParam表列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /mcpServerParam/getMcpServerParamList [get]
export const getMcpServerParamList = (params) => {
  return service({
    url: '/mcpServerParam/getMcpServerParamList',
    method: 'get',
    params
  })
}
// @Tags McpServerParam
// @Summary 不需要鉴权的mcpServerParam表接口
// @Accept application/json
// @Produce application/json
// @Param data query request.McpServerParamSearch true "分页获取mcpServerParam表列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /mcpServerParam/getMcpServerParamPublic [get]
export const getMcpServerParamPublic = () => {
  return service({
    url: '/mcpServerParam/getMcpServerParamPublic',
    method: 'get',
  })
}
