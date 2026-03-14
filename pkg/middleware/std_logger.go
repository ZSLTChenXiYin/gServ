package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

func StdoutLogger(logger *log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) <= 0 {
			logger.Printf("Status: %d, Method: %s, Path: %s, Query: %v", c.Writer.Status(), c.Request.Method, c.Request.URL.Path, c.Request.URL.RawQuery)
		}
	}
}

func StderrLogger(logger *log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			logger.Printf("Status: %d, Method: %s, Path: %s, Query: %v, Error: %s", c.Writer.Status(), c.Request.Method, c.Request.URL.Path, c.Request.URL.RawQuery, c.Errors.String())
		}
	}
}
