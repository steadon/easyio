package middleware

import (
	"EasyIO/biz/pkg/setting"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// GenerateToken 签发 JWT 令牌
func GenerateToken(userID int64) (string, error) {
	// 创建声明（Claims）
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(), // 设置过期时间为 24 小时
	}

	// 创建令牌（Token）
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名令牌
	secret := []byte(setting.SIGN) // 设置自己的密钥
	signedToken, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	// 返回令牌
	return signedToken, nil
}

// VerifyAndParseToken 校验和解析 JWT 令牌
func VerifyAndParseToken(tokenString string) (*jwt.Token, error) {
	// 解析令牌
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 校验签名算法是否匹配
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}

		// 返回密钥
		return []byte(setting.SIGN), nil // 设置与签发时相同的密钥
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

// ExtractUserIDFromToken 从令牌中提取用户ID
func ExtractUserIDFromToken(token *jwt.Token) (int64, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("invalid token claims")
	}

	userID, ok := claims["sub"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid user ID")
	}

	return int64(userID), nil
}

// CheckRole 鉴权
func CheckRole(c *gin.Context) {
	// 鉴权
	token := c.GetHeader("Authorization")
	_, err := VerifyAndParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "身份过期"})
		return
	}
}
