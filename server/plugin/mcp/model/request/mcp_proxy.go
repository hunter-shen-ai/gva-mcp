package request

// McpProxyServerRequest 代理服务器请求
type McpProxyServerRequest struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	SseUrl  string `json:"sseUrl"`
	Enabled bool   `json:"enabled"`
}

// McpProxyServerListQuery 代理服务器列表查询
type McpProxyServerListQuery struct {
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"pageSize" form:"pageSize"`
	Name     string `json:"name" form:"name"`
	Enabled  *bool  `json:"enabled" form:"enabled"`
}

// McpProxyServerIDsRequest 代理服务器ID列表请求
type McpProxyServerIDsRequest struct {
	IDs []uint `json:"ids"`
}
