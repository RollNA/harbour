package middleware

import (
	"net/http"

	"github.com/RollNA/harbour/zLog"
	"github.com/gin-contrib/pprof"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

type Conf struct {
	ApplicationName string
	UsePprof        bool
}

func InitMiddleware(r *gin.Engine, conf Conf) {
	r.Use(Cors())
	//r.Use(gzip.Gzip(gzip.DefaultCompression))
	// Set X-Request-Id header
	r.Use(RequestId())
	if conf.UsePprof {
		pprof.Register(r)
	}
	p := NewPrometheus("")
	p.Use(r)
	r.Use(otelgin.Middleware(conf.ApplicationName))
	// 日志处理
	r.Use(LoggerToFile())
	r.Use(panicApi)
}

func panicApi(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			zLog.TraceError(
				c.Request.Context(),
				"HttpPanic",
				zap.Any("panic", r),
				zap.Any("url", c.Request.URL),
				zap.Stack("stack"),
			)
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{
					"code":    500,
					"message": "Internal Server Error",
					"success": false,
				},
			)
		}
	}()
	c.Next()
}

func authMiddleware(c *gin.Context) {
	if false {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{
				"code":    401,
				"message": "Invalid token",
				"success": false,
			},
		)
	}
	c.Next()
}
