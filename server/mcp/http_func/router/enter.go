package router

import "github.com/flipped-aurora/gin-vue-admin/server/mcp/http_func/api"

var (
	Router       = new(router)
	apiHttpTools = api.Api.HttpApiTools
)

type router struct {
	HttpApiTools httpApiToolsRouter
}
