package mcpTool

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/mcp/model"
	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"
)

// parseValueWithType attempts to parse a string value into a Go type suitable for JSON marshalling.
// dataTypeLC should be a lowercased string like "string", "int", "bool", "float".
func parseValueWithType(valueStr string, dataTypeLC string) (interface{}, error) {
	switch dataTypeLC {
	case "string":
		return valueStr, nil
	case "int", "integer":
		i, err := strconv.ParseInt(valueStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot parse '%s' as int: %v", valueStr, err)
		}
		return i, nil
	case "number", "float", "double":
		f, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot parse '%s' as float: %v", valueStr, err)
		}
		return f, nil
	case "bool", "boolean":
		b, err := strconv.ParseBool(valueStr)
		if err != nil {
			return nil, fmt.Errorf("cannot parse '%s' as bool: %v", valueStr, err)
		}
		return b, nil
	case "": // No data type specified, treat as string.
		return valueStr, nil
	default:
		// For unrecognized data types, treat as string.
		// This could be made stricter in the future if needed.
		// global.GVA_LOG.Warn("Unrecognized data type for default value parsing, treating as string", zap.String("dataType", dataTypeLC), zap.String("value", valueStr))
		return valueStr, nil
	}
}

// HandleHttpApiTool 处理类型为 HTTP API 的 MCP 工具调用
func HandleHttpApiTool(ctx context.Context, request mcp.CallToolRequest, toolInfo *model.McpServerTool) (*mcp.CallToolResult, error) {
	if toolInfo.ApiUrl == nil || *toolInfo.ApiUrl == "" {
		return nil, fmt.Errorf("API URL is not defined for tool %s", request.Params.Name)
	}
	apiUrl := *toolInfo.ApiUrl
	requestMethod := http.MethodGet // Default method
	if toolInfo.RequestMethod != nil && *toolInfo.RequestMethod != "" {
		requestMethod = strings.ToUpper(*toolInfo.RequestMethod)
	}

	// Validate the HTTP method against the supported ones
	switch requestMethod {
	case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete:
		// Method is valid, proceed
	default:
		global.GVA_LOG.Error("Unsupported HTTP method specified for tool",
			zap.String("toolName", request.Params.Name),
			zap.String("method", requestMethod))
		return nil, fmt.Errorf("unsupported HTTP method '%s' for tool '%s'. Supported methods are GET, POST, PUT, DELETE", requestMethod, request.Params.Name)
	}

	paramsFromDB := toolInfo.Parameters // 使用预加载的参数

	httpClient := &http.Client{Timeout: 30 * time.Second} // 增加超时时间
	reqHeaders := make(http.Header)
	reqBodyMap := make(map[string]interface{})
	urlValues := url.Values{}

	// 处理参数
	for _, paramDef := range paramsFromDB {
		paramName := safeString(paramDef.ParamName)
		if paramName == "" {
			global.GVA_LOG.Warn("Parameter definition found with empty name", zap.Uint("tool_id", toolInfo.ID))
			continue // 跳过无效的参数定义
		}

		var processedVal interface{}
		dataType := strings.ToLower(safeString(paramDef.ParamDataType)) // Get data type and lowercase once

		valFromArgs, argExists := request.Params.Arguments[paramName]

		if argExists {
			// Convert the argument from interface{} to string to be parsed
			valFromArgsStr := fmt.Sprintf("%v", valFromArgs)
			parsedValFromArgs, err := parseValueWithType(valFromArgsStr, dataType) // dataType is already defined above
			if err != nil {
				global.GVA_LOG.Error("Failed to parse value provided by caller for parameter as per its defined type",
					zap.String("paramName", paramName),
					zap.Any("providedValueRaw", valFromArgs),
					zap.String("providedValueStr", valFromArgsStr),
					zap.String("expectedDataType", dataType),
					zap.String("toolName", request.Params.Name),
					zap.Error(err))
				return nil, fmt.Errorf("parameter '%s' (expected type: %s) for tool '%s': error processing provided value '%v': %v", paramName, dataType, request.Params.Name, valFromArgs, err)
			}
			processedVal = parsedValFromArgs

			global.GVA_LOG.Debug("Parameter provided by caller, parsed according to data type",
				zap.String("paramName", paramName),
				zap.Any("originalValueFromCaller", valFromArgs),
				zap.Any("processedValue", processedVal),
				zap.String("dataType", dataType),
				zap.String("toolName", request.Params.Name))
		} else { // Argument not provided, handle default value
			if paramDef.DefaultValue != nil && *paramDef.DefaultValue != "" {
				defaultValueStr := *paramDef.DefaultValue
				parsedDefault, err := parseValueWithType(defaultValueStr, dataType)
				if err != nil {
					// If parsing the default value fails, it's a configuration error for the tool,
					// as the defined type and default value are inconsistent.
					global.GVA_LOG.Error("Failed to parse default value for parameter",
						zap.String("paramName", paramName),
						zap.String("defaultValue", defaultValueStr),
						zap.String("dataType", dataType),
						zap.String("toolName", request.Params.Name),
						zap.Error(err))
					return nil, fmt.Errorf("parameter '%s' (type: %s) for tool '%s': error processing default value '%s': %v", paramName, dataType, request.Params.Name, defaultValueStr, err)
				}
				processedVal = parsedDefault
				logMsg := "Parameter not provided by caller, using successfully parsed default value."
				if paramDef.ParamRequired != nil && *paramDef.ParamRequired {
					logMsg = "Required parameter not provided by caller, using successfully parsed default value."
				}
				global.GVA_LOG.Info(logMsg,
					zap.String("paramName", paramName),
					zap.Any("defaultValueUsed", processedVal),
					zap.String("dataType", dataType),
					zap.String("toolName", request.Params.Name))
			} else { // No argument from caller, and no default value configured
				if paramDef.ParamRequired != nil && *paramDef.ParamRequired {
					global.GVA_LOG.Error("Required parameter missing and no default value configured",
						zap.String("paramName", paramName),
						zap.String("toolName", request.Params.Name))
					return nil, fmt.Errorf("required parameter '%s' is missing for tool '%s' (and no default value is configured)", paramName, request.Params.Name)
				} else {
					// Optional parameter, not provided, no default value. Skip this parameter.
					global.GVA_LOG.Debug("Optional parameter not provided and no default value, skipping.",
						zap.String("paramName", paramName),
						zap.String("toolName", request.Params.Name))
					continue
				}
			}
		}

		// 新增：处理 IsSecure 参数
		if paramDef.IsSecure != nil && *paramDef.IsSecure {
			global.GVA_LOG.Info("Parameter is marked as secure. Overriding any provided value with its configured default value.",
				zap.String("paramName", paramName),
				zap.String("toolName", request.Params.Name))

			if paramDef.DefaultValue == nil || *paramDef.DefaultValue == "" {
				global.GVA_LOG.Error("Secure parameter has no default value configured, which is mandatory.",
					zap.String("paramName", paramName),
					zap.String("toolName", request.Params.Name))
				return nil, fmt.Errorf("secure parameter '%s' for tool '%s' must have a default value, but none is configured", paramName, request.Params.Name)
			}

			defaultValueStr := *paramDef.DefaultValue
			// dataType is already defined from earlier in the loop
			parsedSecureDefault, err := parseValueWithType(defaultValueStr, dataType)
			if err != nil {
				global.GVA_LOG.Error("Failed to parse default value for secure parameter. This is a configuration error.",
					zap.String("paramName", paramName),
					zap.String("defaultValue", defaultValueStr),
					zap.String("dataType", dataType),
					zap.String("toolName", request.Params.Name),
					zap.Error(err))
				return nil, fmt.Errorf("secure parameter '%s' (type: %s) for tool '%s': error processing its mandatory default value '%s': %v", paramName, dataType, request.Params.Name, defaultValueStr, err)
			}
			processedVal = parsedSecureDefault // Override with the parsed secure default value

			global.GVA_LOG.Info("Secure parameter value has been overridden with its parsed default value.",
				zap.String("paramName", paramName),
				zap.Any("overriddenValue", processedVal), // Log the actual value being used
				zap.String("toolName", request.Params.Name))
		}
		// 结束：IsSecure 参数处理

		// 如果执行到这里，processedVal 包含了要使用的值（如果参数为 nil，则可能为 nil）。
		// 无论是来自调用者参数，还是来自解析后的默认值。如果参数是可选的且调用者未提供且没有默认值，我们之前已经 'continue' 了。

		strValForHeaderOrQuery := fmt.Sprintf("%v", processedVal) // 用于 header/query 的字符串转换。
		requestType := safeString(paramDef.RequestType)

		switch requestType {
		case "Header":
			reqHeaders.Set(paramName, strValForHeaderOrQuery)
		case "Body":
			reqBodyMap[paramName] = processedVal // Use the correctly typed value
		case "Param": // URL 查询参数
			urlValues.Add(paramName, strValForHeaderOrQuery)
		default: // Includes empty string or other unspecified values
			global.GVA_LOG.Warn("Unknown or unspecified parameter requestType, defaulting to URL query parameter",
				zap.String("paramName", paramName),
				zap.String("specifiedRequestType", requestType), // Log what was actually specified
				zap.String("toolName", request.Params.Name))
			urlValues.Add(paramName, strValForHeaderOrQuery) // Default to query parameter
		}
	}

	// 准备请求体
	var reqBodyReader io.Reader
	if len(reqBodyMap) > 0 {
		if requestMethod == http.MethodGet || requestMethod == http.MethodHead {
			global.GVA_LOG.Warn("Request body provided for GET/HEAD request, will be ignored by most servers.",
				zap.String("method", requestMethod),
				zap.String("toolName", request.Params.Name))
		} else {
			jsonBody, err := json.Marshal(reqBodyMap)
			if err != nil {
				global.GVA_LOG.Error("Failed to marshal request body to JSON", zap.Error(err), zap.String("toolName", request.Params.Name))
				return nil, fmt.Errorf("failed to marshal request body for tool '%s': %v", request.Params.Name, err)
			}
			reqBodyReader = bytes.NewReader(jsonBody)
			if reqHeaders.Get("Content-Type") == "" { // 仅当用户未在参数中指定时设置
				reqHeaders.Set("Content-Type", "application/json")
			}
		}
	}

	// 将查询参数附加到URL
	if len(urlValues) > 0 {
		if strings.Contains(apiUrl, "?") {
			apiUrl += "&" + urlValues.Encode()
		} else {
			apiUrl += "?" + urlValues.Encode()
		}
	}

	// 创建请求
	httpReq, err := http.NewRequestWithContext(ctx, requestMethod, apiUrl, reqBodyReader)
	if err != nil {
		global.GVA_LOG.Error("Failed to create HTTP request", zap.Error(err), zap.String("url", apiUrl), zap.String("toolName", request.Params.Name))
		return nil, fmt.Errorf("failed to create HTTP request for tool '%s': %v", request.Params.Name, err)
	}
	httpReq.Header = reqHeaders

	// --> 添加调试日志 <--
	bodyLog := "(no body)"
	if reqBodyReader != nil {
		// Attempt to log the body content if it was generated
		// We re-marshal reqBodyMap for logging purposes, as reqBodyReader might have been consumed or is complex to read directly without consuming.
		bodyBytesForLog, marshalErr := json.Marshal(reqBodyMap)
		if marshalErr == nil {
			bodyLog = string(bodyBytesForLog)
		} else {
			bodyLog = fmt.Sprintf("(failed to marshal body map for logging: %v)", marshalErr)
		}
	}
	global.GVA_LOG.Debug("Sending HTTP request for tool - Details",
		zap.String("toolName", request.Params.Name),
		zap.String("method", requestMethod),
		zap.String("url", apiUrl),          // apiUrl already includes query params
		zap.Any("headers", httpReq.Header), // http.Header is map[string][]string
		zap.String("body", bodyLog),
	)
	// <-- 结束调试日志 -->

	// 发送请求
	global.GVA_LOG.Info("Sending HTTP request for tool",
		zap.String("toolName", request.Params.Name),
		zap.String("method", requestMethod),
		zap.String("url", apiUrl))
	// zap.Any("headers", reqHeaders), // 考虑移除或脱敏headers日志
	// zap.Any("bodyMap", reqBodyMap)) // 考虑移除或脱敏body日志

	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		global.GVA_LOG.Error("HTTP request failed", zap.Error(err), zap.String("url", apiUrl), zap.String("toolName", request.Params.Name))
		return nil, fmt.Errorf("HTTP request failed for tool '%s': %v", request.Params.Name, err)
	}
	defer httpResp.Body.Close()

	respBodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		global.GVA_LOG.Error("Failed to read HTTP response body", zap.Error(err), zap.String("url", apiUrl), zap.String("toolName", request.Params.Name))
		return nil, fmt.Errorf("failed to read response body for tool '%s': %v", request.Params.Name, err)
	}

	// 恢复 respData 和相关的 JSON 解析逻辑
	var respData interface{}
	contentType := httpResp.Header.Get("Content-Type")
	if strings.Contains(strings.ToLower(contentType), "application/json") {
		var jsonData interface{}
		if err := json.Unmarshal(respBodyBytes, &jsonData); err == nil {
			respData = jsonData
		} else {
			global.GVA_LOG.Warn("Failed to parse JSON response body despite Content-Type, returning as string",
				zap.Error(err),
				zap.String("contentType", contentType),
				zap.ByteString("body", respBodyBytes[:min(200, len(respBodyBytes))])) // Log first 200 bytes or less
			respData = string(respBodyBytes) // 回退到字符串
		}
	} else {
		respData = string(respBodyBytes) // 非JSON响应，直接作为字符串
	}

	global.GVA_LOG.Info("Received HTTP response for tool",
		zap.String("toolName", request.Params.Name),
		zap.Int("statusCode", httpResp.StatusCode))
	// zap.Any("responsePreview", respData)) // 响应可能很大，谨慎记录

	// 恢复 result map
	result := map[string]interface{}{
		"status_code": httpResp.StatusCode,
		"headers":     httpResp.Header, // 注意：httpResp.Header 是 http.Header (map[string][]string) 类型
		"body":        respData,
	}

	// 恢复将 map 结果序列化为 JSON 字符串
	jsonResultBytes, err := json.Marshal(result)
	if err != nil {
		global.GVA_LOG.Error("Failed to marshal result to JSON for tool response",
			zap.String("toolName", request.Params.Name),
			zap.Error(err))
		return nil, fmt.Errorf("failed to marshal result for tool '%s': %v", request.Params.Name, err)
	}

	// 使用 mcp.NewToolResultText 构造返回结果，参考 mcp_get_weather_tool.go
	return mcp.NewToolResultText(string(jsonResultBytes)), nil
}

// safeString 处理空指针字符串
func safeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// min function to safely slice bytes for logging
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
