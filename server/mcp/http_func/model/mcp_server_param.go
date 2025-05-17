package model

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// McpServerParam mcpServerParam表 结构体
type McpServerParam struct {
	global.GVA_MODEL
	ToolID           uint    `json:"toolID" form:"toolID" gorm:"comment:关联的MCP工具ID;column:tool_id;index;"`                             // 关联的MCP工具ID
	ParamName        *string `json:"paramName" form:"paramName" gorm:"comment:参数名称;column:param_name;size:512;"`                       //参数名称
	ParamDescription *string `json:"paramDescription" form:"paramDescription" gorm:"comment:参数描述;column:param_description;size:2048;"` //参数描述
	ParamDataType    *string `json:"paramDataType" form:"paramDataType" gorm:"comment:参数数据类型;column:param_data_type;size:50;"`         //参数数据类型 (e.g., string, int, bool)
	ParamRequired    *bool   `json:"paramRequired" form:"paramRequired" gorm:"comment:是否必需：0 否，1 是;column:param_required;"`            //是否必需
	DefaultValue     *string `json:"defaultValue" form:"defaultValue" gorm:"comment:默认值;column:default_value;size:255;"`               //默认值
	RequestType      *string `json:"requestType" form:"requestType" gorm:"comment:请求是params还是body;column:request_type;size:10;"`       // 请求是params还是body
	IsSecure         *bool   `json:"isSecure" form:"isSecure" gorm:"comment:是否为安全字段;column:is_secure;"`                                // 是否安全
}

// TableName mcpServerParam表 McpServerParam自定义表名 mcp_server_param
func (McpServerParam) TableName() string {
	return "mcp_server_param"
}
