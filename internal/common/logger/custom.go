package logger

import (
	"context"
	"log/slog"
	"sync"
)

type contextKey string

const ContextKeyLogMap contextKey = "LogMap"

type LogType string

const (
	LogTypeApp    LogType = "app_log"
	LogTypeAccess LogType = "access_log"
)

// カスタムslog Handler: contextにセットされたログMap(sync.Map)から自動的にフィールドを追加する機能を持つ
type AppLogHandler struct {
	slog.Handler
}

func NewAppLogHandler(handler slog.Handler) *AppLogHandler {
	return &AppLogHandler{
		Handler: handler,
	}
}

func (h *AppLogHandler) Handle(ctx context.Context, r slog.Record) error { //nolint:gocritic //slogのinterface仕様なので第2引数はポインタ型にできない
	// contextからログMapを取得
	logMap, ok := ctx.Value(ContextKeyLogMap).(*sync.Map)
	if ok {
		// ログMapの全エントリをログに追加
		logMap.Range(func(key, value interface{}) bool {
			keyStr, ok2 := key.(string)
			if ok2 {
				r.AddAttrs(slog.Attr{Key: keyStr, Value: slog.AnyValue(value)})
			}
			return true // 継続
		})
	}

	// info: OpenTelemetryと連携する場合は有効化
	//// OpenTelemetry trace情報を追加
	// if span := trace.SpanFromContext(ctx); span.SpanContext().IsValid() {
	//	spanCtx := span.SpanContext()
	//	r.AddAttrs(
	//		slog.String("trace_id", spanCtx.TraceID().String()),
	//		slog.String("span_id", spanCtx.SpanID().String()),
	//	)
	//}

	return h.Handler.Handle(ctx, r)
}
