package main

import (
	"fmt"
	"openlist"
	"os"
)

func main() {
	// 创建客户端实例
	api := openlist.NewOpenListAPI(
		"http://localhost:5244", // OpenList服务地址
		"admin",                 // 用户名
		"123456",                // 密码
		"",                      // 代理地址（可选）
	)

	// 定义进度回调函数
	progressFunc := func(downloaded, total int64) {
		if total > 0 {
			percentage := float64(downloaded) / float64(total) * 100
			fmt.Printf("\r下载进度: %.2f%% (%d/%d bytes)", percentage, downloaded, total)
		} else {
			fmt.Printf("\r已下载: %d bytes", downloaded)
		}
	}

	// 确保下载目录存在
	err := os.MkdirAll("./downloads", 0755)
	if err != nil {
		fmt.Printf("创建下载目录失败: %v\n", err)
		return
	}

	// 下载文件
	fmt.Println("开始下载文件...")
	
	err = api.DownloadFile(
		"/test/example.txt",          // 远程文件路径
		"./downloads/example.txt",    // 本地保存路径
		progressFunc,                 // 进度回调函数
	)
	
	if err != nil {
		fmt.Printf("\n文件下载失败: %v\n", err)
		return
	}

	fmt.Println("\n文件下载成功!")
	
	// 验证文件是否存在
	if _, err := os.Stat("./downloads/example.txt"); err == nil {
		fmt.Println("文件验证成功: ./downloads/example.txt")
	} else {
		fmt.Printf("文件验证失败: %v\n", err)
	}
}