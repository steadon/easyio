package user

import (
	"EasyIO/biz/dal/model"
	"EasyIO/biz/dal/mysql"
	"EasyIO/biz/dal/param"
	"EasyIO/biz/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Sign 用户注册
func Sign(c *gin.Context) {
	var req param.Sign
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// 校验用户
	check := mysql.QueryUserByName(req.Username)
	if check != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "该用户名已被注册"})
		return
	}
	user := &model.User{
		Username: req.Username,
		Password: req.Password,
		PhoneNum: req.PhoneNum,
	}

	// 新增用户
	err := mysql.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "新增用户失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}

// Login 用户登录
func Login(c *gin.Context) {
	var req param.Login
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// 校验用户
	check := mysql.QueryUserByName(req.Username)
	if check == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "用户名尚未注册"})
		return
	}

	// 校验密码
	if req.Password != check.Password {
		c.JSON(http.StatusBadRequest, gin.H{"message": "用户名或密码错误"})
		return
	}

	// 签发令牌
	token, err := middleware.GenerateToken(int64(check.ID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "令牌签发错误"})
	}

	// 返回令牌
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func Info(c *gin.Context) {
	// 鉴权
	token := c.GetHeader("Authorization")
	load, err := middleware.VerifyAndParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "身份过期"})
		return
	}
	// 获取ID
	ID, _ := middleware.ExtractUserIDFromToken(load)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "token解析失败"})
		return
	}
	// 查询用户名
	username := mysql.QueryUserByID(ID).Username
	// 返回用户名
	c.JSON(http.StatusOK, gin.H{"username": username})
}
