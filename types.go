package openlist

// FileInfo 文件信息结构体（对应原Python的Dict返回）
type FileInfo struct {
	Path     string `json:"path"`     // 文件路径
	Name     string `json:"name"`     // 文件名
	Size     int64  `json:"size"`     // 文件大小（字节）
	IsDir    bool   `json:"is_dir"`   // 是否为目录
	URL      string `json:"url"`      // 下载地址
	Modified string `json:"modified"` // 修改时间
}

// SearchResult 搜索结果结构体
type SearchResult struct {
	Path     string `json:"path"`     // 文件路径
	Name     string `json:"name"`     // 文件名
	Size     int64  `json:"size"`     // 文件大小
	IsDir    bool   `json:"is_dir"`   // 是否为目录
	Modified string `json:"modified"` // 修改时间
}

// ListResponse 目录列表响应结构体
type ListResponse struct {
	Items   []FileInfo `json:"items"`    // 文件/目录列表
	Total   int        `json:"total"`    // 总数量
	Page    int        `json:"page"`     // 当前页码
	PerPage int        `json:"per_page"` // 每页条数
}

// LoginResponse 登录接口响应结构体
type LoginResponse struct {
	Code    int    `json:"code"`    // 状态码（200为成功）
	Message string `json:"message"` // 描述信息
	Data    struct {
		Token string `json:"token"` // 登录令牌
	} `json:"data"` // 业务数据
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
}

// ListRequest 列表请求参数
type ListRequest struct {
	Path     string `json:"path"`
	Password string `json:"password"`
	Page     int    `json:"page"`
	PerPage  int    `json:"per_page"`
	Refresh  bool   `json:"refresh"`
}

// HTTPRequest 通用HTTP请求配置
type HTTPRequest struct {
	Method  string            // HTTP方法 (GET, POST, PUT, DELETE等)
	URL     string            // 请求URL
	Body    interface{}       // 请求体数据（会自动序列化为JSON）
	Headers map[string]string // 请求头
}
