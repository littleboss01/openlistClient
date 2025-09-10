package test

import (
	"fmt"
	"time"

	openlist "github.com/littleboss01/openlistClient"
)

// DownloadExample 演示如何使用下载功能
func DownloadExample() {
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
			fmt.Printf("下载进度: %.2f%% (%d/%d bytes)\n", percentage, downloaded, total)
		} else {
			fmt.Printf("已下载: %d bytes\n", downloaded)
		}
	}

	// 下载文件
	fmt.Println("开始下载文件...")
	startTime := time.Now()

	err := api.DownloadFile(
		"/remote/docs/test.txt", // 远程文件路径
		"./downloaded_test.txt", // 本地保存路径
		progressFunc,            // 进度回调函数
	)

	if err != nil {
		fmt.Printf("文件下载失败: %v\n", err)
		return
	}

	elapsedTime := time.Since(startTime)
	fmt.Printf("文件下载成功! 耗时: %v\n", elapsedTime)
}
