package middleware

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/common/logger"
)

const (
	HttpHeaderXRequestID = "X-Request-ID"
)

func Log() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// RequestID生成
		requestID := uuid.New().String()
		// レスポンスヘッダーにRequestIDをセット
		c.Header(HttpHeaderXRequestID, requestID)

		// ログMapの設定
		logMap := &sync.Map{}
		logMap.Store("log_type", logger.LogTypeApp)
		logMap.Store("request_id", requestID)
		logMap.Store("method", c.Request.Method)
		logMap.Store("path", c.Request.URL.Path)

		// contextにログMapをセット
		ctx := context.WithValue(c.Request.Context(), logger.ContextKeyLogMap, logMap)
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		// アクセスログを出力
		logMap.Store("log_type", logger.LogTypeAccess)
		slog.InfoContext(ctx, "access log",
			slog.String("host", c.Request.Host),
			slog.String("uri", c.Request.URL.RequestURI()),
			slog.Int("status", c.Writer.Status()),
			slog.Int("response_size", c.Writer.Size()),
			slog.String("referer", c.Request.Referer()),
			slog.String("user_agent", c.Request.UserAgent()),
			slog.Int64("duration_ms", time.Since(start).Milliseconds()),
		)
	}
}
