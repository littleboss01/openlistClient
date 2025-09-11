package test

import (
	"fmt"
	"sort"
	"strings"
	"time"

	openlist "github.com/littleboss01/openlistClient"
)

// VersionCheckExample 演示如何检测目录、判断最新版本并下载最新版本
func VersionCheckExample() {
	// 创建客户端实例
	api := openlist.NewOpenListAPI(
		"http://localhost:5244", // OpenList服务地址
		"admin",                 // 用户名
		"123456",                // 密码
		"",                      // 代理地址（可选）
	)

	// 登录
	if ok, err := api.Login(); !ok {
		fmt.Printf("登录失败: %v\n", err)
		return
	}

	// 要检查的目录路径（示例目录）
	checkDir := "/releases"

	// 列出目录中的文件
	fmt.Printf("检查目录 %s 中的文件...\n", checkDir)
	listResp, err := api.ListFiles(checkDir, 1, 0, true) // 获取所有文件，不分页
	if err != nil {
		fmt.Printf("列出目录失败: %v\n", err)
		return
	}

	// 筛选出版本文件（以v开头的文件，假设版本文件命名规则为 v1.0.0.zip, v1.1.0.zip 等）
	var versionFiles []openlist.FileInfo
	for _, item := range listResp.Content {
		if !item.IsDir && strings.HasPrefix(item.Name, "v") && strings.HasSuffix(item.Name, ".zip") {
			versionFiles = append(versionFiles, item)
		}
	}

	if len(versionFiles) == 0 {
		fmt.Println("未找到版本文件")
		return
	}

	fmt.Printf("找到 %d 个版本文件\n", len(versionFiles))

	// 按版本号排序，找出最新版本
	latestVersion := findLatestVersion(versionFiles)
	if latestVersion == nil {
		fmt.Println("无法确定最新版本")
		return
	}

	fmt.Printf("最新版本: %s, 修改时间: %s\n", latestVersion.Name, latestVersion.Modified)

	// 下载最新版本
	fmt.Printf("开始下载最新版本 %s...\n", latestVersion.Name)
	localPath := "./" + latestVersion.Name
	startTime := time.Now()

	// 定义进度回调函数
	progressFunc := func(downloaded, total int64) {
		if total > 0 {
			percentage := float64(downloaded) / float64(total) * 100
			fmt.Printf("下载进度: %.2f%% (%d/%d bytes)\n", percentage, downloaded, total)
		} else {
			fmt.Printf("已下载: %d bytes\n", downloaded)
		}
	}

	// 构造远程文件路径
	remotePath := fmt.Sprintf("%s/%s", checkDir, latestVersion.Name)

	// 下载文件
	err = api.DownloadFile(remotePath, localPath, progressFunc)
	if err != nil {
		fmt.Printf("下载最新版本失败: %v\n", err)
		return
	}

	elapsedTime := time.Since(startTime)
	fmt.Printf("最新版本下载成功! 耗时: %v\n", elapsedTime)
	fmt.Printf("本地文件路径: %s\n", localPath)
}

// findLatestVersion 查找最新版本文件
func findLatestVersion(files []openlist.FileInfo) *openlist.FileInfo {
	if len(files) == 0 {
		return nil
	}

	// 按修改时间排序，最新的在前
	sort.Slice(files, func(i, j int) bool {
		//return files[i].Modified > files[j].Modified
		return files[i].Modified.After(files[j].Modified)
	})

	// 返回最新版本
	return &files[0]
}
