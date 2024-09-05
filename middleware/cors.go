package middleware

import (
	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != "OPTIONS" {
			// 接收客户端发送的origin
			c.Header("Access-Control-Allow-Origin", "*")
			// 服务器支持的所有跨域请求的方法
			c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			// 允许跨域设置可以返回其他子段，可以自定义字段
			c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
			c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
			c.Header("Content-Type", "application/json")
			c.Next()
		} else {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
			c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
			c.Header("Content-Type", "application/json")
			// 设置缓存时间
			c.Header("Access-Control-Max-Age", "3600")
			c.AbortWithStatus(200)
		}
	}
}
