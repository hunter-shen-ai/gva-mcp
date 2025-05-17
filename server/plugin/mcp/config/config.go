package config

type Config struct {
	McpConfig McpConfig `mapstructure:"mcp-config" json:"mcp-config" yaml:"mcp-config"`
}
