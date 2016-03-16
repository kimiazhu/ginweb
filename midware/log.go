// Description: conf log 实现日志相关的配置项
// Author: ZHU HAIHUA
// Since: 2016-02-27 15:10
package conf

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/kimiazhu/log4go"
	"time"
)

var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow  = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset   = string([]byte{27, 91, 48, 109})
)

func ErrorLoggerT(typ gin.ErrorType) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		// avoid writting if we already wrote into the response body
		if !c.Writer.Written() {
			errors := c.Errors.ByType(typ)
			if len(errors) > 0 {
				c.JSON(-1, errors)
			}
		}
	}
}

func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path

		// Process request
		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		statusColor := colorForStatus(statusCode)
		methodColor := colorForMethod(method)
		comment := c.Errors.ByType(gin.ErrorTypePrivate).String()

		log.Access(fmt.Sprintf("%v | %s  %s %-7s %s | %s %3d %s | %13v | %s | %s",
			start.Format("2006/01/02 15:04:05.000000000"),
			methodColor, reset, method, path,
			statusColor, statusCode, reset,
			latency,
			clientIP,
			comment,
		))
	}
}

func colorForStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return green
	case code >= 300 && code < 400:
		return white
	case code >= 400 && code < 500:
		return yellow
	default:
		return red
	}
}

func colorForMethod(method string) string {
	switch method {
	case "GET":
		return blue
	case "POST":
		return cyan
	case "PUT":
		return yellow
	case "DELETE":
		return red
	case "PATCH":
		return green
	case "HEAD":
		return magenta
	case "OPTIONS":
		return white
	default:
		return reset
	}
}
