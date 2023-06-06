package user

import "github.com/gin-gonic/gin"

func Register(router *gin.Engine) {
	user := router.Group("/user")
	{
		user.POST("/login", Login) // 用户登录
	}
}
