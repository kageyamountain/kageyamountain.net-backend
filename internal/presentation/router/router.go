package router

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/kageyamountain/kageyamountain.net-backend/common/config"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/application/usecase"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/infrastructure/gateway"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/infrastructure/repository"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/presentation/handler"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/presentation/middleware"
)

func Setup(ctx context.Context, appConfig *config.AppConfig) (*gin.Engine, error) {
	articles, err := initializeHandler(ctx, appConfig)
	if err != nil {
		return nil, err
	}

	r := gin.Default()
	r.Use(middleware.Logging())

	r.GET("/articles", articles.Execute)                    // 一覧
	r.GET("/articles/:article_id", func(c *gin.Context) {}) // 詳細

	return r, nil
}

func initializeHandler(ctx context.Context, appConfig *config.AppConfig) (*handler.ArticlesGetHandler, error) {
	// gateway
	dynamoDB, err := gateway.NewDynamoDB(ctx, appConfig)
	if err != nil {
		return nil, err
	}

	// repository
	articleRepository := repository.NewArticleRepository(dynamoDB)

	// UseCase
	articlesUsecase := usecase.NewArticlesUseCase(articleRepository)

	// Handler
	articlesHandler := handler.NewArticlesGetHandler(articlesUsecase)

	return articlesHandler, nil
}
