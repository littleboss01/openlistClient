package test

import (
	"fmt"
	"openlist"
)

// ClientTest 测试客户端基本功能
func ClientTest() {
	// 1. 创建客户端实例
	api := openlist.NewOpenListAPI(
		"http://localhost:5244", // OpenList服务地址
		"admin",                 // 用户名
		"123456",                // 密码
		"http://127.0.0.1:8080", // 代理地址（可选，为空则不使用代理）
	)

	// 2. 测试代理（可选）
	if api.TestProxy() {
		fmt.Println("代理可用")
	} else {
		fmt.Println("代理不可用")
	}

	// 3. 登录（自动触发，也可手动调用api.Login()）
	_, err := api.Login()
	if err != nil {
		fmt.Println("登录失败")
		return
	}

	// 4. 上传文件
	remotePath, err := api.UploadFile(
		"/local/path/test.txt", // 本地文件路径
		"/remote/docs",         // 远程目录
	)
	if err != nil {
		fmt.Printf("文件上传失败: %v\n", err)
		return
	}
	fmt.Printf("文件上传成功，远程路径: %s\n", remotePath)

	// 5. 获取文件信息
	fileInfo, err := api.GetFileInfo(remotePath)
	if err != nil {
		fmt.Printf("获取文件信息失败: %v\n", err)
		return
	}
	fmt.Printf("文件大小: %d字节，下载地址: %s\n", fileInfo.Size, fileInfo.URL)

	// 6. 搜索文件
	results, err := api.SearchFiles("test", "/remote/docs")
	if err != nil {
		fmt.Printf("文件搜索失败: %v\n", err)
		return
	}
	fmt.Printf("搜索到%d个结果:\n", len(results))
	for _, res := range results {
		fmt.Printf("  %s (是否目录: %t)\n", res.Path, res.IsDir)
	}

	// 7. 列出目录
	listResp, err := api.ListFiles("/remote/docs", 1, 10, true)
	if err != nil {
		fmt.Printf("列出目录失败: %v\n", err)
		return
	}
	fmt.Printf("目录总文件数: %d，当前页: %d\n", listResp.Total, listResp.Page)

	// 演示使用新的请求参数结构体

	// 演示直接使用通用HTTP请求方法
	fmt.Println("\n=== 使用通用HTTP请求方法 ===")

	// 使用通用HTTP请求方法获取文件信息
	fileInfo2, err := api.GetFileInfo("/remote/docs/test.txt")
	if err != nil {
		fmt.Printf("使用通用HTTP请求方法获取文件信息失败: %v\n", err)
	} else {
		fmt.Printf("使用通用HTTP请求方法获取文件信息成功\n", fileInfo2)
	}
}
