package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type McpQaSessionSearch struct {
	Content      *string    `json:"content" form:"content"`
	StartAskTime *time.Time `json:"startAskTime" form:"startAskTime"`
	EndAskTime   *time.Time `json:"endAskTime" form:"endAskTime"`
	request.PageInfo
}

type EventSSEData struct {
	Delta   string `json:"delta" form:"delta"`
	MsgType string `json:"msgType" form:"msgType"`
}
