package action_server

import (
	"EasyIO/biz/dal/param"
	"EasyIO/biz/dal/result"
	"EasyIO/biz/pkg/setting"
	"crypto/rand"
	"encoding/base64"
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
	extension := filepath.Ext(file.Filename)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件上传失败"})
		return
	}

	// 创建目标目录
	folder := c.PostForm("group") // 指定目录分层
	folderPath := filepath.Join(setting.StorageDir, folder)

	// 获取文件名称
	name := c.PostForm("name")

	// 不传name则使用随机串
	if name == "" {
		name, _ = generateRandomString(8)
	}

	// 如果有后缀则去掉后缀
	if strings.ContainsAny(name, ".") {
		name = getFileNameWithoutExtension(name)
	}

	// 拼接图片全名
	fullName := name + extension

	if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "目录创建失败"})
		return
	}

	// 创建目标文件路径
	dstPath := filepath.Join(folderPath, fullName)

	// 创建一个新文件来保存上传的文件
	dst, err := os.Create(dstPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件创建失败"})
		log.Printf("err: %v\n", err)
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
	filePaths := make([]result.Image, 0)

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Failed to walk directory: %v\n", err)
			return nil
		}

		// 如果是文件，将文件路径添加到切片中
		if !info.IsDir() {
			path := setting.Prefix + strings.Replace(path, setting.StorageDir, setting.Proxy, 1)
			index := strings.LastIndex(path, "/")
			if index == -1 {
				log.Printf("Invalid param: %v\n", err)
			}
			image := result.Image{Name: path[index+1:], Path: path}
			filePaths = append(filePaths, image)
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

// AddDir 创建目录
func AddDir(c *gin.Context) {
	var req param.Group

	// 创建目标目录
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	folderPath := filepath.Join(setting.StorageDir, req.Name)

	if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "目录创建失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "目录创建成功"})
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

// ShowRoot 获取目录
func ShowRoot(c *gin.Context) {
	// 指定目录路径
	dirPath := setting.StorageDir

	// 创建切片返回响应
	dirPaths := make([]result.Root, 0)

	// 遍历目录下的所有文件和子目录
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// 处理遍历过程中的错误
			log.Printf("Encountered error: %v\n", err)
			return nil
		}

		if info.IsDir() && path != dirPath {
			// 获取相对路径
			relPath, err := filepath.Rel(dirPath, path)
			if err != nil {
				return err
			}

			// 计算层级(通过相对路径计算分隔符)
			level := strings.Count(relPath, string(os.PathSeparator))

			// 封装后添加的返回列表
			dirInfo := result.Root{
				Name:  filepath.Base(path),
				Level: level,
			}
			dirPaths = append(dirPaths, dirInfo)
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

// ShowDir 获取目录
func ShowDir(c *gin.Context) {
	// 指定目录路径
	path := c.Query("group")

	// 拼接完整目录路径
	dirPath := filepath.Join(setting.StorageDir, path)

	// 创建切片返回响应
	dirPaths := make([]result.Group, 0)

	// 打开指定目录
	dir, err := os.Open(dirPath)
	if err != nil {
		// 处理打开目录失败的错误
		log.Printf("Failed to open directory: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "打开目录失败"})
		return
	}
	defer dir.Close()

	// 读取目录中的文件和子目录
	fileInfo, err := dir.Readdir(-1)
	if err != nil {
		// 处理读取目录内容失败的错误
		log.Printf("Failed to read directory: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取目录失败"})
		return
	}

	// 遍历目录中的文件和子目录
	for _, info := range fileInfo {
		if info.IsDir() {
			// 添加子目录名称
			dirName := result.Group{Name: info.Name()}
			dirPaths = append(dirPaths, dirName)
		}
	}
	c.JSON(http.StatusOK, gin.H{"dirPaths": dirPaths})
}

// 生成随机字符串
func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	randomString := base64.URLEncoding.EncodeToString(bytes)[:length]
	randomString = strings.ReplaceAll(randomString, "-", "")
	randomString = strings.ReplaceAll(randomString, "_", "")
	return randomString, nil
}

// 去掉图片后缀名
func getFileNameWithoutExtension(filePath string) string {
	return strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
}
