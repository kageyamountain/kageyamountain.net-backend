package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kageyamountain/kageyamountain.net-backend/common/config"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/presentation/router"
)

func main() {
	ctx := context.Background()

	err := loadEnvFileOnlyDev()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
		return
	}

	appConfig, err := config.Load()
	if err != nil {
		log.Fatal("Error AppConfig Load. err:", err)
		return
	}
	log.Println("appConfig:", appConfig)

	// TODO ENV毎にmode変更
	gin.SetMode(gin.DebugMode)

	r, err := router.Setup(ctx, appConfig)
	if err != nil {
		log.Fatal("Error router setup. err:", err)
		return
	}

	err = r.Run("localhost:8080")
	if err != nil {
		log.Fatal("Error gin.Run. err:", err)
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
