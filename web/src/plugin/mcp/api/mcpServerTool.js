import service from '@/utils/request'
// @Tags McpServerTool
// @Summary 创建mcpServerTool表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.McpServerTool true "创建mcpServerTool表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /mcpServerTool/createMcpServerTool [post]
export const createMcpServerTool = (data) => {
  return service({
    url: '/mcpServerTool/createMcpServerTool',
    method: 'post',
    data
  })
}

// @Tags McpServerTool
// @Summary 删除mcpServerTool表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.McpServerTool true "删除mcpServerTool表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /mcpServerTool/deleteMcpServerTool [delete]
export const deleteMcpServerTool = (params) => {
  return service({
    url: '/mcpServerTool/deleteMcpServerTool',
    method: 'delete',
    params
  })
}

// @Tags McpServerTool
// @Summary 批量删除mcpServerTool表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除mcpServerTool表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /mcpServerTool/deleteMcpServerTool [delete]
export const deleteMcpServerToolByIds = (params) => {
  return service({
    url: '/mcpServerTool/deleteMcpServerToolByIds',
    method: 'delete',
    params
  })
}

// @Tags McpServerTool
// @Summary 更新mcpServerTool表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.McpServerTool true "更新mcpServerTool表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /mcpServerTool/updateMcpServerTool [put]
export const updateMcpServerTool = (data) => {
  return service({
    url: '/mcpServerTool/updateMcpServerTool',
    method: 'put',
    data
  })
}

// @Tags McpServerTool
// @Summary 用id查询mcpServerTool表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.McpServerTool true "用id查询mcpServerTool表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /mcpServerTool/findMcpServerTool [get]
export const findMcpServerTool = (params) => {
  return service({
    url: '/mcpServerTool/findMcpServerTool',
    method: 'get',
    params: { id: params.ID }
  })
}

// @Tags McpServerTool
// @Summary 分页获取mcpServerTool表列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取mcpServerTool表列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /mcpServerTool/getMcpServerToolList [get]
export const getMcpServerToolList = (params) => {
  return service({
    url: '/mcpServerTool/getMcpServerToolList',
    method: 'get',
    params
  })
}
// @Tags McpServerTool
// @Summary 不需要鉴权的mcpServerTool表接口
// @Accept application/json
// @Produce application/json
// @Param data query request.McpServerToolSearch true "分页获取mcpServerTool表列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /mcpServerTool/getMcpServerToolPublic [get]
export const getMcpServerToolPublic = () => {
  return service({
    url: '/mcpServerTool/getMcpServerToolPublic',
    method: 'get',
  })
}
