package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gwid.io/gwid-core/internals/config"
	"gwid.io/gwid-core/internals/types"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.GetHeader("Authorization")

		if bearerToken == "" || !strings.HasPrefix(bearerToken, "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "unauthorized",
			})

			return
		}

		tokenString := extractToken(bearerToken)

		token, err := jwt.ParseWithClaims(tokenString, &types.JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.GetEnv("JWT_SECRET", "the-fallback-key")), nil
		})
		if err != nil {

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "unauthorized",
			})

			return
		}

		if claims, ok := token.Claims.(*types.JwtCustomClaims); ok {
			c.Set("user", claims)
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "unauthorized",
			})

			return
		}

		c.Next()
	}
}

func extractToken(bearerToken string) string {
	bearerTokenSlice := strings.Split(bearerToken, " ")

	if len(bearerTokenSlice) == 2 {
		return bearerTokenSlice[1]
	}

	return ""
}
