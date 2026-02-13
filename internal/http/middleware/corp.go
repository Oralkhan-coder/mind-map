package middleware

import "github.com/gin-gonic/gin"

func CORP() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set(
			"Cross-Origin-Opener-Policy",
			"same-origin-allow-popups",
		)
		c.Next()
	}
}
