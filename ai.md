# OpenList Go Client 项目使用方法总结

## 项目概述
OpenList Go Client 是一个用于与 OpenList 文件管理服务进行交互的 Go 语言客户端库，提供了文件上传、下载、搜索、删除和管理等功能。

## 核心功能
1. **用户认证**：自动处理登录和令牌管理
2. **文件管理**：上传、下载、删除、获取文件信息、列出目录内容
3. **文件搜索**：根据关键词搜索文件
4. **备份管理**：自动备份目录并保留最新3份备份
5. **代理支持**：可配置 HTTP 代理

## 安装方法
```bash
go get -u github.com/littleboss01/openlistClient
```

## 基本使用方法

### 1. 创建客户端实例
```go
api := openlist.NewOpenListAPI(
    "http://localhost:5244", // OpenList服务地址
    "admin",                 // 用户名
    "123456",                // 密码
    "",                      // 代理地址（可选）
)
```

### 2. 登录
```go
ok, err := api.Login()
```

### 3. 文件上传
```go
remotePath, err := api.UploadFile("/local/path/test.txt", "/remote/docs")
```

### 4. 文件下载（带进度回调）
```go
// 定义进度回调函数
progressFunc := func(downloaded, total int64) {
    fmt.Printf("下载进度: %d/%d bytes\n", downloaded, total)
}

// 下载文件
err := api.DownloadFile("/remote/docs/test.txt", "./downloaded_test.txt", progressFunc)
```

### 5. 删除文件或文件夹
```go
// 删除单个文件
err := api.Remove("/remote/docs", []string{"test.txt"})

// 删除多个文件
err := api.Remove("/remote/docs", []string{"test1.txt", "test2.txt"})

// 删除文件夹
err := api.Remove("/remote", []string{"docs"})
```

### 6. 获取文件信息
```go
fileInfo, err := api.GetFileInfo("/remote/docs/test.txt")
```

### 7. 搜索文件
```go
results, err := api.SearchFiles("keyword", "/remote/docs")
```

### 8. 列出目录内容
```go
listResp, err := api.ListFiles("/remote/docs", 1, 10, true)
```

## 主要数据结构

### FileInfo（文件信息）
```go
type FileInfo struct {
    Path     string // 文件路径
    Name     string // 文件名
    Size     int64  // 文件大小（字节）
    IsDir    bool   // 是否为目录
    URL      string // 下载地址
    Modified string // 修改时间
}
```

### SearchResult（搜索结果）
```go
type SearchResult struct {
    Path     string // 文件路径
    Name     string // 文件名
    Size     int64  // 文件大小
    IsDir    bool   // 是否为目录
    Modified string // 修改时间
}
```

## 错误处理
所有 API 方法都会返回详细的错误信息，应根据需要进行处理：
```go
if _, err := api.Login(); err != nil {
    // 处理登录错误
    log.Printf("登录失败: %v", err)
}
```

## 特殊功能示例

### 备份目录管理
备份目录并自动管理备份文件，上传新的备份文件，并自动删除旧的备份文件，只保留最新的3份。

### 版本检测与下载
检测目录中的版本文件，找出最新版本并下载。