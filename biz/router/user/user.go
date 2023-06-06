package user

import (
	"EasyIO/biz/dal/param"
	"EasyIO/biz/middleware"
	"EasyIO/biz/pkg/setting"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Login 用户登录
func Login(c *gin.Context) {
	var req param.Login

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// 校验用户名和密码
	if req.Username != setting.Username || req.Password != setting.Password {
		c.JSON(http.StatusBadRequest, gin.H{"message": "不存在该用户"})
		return
	}

	// 签发令牌
	token, err := middleware.GenerateToken(req.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "令牌签发错误"})
	}

	// 返回令牌
	c.JSON(http.StatusOK, gin.H{"token": token})
}
