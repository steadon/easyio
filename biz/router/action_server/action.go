package action_server

import (
	"EasyIO/biz/pkg/setting"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// UploadImg 上传图片
func UploadImg(c *gin.Context) {
	// 从请求中获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件上传失败"})
		return
	}

	// 创建目标目录
	folder := c.PostForm("group") // 指定目录分层
	folderPath := filepath.Join(setting.StorageDir, folder)

	if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "目录创建失败"})
		return
	}

	// 创建目标文件路径
	dstPath := filepath.Join(folderPath, file.Filename)

	// 创建一个新文件来保存上传的文件
	dst, err := os.Create(dstPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件创建失败"})
		return
	}
	defer dst.Close()

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法打开文件"})
		return
	}
	defer src.Close()

	// 将上传的文件内容拷贝到目标文件
	if _, err := io.Copy(dst, src); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件拷贝失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "文件上传成功"})
}

// ShowImg 下载图片
func ShowImg(c *gin.Context) {
	// 获取参数凭借路径
	group := c.Query("group")
	folderPath := filepath.Join(setting.StorageDir, group)

	// 遍历文件夹下的所有文件
	filePaths := make([]string, 0)

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Failed to walk directory: %v\n", err)
			return nil
		}

		// 如果是文件，将文件路径添加到切片中
		if !info.IsDir() {
			filePaths = append(filePaths, setting.Prefix+strings.Replace(path, setting.StorageDir, setting.Proxy, 1))
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件遍历失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"filePaths": filePaths})
}

// DeleteImg 删除图片
func DeleteImg(c *gin.Context) {
	// 获取文件或目录路径
	path := c.Query("path")

	// 拼接完整路径
	filePath := filepath.Join(setting.StorageDir, path)

	// 删除文件
	err := os.Remove(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除文件失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "文件删除成功"})
}

// DeleteDir 删除目录
func DeleteDir(c *gin.Context) {
	// 获取文件或目录路径
	path := c.Query("path")

	// 拼接完整路径
	dirPath := filepath.Join(setting.StorageDir, path)

	// 删除目录
	err := os.RemoveAll(dirPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除目录失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "目录删除成功"})
}

// ShowDir 获取目录
func ShowDir(c *gin.Context) {
	// 指定目录路径
	dirPath := setting.StorageDir

	// 创建切片返回响应
	dirPaths := make([]string, 0)

	// 遍历目录下的所有文件和子目录
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// 处理遍历过程中的错误
			log.Printf("Encountered error: %v\n", err)
			return nil
		}
		if info.IsDir() && path != dirPath {
			// 获取子目录的名称
			dirName := filepath.Base(path)
			// 添加目录名称
			dirPaths = append(dirPaths, dirName)
		}
		return nil
	})

	if err != nil {
		// 处理遍历目录的错误
		log.Printf("Failed to walk directory: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "遍历目录失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"dirPaths": dirPaths})
}
