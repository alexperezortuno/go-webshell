package logging

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		// Request URL path
		path := c.Request.URL.Path

		if c.Request.URL.RawQuery != "" {
			path = path + "?" + c.Request.URL.RawQuery
		}

		// Process request
		c.Next()

		// Results
		timestamp := time.Now()
		latency := timestamp.Sub(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		userAgent := c.Request.UserAgent()

		fmt.Printf("%v | %3d | %s | %15s | %-7s %#v | %13v\n",
			timestamp.Format("2006/01/02 - 15:04:05"),
			statusCode,
			userAgent,
			clientIP,
			method,
			path,
			latency,
		)
	}
}
