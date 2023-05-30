package middleware

import "github.com/gin-gonic/gin"

// CorsMiddleware 跨域请求处理中间件
func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置允许跨域请求的头部字段
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			// 处理预检请求，返回空响应
			c.AbortWithStatus(200)
			return
		}
		// 继续处理其他请求
		c.Next()
	}
}
