package middleware

import (
	"net/http"

	"github.com/Oralkhan-coder/mind-map/internal/core"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			code := http.StatusInternalServerError
			msg := err.Error()

			if appErr, ok := err.Err.(*core.AppError); ok {
				code = appErr.Code
				msg = appErr.Message
			}

			c.JSON(code, gin.H{
				"error": msg,
			})
		}
	}
}
