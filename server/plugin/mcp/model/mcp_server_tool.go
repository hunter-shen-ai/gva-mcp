package model

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// McpServerTool mcpServerTool表 结构体
type McpServerTool struct {
	global.GVA_MODEL
	FuncName        *string          `json:"funcName" form:"funcName" gorm:"comment:函数名称/工具名称;column:func_name;size:128;"`                       //函数名称/工具名称
	FuncDescription *string          `json:"funcDescription" form:"funcDescription" gorm:"comment:函数描述/工具描述;column:func_description;size:2048;"` //函数描述/工具描述
	FuncType        *string          `json:"funcType" form:"funcType" gorm:"comment:函数类型/工具类型;column:func_type;size:50;"`                        //函数类型/工具类型 (e.g., MCP_AGENT, SYSTEM_BUILTIN, THIRD_PARTY)
	ApiUrl          *string          `json:"apiUrl" form:"apiUrl" gorm:"comment:API地址/来源;column:api_url;size:255;"`                              // API地址/来源
	RequestMethod   *string          `json:"requestMethod" form:"requestMethod" gorm:"comment:请求方法;column:request_method;size:10;"`              // 请求方法 (e.g., GET, POST)
	Parameters      []McpServerParam `json:"parameters" gorm:"foreignKey:ToolID;comment:工具参数列表"`                                                 // 工具参数列表, 通过 McpServerParam.ToolID 关联
	Enabled         *bool            `json:"enabled" form:"enabled" gorm:"column:enabled;default:true;comment:是否启用"`
}

// TableName McpServerTool 表 McpServerTool自定义表名 mcp_server_tool
func (McpServerTool) TableName() string {
	return "mcp_server_tool"
}
