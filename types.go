package openlist

import "time"

// ProgressFunc 进度回调函数类型
// 参数: 已下载字节数, 总字节数
type ProgressFunc func(downloaded, total int64)

// FileInfo 文件信息结构体（对应原Python的Dict返回）
type FileInfo struct {
	Name      string      `json:"name"`     // 文件名
	Size      int64       `json:"size"`     // 文件大小（字节）
	IsDir     bool        `json:"is_dir"`   // 是否为目录
	Modified  time.Time   `json:"modified"` // 修改时间
	Created   time.Time   `json:"created"`
	Sign      string      `json:"sign"`
	Thumb     string      `json:"thumb"`
	Type      int64       `json:"type"`
	HashInfo  interface{} `json:"hashinfo"`
	Hash_info interface{} `json:"hash_info"`

	//list没有get才有
	Raw_url string      `json:"raw_url"`
	Related interface{} `json:"related"`
}

// SearchResult 搜索结果结构体
type SearchResult struct {
	Content []struct {
		Parent string `json:"parent"`
		Name   string `json:"name"`   // 文件名
		IsDir  bool   `json:"is_dir"` // 是否为目录
		Size   int64  `json:"size"`   // 文件大小（字节）
		Type   int64  `json:"type"`
	}
}

// ListResponse 目录列表响应结构体
type ListResponse struct {
	Content  []FileInfo `json:"content"`  // 文件/目录列表
	Total    int        `json:"total"`    // 总数量
	Page     int        `json:"page"`     // 当前页码
	PerPage  int        `json:"per_page"` // 每页条数
	Write    bool       `json:"write"`
	Provider string     `json:"provider"`
	Readme   string     `json:"readme"`
	Header   string     `json:"header"`
}

// LoginResponse 登录接口响应结构体
type LoginResponse struct {
	Token string `json:"token"` // 登录令牌
}

// APIResponse 通用API响应结构体（用于解析非登录接口的返回）
type APIResponse struct {
	Code    int         `json:"code"`    // 状态码（200为成功）
	Message string      `json:"message"` // 描述信息
	Data    interface{} `json:"data"`    // 业务数据（动态解析）
}

// LoginRequest 登录请求参数
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UploadRequest 上传文件请求参数
type UploadRequest struct {
	FilePath   string `json:"file_path"`
	RemotePath string `json:"remote_path"`
}

// FileInfoRequest 获取文件信息请求参数
type FileInfoRequest struct {
	Path     string `json:"path"`
	Password string `json:"password"`
}

// SearchRequest 搜索文件请求参数
type SearchRequest struct {
	Parent   string `json:"parent"`
	Keywords string `json:"keywords"`
	Scope    int    `json:"scope" default:"0"`

	Page     int `json:"page" default:"1"`
	Per_page int `json:"per_page" default:"50"`
}

// ListRequest 列表请求参数
type ListRequest struct {
	Path     string `json:"path"`
	Password string `json:"password"`
	Page     int    `json:"page"`
	PerPage  int    `json:"per_page"`
	Refresh  bool   `json:"refresh"`
}

// RemoveRequest 删除文件或文件夹请求参数
type RemoveRequest struct {
	Dir   string   `json:"dir"`   // 目录
	Names []string `json:"names"` // 文件名列表
}

// HTTPRequest 通用HTTP请求配置
type HTTPRequest struct {
	Method  string            // HTTP方法 (GET, POST, PUT, DELETE等)
	URL     string            // 请求URL
	Body    interface{}       // 请求体数据（会自动序列化为JSON）
	Headers map[string]string // 请求头
}

// MkdirRequest 创建文件夹请求参数
type MkdirRequest struct {
	Path string `json:"path"` // 新目录路径
}
