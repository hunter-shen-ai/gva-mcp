package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type McpServerParamSearch struct {
	// FuncName         *string `json:"funcName" form:"funcName"`
	ToolID           *uint   `json:"toolID" form:"toolID"`
	ParamName        *string `json:"paramName" form:"paramName"`
	ParamDescription *string `json:"paramDescription" form:"paramDescription"`
	ParamRequired    *bool   `json:"paramRequired" form:"paramRequired"`
	// Can add ParamDataType and DefaultValue if needed for search
	request.PageInfo
}
