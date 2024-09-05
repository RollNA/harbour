package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/quincy0/harbour/zLog"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// 打印response，refer：https://stackoverflow.com/questions/38501325/how-to-log-response-body-in-gin
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

// LoggerToFile 日志记录到文件
func LoggerToFile() gin.HandlerFunc {

	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		if c.Request.Method == http.MethodOptions {
			//options 不打印日志
			c.Next()
			return
		}

		// payload
		payload := ""
		if c.Request.Method == http.MethodPost {
			payloadBytes, _ := c.GetRawData()
			c.Request.Body = io.NopCloser(bytes.NewBuffer(payloadBytes))
			var jsonData interface{}
			_ = json.Unmarshal(payloadBytes, &jsonData)
			marshalledBytes, _ := json.Marshal(jsonData)
			payload = string(marshalledBytes)
		}

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime)
		// 请求方式
		reqMethod := c.Request.Method
		// 请求路由
		reqUri := c.Request.RequestURI
		// 状态码
		statusCode := c.Writer.Status()
		// 请求IP
		clientIP := c.ClientIP()
		// ua
		ua := c.Request.UserAgent()
		// body string

		bodyString := blw.body.String()
		encoding := blw.Header().Get("Content-Encoding")
		if encoding == "gzip" {
			bodyString = "gziped body"
		}
		if len(bodyString) > 1024 {
			bodyString = bodyString[0:1024]
		}
		// 日志格式
		zLog.TraceInfo(
			c.Request.Context(),
			"gin request",
			zap.String("start", startTime.Format("2006-01-02 15:04:05.9999")),
			zap.Any("statusCode", statusCode),
			zap.Any("cost", latencyTime),
			zap.String("clientIP", clientIP),
			zap.String("method", reqMethod),
			zap.String("uri", reqUri),
			zap.String("ua", ua),
			zap.String("payload", payload),
			zap.String("response", bodyString),
		)
	}
}
