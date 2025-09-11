# OpenList Go Client

OpenList Go Client 是一个用于与 OpenList 文件管理服务进行交互的 Go 语言客户端库。它提供了简洁的 API 来执行文件上传、下载、搜索、删除和管理等操作。

## 功能特性

- 🔐 用户认证：自动处理登录和令牌管理
- 📁 文件管理：上传、下载、删除、获取文件信息、列出目录内容
- 🔍 文件搜索：根据关键词搜索文件
- 🔄 备份管理：自动备份目录并保留最新3份备份
- 🌐 代理支持：可配置 HTTP 代理
- 🔄 自动重试：登录状态自动维护
- 📦 易于集成：简洁的 API 设计，易于集成到您的 Go 项目中

## 安装

确保您已经安装了 Go 1.16 或更高版本。

```bash
go get -u github.com/littleboss01/openlistClient
```

或者在您的项目目录中初始化 Go 模块：

```bash
go mod init your-project-name
go get github.com/littleboss01/openlistClient
```

## 快速开始

```go
package main

import (
    "fmt"
    "log"
    "openlist"
)

func main() {
    // 创建客户端实例
    api := openlist.NewOpenListAPI(
        "http://localhost:5244", // OpenList服务地址
        "admin",                 // 用户名
        "123456",                // 密码
        "",                      // 代理地址（可选）
    )

    // 登录
    if ok, err := api.Login(); !ok {
        log.Fatal("登录失败:", err)
    }

    // 上传文件
    remotePath, err := api.UploadFile("/local/path/test.txt", "/remote/docs")
    if err != nil {
        log.Fatal("文件上传失败:", err)
    }
    fmt.Printf("文件上传成功，远程路径: %s\n", remotePath)

    // 获取文件信息
    fileInfo, err := api.GetFileInfo(remotePath)
    if err != nil {
        log.Fatal("获取文件信息失败:", err)
    }
    fmt.Printf("文件大小: %d字节，下载地址: %s\n", fileInfo.Size, fileInfo.URL)
}
```

## API 参考

### 创建客户端

```go
api := openlist.NewOpenListAPI(baseURL, username, password, proxy)
```

### 登录

```go
ok, err := api.Login()
```

### 上传文件

```go
remotePath, err := api.UploadFile(localFilePath, remoteDirectory)
```

### 下载文件（带进度回调）

```go
// 定义进度回调函数
progressFunc := func(downloaded, total int64) {
    fmt.Printf("下载进度: %d/%d bytes\n", downloaded, total)
}

// 下载文件
err := api.DownloadFile(remoteFilePath, localFilePath, progressFunc)
```

### 删除文件或文件夹

```go
// 删除单个文件
err := api.Remove("/remote/docs", []string{"test.txt"})

// 删除多个文件
err := api.Remove("/remote/docs", []string{"test1.txt", "test2.txt"})

// 删除文件夹
err := api.Remove("/remote", []string{"docs"})
```

### 备份目录并保留最新3份

```go
// 备份目录并自动管理备份文件
// 该功能会上传新的备份文件，并自动删除旧的备份文件，只保留最新的3份
err := backupExample() // 参见test/backup_example.go
```

### 检测目录并下载最新版本

```go
// 检测目录中的版本文件，找出最新版本并下载
err := versionCheckExample() // 参见test/version_check_example.go
```

### 获取文件信息

```go
fileInfo, err := api.GetFileInfo(filePath)
```

### 搜索文件

```go
results, err := api.SearchFiles(keyword, parentPath)
```

### 列出目录内容

```go
listResp, err := api.ListFiles(path, page, perPage, refresh)
```

## 错误处理

所有 API 方法都会返回详细的错误信息，您可以根据需要进行处理：

```go
if _, err := api.Login(); err != nil {
    // 处理登录错误
    log.Printf("登录失败: %v", err)
}
```

## 许可证

MIT License

## 参考
https://openlist.apifox.cn/
https://github.com/OpenListTeam/OpenList

## 贡献

欢迎提交 Issue 和 Pull Request 来改进这个项目。