package main

import (
	"context"
	"log/slog"

	"github.com/joho/godotenv"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/common/config"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/common/logger"
)

func main() {
	ctx := context.Background()

	// dev環境変数のロード
	err := godotenv.Load(".env.dev")
	if err != nil {
		slog.ErrorContext(ctx, "failed to load .env file.", slog.Any(logger.AttrKeyError, err))
		return
	}

	// 環境変数をAppConfigへマッピング
	appConfig, err := config.Load()
	if err != nil {
		slog.ErrorContext(ctx, "failed to AppConfig Load.", slog.Any(logger.AttrKeyError, err))
		return
	}

	// 出力
	slog.InfoContext(ctx, "AppConfig loaded.", slog.Any("app_config", appConfig))
}
