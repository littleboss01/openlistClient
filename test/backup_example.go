package test

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"openlist"
	"path/filepath"
	"sort"
	"time"
)

// BackupExample 演示备份功能
func BackupExample() {
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

	// 要备份的目录（示例目录）
	sourceDir := "/data"
	// 备份文件存储目录
	backupDir := "/backups"
	
	// 确保备份目录存在
	err := ensureBackupDir(api, backupDir)
	if err != nil {
		fmt.Printf("确保备份目录存在失败: %v\n", err)
		return
	}
	
	// 生成备份文件名（包含时间戳）
	backupFileName := fmt.Sprintf("backup_%s.zip", time.Now().Format("20060102_150405"))
	
	// 创建本地备份文件
	localBackupPath := "./" + backupFileName
	fmt.Printf("创建备份文件: %s\n", localBackupPath)
	
	err = createBackupFile(localBackupPath, sourceDir)
	if err != nil {
		fmt.Printf("创建备份文件失败: %v\n", err)
		return
	}
	
	// 上传备份文件到OpenList
	fmt.Println("上传备份文件...")
	remotePath, err := api.UploadFile(
		localBackupPath, // 本地备份文件路径
		backupDir,       // 远程备份目录
	)
	
	if err != nil {
		fmt.Printf("上传备份文件失败: %v\n", err)
		// 清理本地文件
		os.Remove(localBackupPath)
		return
	}
	
	fmt.Printf("备份文件上传成功，远程路径: %s\n", remotePath)
	
	// 清理本地文件
	os.Remove(localBackupPath)
	
	// 检查备份目录中的文件列表
	fmt.Println("检查备份目录中的文件...")
	listResp, err := api.ListFiles(backupDir, 1, 0, true) // 获取所有文件，不分页
	if err != nil {
		fmt.Printf("列出备份目录失败: %v\n", err)
		return
	}
	
	// 筛选出备份文件（以backup_开头的文件）
	var backupFiles []openlist.FileInfo
	for _, item := range listResp.Items {
		if !item.IsDir && len(item.Name) > 7 && item.Name[:7] == "backup_" {
			backupFiles = append(backupFiles, item)
		}
	}
	
	fmt.Printf("找到 %d 个备份文件\n", len(backupFiles))
	
	// 如果备份文件超过3个，删除最旧的
	if len(backupFiles) > 3 {
		// 按修改时间排序，最新的在前
		sort.Slice(backupFiles, func(i, j int) bool {
			return backupFiles[i].Modified > backupFiles[j].Modified
		})
		
		// 删除超过3个的旧备份文件
		for i := 3; i < len(backupFiles); i++ {
			fileToDelete := backupFiles[i]
			fmt.Printf("删除旧备份文件: %s\n", fileToDelete.Name)
			
			// 删除文件
			err := api.Remove(backupDir, []string{fileToDelete.Name})
			if err != nil {
				fmt.Printf("删除备份文件失败: %v\n", err)
				// 继续处理其他文件
			} else {
				fmt.Printf("备份文件删除成功: %s\n", fileToDelete.Name)
			}
		}
	} else {
		fmt.Println("备份文件数量未超过3个，无需清理")
	}
	
	fmt.Println("备份演示完成")
}

// ensureBackupDir 确保备份目录存在
func ensureBackupDir(api *openlist.OpenListAPI, backupDir string) error {
	// 尝试列出目录来检查是否存在
	_, err := api.ListFiles(backupDir, 1, 1, true)
	if err != nil {
		// 目录可能不存在，这里我们假设目录会自动创建
		// 在实际应用中，您可能需要根据具体的服务行为进行处理
		fmt.Printf("备份目录可能不存在，将在上传时自动创建: %s\n", backupDir)
	}
	return nil
}

// createBackupFile 创建备份文件（示例实现）
func createBackupFile(backupFilePath, sourceDir string) error {
	// 创建备份文件
	file, err := os.Create(backupFilePath)
	if err != nil {
		return fmt.Errorf("创建备份文件失败: %w", err)
	}
	defer file.Close()

	// 创建zip writer
	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()

	// 添加一些示例文件到备份中（在实际应用中，您需要遍历sourceDir目录）
	// 这里我们只添加一个示例文件
	err = addFileToZip(zipWriter, "backup_info.txt", "这是备份文件的示例内容\n备份时间: "+time.Now().Format("2006-01-02 15:04:05")+"\n")
	if err != nil {
		return fmt.Errorf("添加文件到备份失败: %w", err)
	}

	// 添加另一个示例文件
	err = addFileToZip(zipWriter, "source_dir_info.txt", "源目录: "+sourceDir+"\n")
	if err != nil {
		return fmt.Errorf("添加文件到备份失败: %w", err)
	}

	fmt.Printf("备份文件创建成功: %s\n", backupFilePath)
	return nil
}

// addFileToZip 添加文件到zip归档
func addFileToZip(zipWriter *zip.Writer, filename, content string) error {
	// 创建文件头
	header := &zip.FileHeader{
		Name:   filename,
		Method: zip.Deflate,
	}
	header.SetModTime(time.Now())

	// 创建文件
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	// 写入内容
	_, err = io.WriteString(writer, content)
	return err
}