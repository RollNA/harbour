package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ParseJWT func(token string) (map[string]any, bool)

func ParseJWTToken(parse ParseJWT) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
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

func ValidateJWTToken(parse ParseJWT) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
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
