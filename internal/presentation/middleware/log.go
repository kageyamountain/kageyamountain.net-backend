package middleware

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kageyamountain/kageyamountain.net-backend/common/logger"
)

func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// ログMapの設定
		logMap := &sync.Map{}
		logMap.Store("request_id", uuid.New().String())
		logMap.Store("method", c.Request.Method)
		logMap.Store("path", c.Request.URL.Path)

		// contextにログMapをセット
		ctx := context.WithValue(c.Request.Context(), logger.ContextKeyLogMap, logMap)
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		// アクセスログを出力
		duration := time.Since(start)
		slog.InfoContext(ctx, "request started",
			slog.String("host", c.Request.Host),
			slog.String("uri", c.Request.URL.RequestURI()),
			slog.Int("status", c.Writer.Status()),
			slog.Int("response_size", c.Writer.Size()),
			slog.String("referer", c.Request.Referer()),
			slog.String("user_agent", c.Request.UserAgent()),
			slog.Int64("duration_ms", duration.Milliseconds()),
		)
	}
}
