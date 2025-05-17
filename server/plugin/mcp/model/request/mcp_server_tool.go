package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type McpServerToolSearch struct {
	FuncName        *string `json:"funcName" form:"funcName"`
	FuncDescription *string `json:"funcDescription" form:"funcDescription"`
	FuncType        *string `json:"funcType" form:"funcType"`
	request.PageInfo
}
