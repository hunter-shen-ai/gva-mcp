package config

type McpConfig struct {
	McpServer McpServer `mapstructure:"mcp-server" json:"mcp-server" yaml:"mcp-server"`
}

type McpServer struct {
	SseUrl string `mapstructure:"sse-url" json:"sse-url" yaml:"sse-url"`
}
