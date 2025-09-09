package openlist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// OpenListAPI OpenList 服务客户端
type OpenListAPI struct {
	baseURL        string       // 服务基础URL（如 http://localhost:5244）
	username       string       // 管理员用户名
	password       string       // 管理员密码
	proxy          string       // 代理地址（如 http://127.0.0.1:8080）
	token          string       // 登录令牌
	httpClient     *http.Client // HTTP客户端（带代理配置）
	mu             sync.RWMutex // 并发安全锁（保护token、proxy状态）
	proxyTested    bool         // 代理是否已测试
	proxyAvailable bool         // 代理是否可用
}

// NewOpenListAPI 创建OpenListAPI客户端实例
func NewOpenListAPI(baseURL, username, password, proxy string) *OpenListAPI {
	// 处理baseURL末尾的斜杠（确保统一格式）
	baseURL = strings.TrimSuffix(baseURL, "/")

	client := &OpenListAPI{
		baseURL:  baseURL,
		username: username,
		password: password,
		proxy:    proxy,
		// 初始化HTTP客户端（超时时间默认30秒，后续可根据接口需求调整）
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	// 若配置了代理，初始化代理客户端
	if proxy != "" {
		client.initProxyClient()
	}

	return client
}

// initProxyClient 初始化带代理的HTTP客户端
func (c *OpenListAPI) initProxyClient() {
	// 解析代理URL
	proxyURL, err := url.Parse(c.proxy)
	if err != nil {
		// 直接返回，不打印日志
		return
	}

	// 设置代理客户端
	c.httpClient.Transport = &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
		// 基础TCP配置（复用连接、超时控制）
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second, // 拨号超时
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}
}

// TestProxy 测试代理是否可用
func (c *OpenListAPI) TestProxy() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 无代理配置时返回true（视为"无需代理即可用"）
	if c.proxy == "" {
		return true
	}

	// 已测试过代理，直接返回缓存结果
	if c.proxyTested {
		return c.proxyAvailable
	}

	// 解析代理地址（提取主机和端口）
	proxyURL, err := url.Parse(c.proxy)
	if err != nil {
		c.proxyTested = true
		c.proxyAvailable = false
		return false
	}

	// 提取代理的主机和端口（处理无端口的情况，默认http用80，https用443）
	host := proxyURL.Host
	if !strings.Contains(host, ":") {
		if proxyURL.Scheme == "https" {
			host += ":443"
		} else {
			host += ":80"
		}
	}

	// 尝试TCP连接代理（5秒超时）
	conn, err := net.DialTimeout("tcp", host, 5*time.Second)
	if err != nil {
		c.proxyTested = true
		c.proxyAvailable = false
		return false
	}
	defer conn.Close()

	// 代理测试成功
	c.proxyTested = true
	c.proxyAvailable = true
	return true
}

// ResetProxyStatus 重置代理状态（修改代理配置后调用，重新测试）
func (c *OpenListAPI) ResetProxyStatus() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.proxyTested = false
	c.proxyAvailable = false
}

// Login 登录OpenList服务，获取访问令牌
func (c *OpenListAPI) Login() (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 若已存在有效令牌，直接返回成功
	if c.token != "" {
		return true, nil
	}

	// 构造登录请求体
	loginReq := LoginRequest{
		Username: c.username,
		Password: c.password,
	}

	// 执行请求
	loginResp := &LoginResponse{}
	if err := c.doRequest(&HTTPRequest{
		Method: "POST",
		URL:    fmt.Sprintf("%s/api/auth/login", c.baseURL),
		Body:   loginReq,
	}, loginResp); err != nil {
		return false, fmt.Errorf("登录失败: %w", err)
	}

	// 登录成功，保存令牌
	c.token = loginResp.Data.Token
	return true, nil
}

