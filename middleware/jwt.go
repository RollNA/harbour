package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ParseJWT func(token string) (map[string]any, bool)

func ParseJWTToken(key string, parse ParseJWT) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if len(key) == 0 {
			key = "Authorization"
		}
		tokenString := ctx.GetHeader(key)

		if strings.Contains(tokenString, "Bearer ") {
			tokenString = strings.Split(tokenString, "Bearer ")[1]
		}
		input, valid := parse(tokenString)
		if valid {
			for k, v := range input {
				ctx.Set(k, v)
			}
			ctx.Set("JWTToken", tokenString)
		}
		ctx.Next()
	}
}

func ValidateJWTToken(key string, parse ParseJWT) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if len(key) == 0 {
			key = "Authorization"
		}
		tokenString := ctx.GetHeader(key)
		if strings.Contains(tokenString, "Bearer ") {
			tokenString = strings.Split(tokenString, "Bearer ")[1]
		}
		input, valid := parse(tokenString)
		if !valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			ctx.Abort()
			return
		}
		for k, v := range input {
			ctx.Set(k, v)
		}
		ctx.Set("JWTToken", tokenString)
		ctx.Next()
	}
}
