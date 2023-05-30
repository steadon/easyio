package action_server

import "github.com/gin-gonic/gin"

func Register(router *gin.Engine) {
	image := router.Group("/action")
	{
		image.GET("/show/root", ShowRoot)      // 检索根下目录
		image.GET("/show/dir", ShowDir)        // 查看目录列表
		image.GET("/show/img", ShowImg)        // 查看图片列表
		image.DELETE("/delete/dir", DeleteDir) // 删除指定目录
		image.DELETE("/delete/img", DeleteImg) // 删除指定图片
		image.POST("/upload", UploadImg)       // 上传一张图片
	}
}
