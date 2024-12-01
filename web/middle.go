package web

import (
	"fmt"
	"net/http"
	"project/logic/controll"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware() gin.HandlerFunc {
	whitelist := map[string]bool{
		"/login":    true,
		"/refresh":  true,
		"/register": true,
	}
	return func(c *gin.Context) {
		if whitelist[c.Request.URL.Path] {
			c.Next()
			return
		}
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token错误"})
			return
		}
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}
		claims := &controll.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return controll.AccessSecret, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token解析失败"})
			return
		}
		// 将用户名存储在上下文中
		c.Set("username", claims.UserName)
		c.Set("userid", claims.UserId)
		c.Set("rolesId", claims.RoleId)
		c.Next()
	}
}
func RoleMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
