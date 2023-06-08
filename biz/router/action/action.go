package action

import (
	"EasyIO/biz/dal/param"
	"EasyIO/biz/dal/result"
	"EasyIO/biz/middleware"
	"EasyIO/biz/pkg/setting"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// UploadImg 上传图片
func UploadImg(c *gin.Context) {
	// 鉴权
	if middleware.CheckRole(c) == false {
		return
	}
	// 从请求中获取上传的文件
	file, err := c.FormFile("file")
	extension := filepath.Ext(file.Filename)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "文件上传失败"})
		return
	}
	// 创建目标目录
	folder := c.PostForm("group") // 指定目录分层
	folderPath := filepath.Join(setting.StorageDir, folder)

	// 删除图片缓存
	go middleware.DeleteCache(folder)

	// 获取文件名称
	name := c.PostForm("name")

	// 定义正则表达式模式
	pattern := "[\\p{P}\\s]+" // 匹配标点符号和空格

	// 编译正则表达式
	reg := regexp.MustCompile(pattern)

	// 检查字符串是否含有标点符号或空格
	if reg.MatchString(name) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "不合法的文件名"})
	}
	// 不传name则使用随机串
	if name == "" {
		name, _ = middleware.GenerateRandomString(8)
	}
	// 如果有后缀则去掉后缀
	if strings.ContainsAny(name, ".") {
		name = middleware.GetFileNameWithoutExtension(name)
	}
	// 拼接图片全名
	fullName := name + extension

	// 递归创建目录
	if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "目录创建失败"})
		return
	}
	// 创建目标文件路径
	dstPath := filepath.Join(folderPath, fullName)

	// 创建一个新文件来保存上传的文件
	dst, err := os.Create(dstPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "文件创建失败"})
		log.Printf("err: %v\n", err)
		return
	}
	defer dst.Close()

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "无法打开文件"})
		return
	}
	defer src.Close()

	// 将上传的文件内容拷贝到目标文件
	if _, err := io.Copy(dst, src); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "文件拷贝失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "文件上传成功"})
}

// ShowImg 获取图片
func ShowImg(c *gin.Context) {
	// 鉴权
	if middleware.CheckRole(c) == false {
		return
	}
	// 获取参数并尝试从缓存获取结果
	group := c.Query("group")
	filePathsCache := middleware.GetCache(group)
	if filePathsCache != nil {
		c.JSON(http.StatusOK, gin.H{"filePaths": filePathsCache})
		return
	}
	// 拼接文件夹路径
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
		c.JSON(http.StatusInternalServerError, gin.H{"message": "文件遍历失败"})
		return
	}
	// 添加缓存
	go middleware.AddCache(group, filePaths)
	c.JSON(http.StatusOK, gin.H{"filePaths": filePaths})
}

// DeleteImg 删除图片
func DeleteImg(c *gin.Context) {
	// 鉴权
	if middleware.CheckRole(c) == false {
		return
	}
	// 获取文件或目录路径
	path := c.Query("path")

	// 删除图片所在缓存区的缓存
	index := strings.LastIndex(path, "/")
	go middleware.DeleteCache(path[:index])

	// 拼接完整路径
	filePath := filepath.Join(setting.StorageDir, path)

	// 删除单个文件
	err := os.Remove(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "删除文件失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "文件删除成功"})
}

// AddDir 创建目录
func AddDir(c *gin.Context) {
	// 鉴权
	if middleware.CheckRole(c) == false {
		return
	}
	var req param.Group
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	// 递归创建目标目录
	folderPath := filepath.Join(setting.StorageDir, req.Name)
	if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "目录创建失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "目录创建成功"})
}

// DeleteDir 删除目录
func DeleteDir(c *gin.Context) {
	// 鉴权
	if middleware.CheckRole(c) == false {
		return
	}
	// 获取文件或目录路径
	path := c.Query("path")

	// 拼接完整路径
	dirPath := filepath.Join(setting.StorageDir, path)

	// 删除目录
	err := os.RemoveAll(dirPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "删除目录失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "目录删除成功"})
}

// ShowRoot 获取目录
func ShowRoot(c *gin.Context) {
	// 鉴权
	if middleware.CheckRole(c) == false {
		return
	}
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
		c.JSON(http.StatusInternalServerError, gin.H{"message": "遍历目录失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"dirPaths": dirPaths})
}

// ShowDir 获取目录
func ShowDir(c *gin.Context) {
	// 鉴权
	if middleware.CheckRole(c) == false {
		return
	}
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
		c.JSON(http.StatusInternalServerError, gin.H{"message": "打开目录失败"})
		return
	}
	defer dir.Close()

	// 读取目录中的文件和子目录
	fileInfo, err := dir.Readdir(-1)
	if err != nil {
		// 处理读取目录内容失败的错误
		log.Printf("Failed to read directory: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "读取目录失败"})
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
