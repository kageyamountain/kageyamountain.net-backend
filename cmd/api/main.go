package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/common/config"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/common/logger"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/presentation/router"
)

func main() {
	ctx := context.Background()

	// logger設定（カスタムslogハンドラを設定）
	handler := logger.NewAppLogHandler(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(slog.New(handler))

	// Ginのデフォルトログを無効化
	gin.SetMode(gin.ReleaseMode)

	// 環境変数のロード（dev環境のみ）
	err := loadEnvFileOnlyDev()
	if err != nil {
		slog.Error("failed to load .env file.", slog.Any("err", err))
		return
	}

	// 環境変数をAppConfigへマッピング
	appConfig, err := config.Load()
	if err != nil {
		slog.Error("failed to AppConfig Load.", slog.Any("err", err))
		return
	}

	// ルーティングの設定
	r, err := router.Setup(ctx, appConfig)
	if err != nil {
		slog.Error("failed to setup router.", slog.Any("err", err))
		return
	}

	// Webサーバーの起動
	err = r.Run("0.0.0.0:8080")
	if err != nil {
		slog.Error("failed to run server.", slog.Any("err", err))
		return
	}
}

// dev環境の場合は .env.dev から環境変数を読み込む
// dev環境以外の場合はコンテナ起動時に環境変数を設定するためenvファイルは読み込まない
func loadEnvFileOnlyDev() error {
	env := os.Getenv("ENV")
	if env == "dev" {
		err := godotenv.Load(".env.dev")
		if err != nil {
			return err
		}
	}
	return nil
}
