package middleware

import (
	"strings"

	"github.com/Oralkhan-coder/mind-map/internal/core"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Error(core.Unauthorized("authorization header required"))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Error(core.Unauthorized("format must be Bearer {token}"))
			c.Abort()
			return
		}

		token, err := core.ValidateJwtToken(parts[1], secret)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("user", claims)

		c.Next()
	}
}
