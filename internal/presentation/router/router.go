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
	openapi "github.com/kageyamountain/kageyamountain.net-backend/internal/presentation/openapi/generate"
)

func Setup(ctx context.Context, appConfig *config.AppConfig) (*gin.Engine, error) {
	siw, err := initializeHandler(ctx, appConfig)
	if err != nil {
		return nil, err
	}

	r := gin.Default()
	r.Use(middleware.Logging())

	r.GET("/articles", siw.ArticlesGet)                     // 記事一覧
	r.GET("/articles/:article_id", func(c *gin.Context) {}) // 記事詳細

	return r, nil
}

func initializeHandler(ctx context.Context, appConfig *config.AppConfig) (*openapi.ServerInterfaceWrapper, error) {
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
	articlesGetHandler := handler.NewArticlesGetHandler(articlesUsecase)

	// OpenAPI生成コードのHandlerのラッパーを作成
	// middlewareはルーティング時にパス毎に個別に設定する
	siw := &openapi.ServerInterfaceWrapper{
		Handler: articlesGetHandler,
		ErrorHandler: func(c *gin.Context, err error, statusCode int) {
			c.AbortWithStatusJSON(statusCode, &openapi.Error{
				Code:    openapi.InternalServerError,
				Message: err.Error(),
			})
		},
	}

	return siw, nil
}
