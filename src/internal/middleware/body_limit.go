package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func BodyLimit(limit int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, limit)
		c.Next()
	}
}
