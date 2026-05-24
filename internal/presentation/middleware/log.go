package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/common/logger"
)

const HttpHeaderXRequestID = "X-Request-ID"

func Log() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// RequestID生成
		requestID := uuid.New().String()
		// レスポンスヘッダーにRequestIDをセット
		c.Header(HttpHeaderXRequestID, requestID)

		// LogContextの設定
		logContext := logger.NewLogContext()
		logContext.Set("log_type", logger.LogTypeApp)
		logContext.Set("request_id", requestID)
		logContext.Set("method", c.Request.Method)
		logContext.Set("path", c.Request.URL.Path)

		// contextにlogContextをセット
		ctx := logger.WithLogContext(c.Request.Context(), logContext)
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		// アクセスログを出力（log_typeはcall-siteで明示することでLogContextより優先される）
		slog.InfoContext(ctx, "access log",
			slog.Any("log_type", logger.LogTypeAccess),
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
