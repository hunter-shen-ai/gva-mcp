package mcpTool

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

func init() {
	RegisterTool(&Example{})
}

type Example struct {
}

// mcp测试
func (t *Example) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// TODO: 实现工具逻辑
	// 参数示例:
	//
	// test_a := request.Params["test_a"]
	//
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				// TODO: 填充text内容
			},
		},
	}, nil
}

func (t *Example) New() mcp.Tool {
	return mcp.NewTool("example",
		mcp.WithDescription("mcp测试"),
		mcp.WithNumber("test_a", mcp.Required(),
			mcp.Description("测试"),
			mcp.DefaultNumber(123),
		),
	)
}
