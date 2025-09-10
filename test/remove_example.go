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

	// 删除单个文件
	fmt.Println("删除单个文件...")
	err := api.Remove("/remote/docs", []string{"test.txt"})
	if err != nil {
		fmt.Printf("删除文件失败: %v\n", err)
	} else {
		fmt.Println("文件删除成功!")
	}

	// 删除多个文件
	fmt.Println("删除多个文件...")
	err = api.Remove("/remote/docs", []string{"test1.txt", "test2.txt", "test3.txt"})
	if err != nil {
		fmt.Printf("删除文件失败: %v\n", err)
	} else {
		fmt.Println("文件删除成功!")
	}

	// 删除文件夹
	fmt.Println("删除文件夹...")
	err = api.Remove("/remote", []string{"temp"})
	if err != nil {
		fmt.Printf("删除文件夹失败: %v\n", err)
	} else {
		fmt.Println("文件夹删除成功!")
	}
}