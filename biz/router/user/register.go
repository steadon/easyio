package user

import "github.com/gin-gonic/gin"

func Register(router *gin.Engine) {
	user := router.Group("/user")
	{
		user.POST("/sign", Sign)   // 用户注册
		user.POST("/login", Login) // 用户登录
		user.GET("/info", Info)    // 用户登录
	}
}
