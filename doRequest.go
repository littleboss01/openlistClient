package openlist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// doRequest 执行通用HTTP请求
func (c *OpenListAPI) doRequest(req *HTTPRequest, result interface{}) error {
	// 序列化请求体
	var bodyReader io.Reader
	if req.Body != nil {
		bodyBytes, err := json.Marshal(req.Body)
		if err != nil {
			return fmt.Errorf("序列化请求体失败: %w", err)
		}
		bodyReader = bytes.NewBuffer(bodyBytes)
	}

	// 创建HTTP请求
	httpReq, err := http.NewRequest(req.Method, req.URL, bodyReader)
	if err != nil {
		return fmt.Errorf("创建HTTP请求失败: %w", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	if c.getToken() != "" {
		httpReq.Header.Set("Authorization", c.getToken())
	}

	// 设置自定义请求头
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	// 发送请求
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("发送HTTP请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应体失败: %w", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP请求失败，状态码: %d, 响应体: %s", resp.StatusCode, string(respBody))
	}

	// 解析响应
	apiResp := &APIResponse{
		Data: result,
	}
	if err := json.Unmarshal(respBody, apiResp); err != nil {
		return fmt.Errorf("解析响应失败，响应体: %s, 原因: %w", string(respBody), err)
	}

	// 检查业务状态码
	if apiResp.Code != 200 {
		return fmt.Errorf("API调用失败，错误码: %d, 消息: %s", apiResp.Code, apiResp.Message)
	}

	return nil
}
