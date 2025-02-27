package middleware

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"time"
)

// 结构体用于捕获响应数据
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b) // 保存响应数据
	return w.ResponseWriter.Write(b)
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		var requestBody string
		var requestParams string

		// 解析 GET 请求参数
		if c.Request.Method == "GET" {
			requestParams = c.Request.URL.RawQuery
		} else if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			// 解析 `multipart/form-data`
			if c.ContentType() == "multipart/form-data" {
				err := c.Request.ParseMultipartForm(32 << 20) // 32MB 限制
				if err == nil {
					requestParams = fmt.Sprintf("%v", c.Request.Form)
				}
			} else {
				// 解析 `application/json` 或 `x-www-form-urlencoded`
				if err := c.Request.ParseForm(); err == nil {
					requestParams = c.Request.Form.Encode()
				}

				// 解析 `application/json`
				if c.Request.Header.Get("Content-Type") == "application/json" {
					bodyBytes, _ := io.ReadAll(c.Request.Body)
					c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // 重新赋值，防止数据丢失
					requestBody = string(bodyBytes)
				}
			}
		}

		// 记录请求信息
		zap.S().Infof("\n请求路径: %s\n 请求方法: %s\n 请求参数: %s\n 请求体: %s",
			c.Request.URL.Path, c.Request.Method, requestParams, requestBody)

		// 捕获响应
		respWriter := &responseWriter{ResponseWriter: c.Writer, body: bytes.NewBufferString("")}
		c.Writer = respWriter

		// 处理请求
		c.Next()

		// 计算耗时
		duration := time.Since(startTime)

		// 记录响应信息
		zap.S().Infof("\n请求路径: %s\n 状态码: %d\n\n 耗时: %v",
			c.Request.URL.Path, c.Writer.Status(), duration)
	}
}
