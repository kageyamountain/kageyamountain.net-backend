package router

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/kageyamountain/kageyamountain.net-backend/common/config"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/application/usecase"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/infrastructure/gateway"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/infrastructure/repository"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/presentation/handler"
)

func Setup(ctx context.Context, appConfig *config.AppConfig) (*gin.Engine, error) {
	r := gin.Default()

	articles, err := initializeHandler(ctx, appConfig)
	if err != nil {
		return nil, err
	}

	r.GET("/articles", articles.Execute)                    // 一覧
	r.GET("/articles/:article_id", func(c *gin.Context) {}) // 詳細

	return r, nil
}

func initializeHandler(ctx context.Context, appConfig *config.AppConfig) (*handler.Articles, error) {
	// gateway
	dynamoDB, err := gateway.NewDynamoDB(ctx, appConfig)
	if err != nil {
		return nil, err
	}

	// repository
	articleRepository := repository.NewArticleRepository(dynamoDB)

	// UseCase
	articlesUsecase := usecase.NewArticles(articleRepository)

	// Handler
	articlesHandler := handler.NewArticles(articlesUsecase)

	return articlesHandler, nil
}
