import service from '@/utils/request'

export const getMcpClientToolsList = (params) => {
  return service({
    url: '/mcpClientTools/getMcpClientToolsList',
    method: 'get',
    params
  })
}

 // 执行mcpClientTool
export const executeTool = (data) => {
    return service({
      url: '/mcpClientTools/executeTool',
      method: 'post',
      data
    })
  }