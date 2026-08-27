package api

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"ricantares.com/ballroom/internal/logger"
	"ricantares.com/ballroom/internal/security"
)

// struttura del body logger
type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// RbacHandler is a middleware that checks if the request is authorized according to the configured RBAC rules.
// If the request is not authorized, it returns a 401 Unauthorized response and aborts the request.
func RbacHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !security.AccessGranted(c) {
			//c.JSON(http.StatusUnauthorized, ApiResponse{Code: http.StatusUnauthorized, Message: "", Data: ""})
			c.AbortWithStatusJSON(http.StatusUnauthorized, ApiResponse{Code: http.StatusUnauthorized, Message: "", Data: ""})
			return
		}

		c.Next()
	}
}

// imposta il logger
func LoggerWrapper(router *gin.Engine) {
	router.Use(requestLogger())
}

/*************  ✨ Windsurf Command ⭐  *************/
// Write implements the http.ResponseWriter interface.
// It writes the response body to the internal buffer and also to the original response writer.
// This allows to log the response body.
/*******  e8872012-5151-4eab-8321-b0c0d055eef5  *******/
func (w bodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

/*************  ✨ Windsurf Command ⭐  *************/
// requestLogger is a middleware that logs the request and response information, including the response body.
// If the response status is less than 400, it logs the information at the INFO level, otherwise it logs at the ERROR level.
/*******  309f4841-0c98-4615-b628-8e334c7c7914  *******/
func requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		bl := &bodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bl
		t := time.Now()
		c.Next()
		latency := time.Since(t)
		status := c.Writer.Status()
		headers := c.Request.Header
		message := fmt.Sprintf("%d %s %s %s %s %s headers: %v",
			status,
			c.Request.Method,
			c.Request.RequestURI,
			c.Request.Proto,
			c.Errors,
			latency,
			headers,
		)
		if status < 400 {
			logger.LogInfo(message + " body: " + bl.body.String())
		} else {
			logger.LogError(message + " body: " + bl.body.String())
		}
	}
}
