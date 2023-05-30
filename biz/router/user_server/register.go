package user_server

import "github.com/gin-gonic/gin"

func Register(router *gin.Engine) {
	user := router.Group("/user")
	{
		user.POST("/login") // 用户登录
		user.POST("/check") // 用户鉴权
	}
}
