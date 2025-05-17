package http_func

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	mcpTool "github.com/flipped-aurora/gin-vue-admin/server/mcp"
	"github.com/flipped-aurora/gin-vue-admin/server/mcp/http_func/service"
	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"
)

func init() {
	mcpTool.RegisterTool(&HttpApiTool{})
}

// HttpApiTool 是一个 HTTP API 工具的实现
type HttpApiTool struct{}

// New 返回工具注册信息
func (t *HttpApiTool) New() mcp.Tool {
	return mcp.NewTool("httpApi",
		mcp.WithDescription("通过HTTP API调用外部服务"),
		mcp.WithObject("parameters",
			mcp.Required(),
			mcp.Description("HTTP请求参数，包括URL、方法、请求体等"),
		),
	)
}

// Handle 处理工具调用请求
func (t *HttpApiTool) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 解析参数
	parameters, ok := request.Params.Arguments["parameters"].(map[string]interface{})
	if !ok {
		return nil, errors.New("参数错误：parameters 必须是对象类型")
	}

	// 提取必要参数
	apiUrl, ok := parameters["url"].(string)
	if !ok || apiUrl == "" {
		return nil, errors.New("参数错误：url 必须是非空字符串")
	}

	// 提取可选参数
	requestMethod := http.MethodGet // 默认方法
	if method, ok := parameters["method"].(string); ok && method != "" {
		requestMethod = strings.ToUpper(method)
	}

	// 验证HTTP方法
	switch requestMethod {
	case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete:
		// 方法有效，继续
	default:
		return nil, fmt.Errorf("不支持的HTTP方法：'%s'。支持的方法包括GET、POST、PUT、DELETE", requestMethod)
	}

	// 初始化HTTP客户端和请求数据
	httpClient := &http.Client{Timeout: 30 * time.Second}
	reqHeaders := make(http.Header)
	reqBodyMap := make(map[string]interface{})
	urlValues := url.Values{}

	// 处理请求头
	if headers, ok := parameters["headers"].(map[string]interface{}); ok {
		for key, value := range headers {
			if strValue, ok := value.(string); ok {
				reqHeaders.Add(key, strValue)
			} else if strArr, ok := value.([]string); ok {
				for _, s := range strArr {
					reqHeaders.Add(key, s)
				}
			} else if valueArr, ok := value.([]interface{}); ok {
				for _, v := range valueArr {
					if s, ok := v.(string); ok {
						reqHeaders.Add(key, s)
					}
				}
			}
		}
	}

	// 处理请求体和URL参数
	if body, ok := parameters["body"].(map[string]interface{}); ok && (requestMethod == http.MethodPost || requestMethod == http.MethodPut) {
		reqBodyMap = body
	}

	if queryParams, ok := parameters["query"].(map[string]interface{}); ok {
		for key, value := range queryParams {
			if strValue, ok := value.(string); ok {
				urlValues.Add(key, strValue)
			} else if valueArr, ok := value.([]interface{}); ok {
				for _, v := range valueArr {
					if s, ok := v.(string); ok {
						urlValues.Add(key, s)
					}
				}
			}
		}
	}

	// 添加URL参数
	if len(urlValues) > 0 {
		if strings.Contains(apiUrl, "?") {
			apiUrl += "&" + urlValues.Encode()
		} else {
			apiUrl += "?" + urlValues.Encode()
		}
	}

	// 准备请求体
	var reqBodyReader io.Reader
	if len(reqBodyMap) > 0 {
		reqBodyBytes, err := json.Marshal(reqBodyMap)
		if err != nil {
			return nil, fmt.Errorf("序列化请求体失败：%v", err)
		}
		reqBodyReader = strings.NewReader(string(reqBodyBytes))
		// 设置Content-Type为JSON，如果没有显式设置
		if _, found := reqHeaders["Content-Type"]; !found {
			reqHeaders.Set("Content-Type", "application/json")
		}
	}

	// 创建请求
	httpReq, err := http.NewRequestWithContext(ctx, requestMethod, apiUrl, reqBodyReader)
	if err != nil {
		global.GVA_LOG.Error("创建HTTP请求失败", zap.Error(err), zap.String("url", apiUrl))
		return nil, fmt.Errorf("创建HTTP请求失败：%v", err)
	}
	httpReq.Header = reqHeaders

	// 发送请求
	global.GVA_LOG.Info("发送HTTP请求",
		zap.String("method", requestMethod),
		zap.String("url", apiUrl))

	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		global.GVA_LOG.Error("HTTP请求失败", zap.Error(err), zap.String("url", apiUrl))
		return nil, fmt.Errorf("HTTP请求失败：%v", err)
	}
	defer httpResp.Body.Close()

	// 读取响应
	respBodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		global.GVA_LOG.Error("读取HTTP响应体失败", zap.Error(err), zap.String("url", apiUrl))
		return nil, fmt.Errorf("读取响应体失败：%v", err)
	}

	// 解析响应数据
	var respData interface{}
	contentType := httpResp.Header.Get("Content-Type")
	if strings.Contains(strings.ToLower(contentType), "application/json") {
		var jsonData interface{}
		if err := json.Unmarshal(respBodyBytes, &jsonData); err == nil {
			respData = jsonData
		} else {
			global.GVA_LOG.Warn("解析JSON响应体失败，以字符串形式返回",
				zap.Error(err),
				zap.String("contentType", contentType))
			respData = string(respBodyBytes)
		}
	} else {
		respData = string(respBodyBytes) // 非JSON响应，直接作为字符串
	}

	// 准备结果
	result := map[string]interface{}{
		"status_code": httpResp.StatusCode,
		"headers":     httpResp.Header,
		"body":        respData,
	}

	// 序列化为JSON
	jsonResultBytes, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("序列化结果失败：%v", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: string(jsonResultBytes),
			},
		},
	}, nil
}

// 创建代理工具，用于从数据库加载工具配置
func CreateHttpApiProxyTool() *HttpApiProxyTool {
	proxy := &HttpApiProxyTool{}
	mcpTool.RegisterTool(proxy)
	return proxy
}

// HttpApiProxyTool 是数据库配置的HTTP API工具的代理
type HttpApiProxyTool struct{}

// New 返回工具注册信息
func (t *HttpApiProxyTool) New() mcp.Tool {
	return mcp.NewTool("httpApiProxy",
		mcp.WithDescription("HTTP API工具代理，从数据库加载配置"),
		mcp.WithString("toolName",
			mcp.Required(),
			mcp.Description("要调用的工具名称"),
		),
		mcp.WithObject("parameters",
			mcp.Description("工具参数"),
		),
	)
}

// Handle 处理工具调用请求
func (t *HttpApiProxyTool) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 调用服务层的处理入口函数，完成对HTTP API工具的调用
	return service.HandleToolEntry(ctx, request)
}
