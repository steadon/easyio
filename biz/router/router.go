package router

import (
	"EasyIO/biz/middleware"
	"EasyIO/biz/pkg/setting"
	"EasyIO/biz/router/action"
	"EasyIO/biz/router/user"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	// 设置跨域请求处理中间件
	router.Use(middleware.CorsMiddleware())

	// 指定图片目录
	imgDir := setting.StorageDir
	// 注册静态文件处理中间件
	router.Static(setting.Proxy, imgDir)
	// 注册接口
	user.Register(router)
	action.Register(router)

	return router
}