// UploadFile 上传文件到OpenList服务
// filePaths: 本地文件路径
// remotePath: 远程存储目录（如 "/docs"）
// 返回值: 远程文件完整路径（如 "/docs/test.txt"），错误信息
func (c *OpenListAPI) UploadFile(filePath, remotePath string) (string, error) {
	// 先检查登录状态
	if ok, err := c.Login(); !ok {
		if err != nil {
			return "", fmt.Errorf("登录失败: %w", err)
		}
		return "", fmt.Errorf("登录失败，无法执行文件上传")
	}

	// 验证本地文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("本地文件不存在: %s", filePath)
	}

	// 打开本地文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("打开本地文件失败: %w", err)
	}
	defer file.Close()

	// 构造远程完整路径（处理重复斜杠，如 "/docs//test.txt" → "/docs/test.txt"）
	fileName := filepath.Base(filePath)
	fullRemotePath := strings.ReplaceAll(fmt.Sprintf("%s/%s", remotePath, fileName), "//", "/")
	// URL编码远程路径（保留斜杠，避免转义）
	encodedPath := url.QueryEscape(fullRemotePath)
	encodedPath = strings.ReplaceAll(encodedPath, "%2F", "/") // 恢复斜杠

	// 构造上传请求URL
	reqURL := fmt.Sprintf("%s/api/fs/form", c.baseURL)

	// 创建multipart/form-data请求体
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	// 添加文件字段（字段名"file"需与服务端一致）
	formFile, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return "", fmt.Errorf("创建表单文件失败: %w", err)
	}
	// 复制文件内容到表单
	if _, err := io.Copy(formFile, file); err != nil {
		return "", fmt.Errorf("复制文件到表单失败: %w", err)
	}
	// 关闭writer，确保边界符正确写入
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("关闭表单写入器失败: %w", err)
	}

	// 构造HTTP请求
	req, err := http.NewRequest("PUT", reqURL, body)
	if err != nil {
		return "", fmt.Errorf("创建上传请求失败: %w", err)
	}
	// 设置请求头（Authorization、Content-Type、file-path）
	req.Header.Set("Authorization", c.getToken())
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("file-path", encodedPath)
	// 延长上传超时（大文件上传可能需要更长时间，此处设5分钟）
	req.Close = true
	client := *c.httpClient
	client.Timeout = 5 * time.Minute // 覆盖默认超时

	// 发送上传请求
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送上传请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取上传响应失败: %w", err)
	}

	// 解析上传响应
	var apiResp APIResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return "", fmt.Errorf("解析上传响应失败，响应体: %s, 原因: %w", string(respBody), err)
	}

	// 检查上传结果
	if resp.StatusCode != http.StatusOK || apiResp.Code != 200 {
		return "", fmt.Errorf("上传失败，HTTP状态码: %d, 错误码: %d, 消息: %s",
			resp.StatusCode, apiResp.Code, apiResp.Message)
	}

	return fullRemotePath, nil
}

// GetFileInfo 获取文件信息（含下载地址）
// filePath: 远程文件路径（如 "/docs/test.txt"）
// 返回值: 文件信息结构体，错误信息
func (c *OpenListAPI) GetFileInfo(filePath string) (*FileInfo, error) {
	// 先检查登录状态
	if ok, err := c.Login(); !ok {
		if err != nil {
			return nil, fmt.Errorf("登录失败: %w", err)
		}
		return nil, fmt.Errorf("登录失败，无法获取文件信息")
	}

	// 构造请求体
	fileInfoReq := FileInfoRequest{
		Path:     filePath,
		Password: "",
	}

	// 执行请求
	fileInfo := &FileInfo{}
	if err := c.doRequest(&HTTPRequest{
		Method: "POST",
		URL:    fmt.Sprintf("%s/api/fs/get", c.baseURL),
		Body:   fileInfoReq,
	}, fileInfo); err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %w", err)
	}

	return fileInfo, nil
}

// SearchFiles 搜索文件
// keyword: 搜索关键词
// parentPath: 搜索父目录（默认 "/"）
// 返回值: 搜索结果列表，错误信息
func (c *OpenListAPI) SearchFiles(keyword, parentPath string) ([]SearchResult, error) {
	// 先检查登录状态
	if ok, err := c.Login(); !ok {
		if err != nil {
			return nil, fmt.Errorf("登录失败: %w", err)
		}
		return nil, fmt.Errorf("登录失败，无法执行文件搜索")
	}

	// 处理默认父目录（为空时设为 "/"）
	if parentPath == "" {
		parentPath = "/"
	}

	// 构造请求体
	searchReq := SearchRequest{
		Parent:   parentPath,
		Keywords: keyword,
	}

	// 执行请求
	var searchResults []SearchResult
	if err := c.doRequest(&HTTPRequest{
		Method: "POST",
		URL:    fmt.Sprintf("%s/api/fs/search", c.baseURL),
		Body:   searchReq,
	}, &searchResults); err != nil {
		return nil, fmt.Errorf("搜索文件失败: %w", err)
	}

	return searchResults, nil
}

// ListFiles 列出目录下的文件/目录
// path: 目录路径（默认 "/"）
// page: 页码（默认 1）
// perPage: 每页条数（0表示不分页）
// refresh: 是否强制刷新（默认 true）
// 返回值: 目录列表响应，错误信息
func (c *OpenListAPI) ListFiles(path string, page, perPage int, refresh bool) (*ListResponse, error) {
	// 先检查登录状态
	if ok, err := c.Login(); !ok {
		if err != nil {
			return nil, fmt.Errorf("登录失败: %w", err)
		}
		return nil, fmt.Errorf("登录失败，无法列出目录")
	}

	// 处理默认参数
	if path == "" {
		path = "/"
	}
	if page <= 0 {
		page = 1
	}

	// 构造请求体
	listReq := ListRequest{
		Path:     path,
		Password: "",
		Page:     page,
		PerPage:  perPage,
		Refresh:  refresh,
	}

	// 执行请求
	listResp := &ListResponse{}
	if err := c.doRequest(&HTTPRequest{
		Method: "POST",
		URL:    fmt.Sprintf("%s/api/fs/list", c.baseURL),
		Body:   listReq,
	}, listResp); err != nil {
		return nil, fmt.Errorf("列出目录失败: %w", err)
	}

	return listResp, nil
}

// getToken 获取当前登录令牌（带读锁，确保并发安全）
func (c *OpenListAPI) getToken() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.token
}
