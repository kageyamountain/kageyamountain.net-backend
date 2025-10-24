package router

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/common/config"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/infrastructure/gateway/dynamodb"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/infrastructure/repository"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/module/article_get"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/module/articles_get"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/presentation/middleware"
	openapi "github.com/kageyamountain/kageyamountain.net-backend/internal/presentation/openapi/v1"
)

func Setup(ctx context.Context, appConfig *config.AppConfig) (*gin.Engine, error) {
	siw, err := initializeHandler(ctx, appConfig)
	if err != nil {
		return nil, err
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CORS(appConfig))
	r.Use(middleware.Log())

	v1 := r.Group("/v1")
	v1.GET("/articles", siw.ArticlesGet)            // 記事一覧
	v1.GET("/articles/:article-id", siw.ArticleGet) // 記事詳細

	return r, nil
}

func initializeHandler(ctx context.Context, appConfig *config.AppConfig) (*openapi.ServerInterfaceWrapper, error) {
	// Gateway
	dynamoDB, err := dynamodb.NewClient(ctx, appConfig)
	if err != nil {
		return nil, err
	}

	// Repository
	articleRepository := repository.NewArticleRepository(dynamoDB, appConfig)

	// module
	articlesUseCase := articles_get.NewUseCase(articleRepository)
	articlesGetHandler := articles_get.NewHandler(articlesUseCase)

	articleUseCase := article_get.NewUseCase(articleRepository)
	articleGetHandler := article_get.NewHandler(articleUseCase)

	// OpenAPI生成コードのServerInterfaceを実装するHandlerを作成
	sih := NewServerInterfaceHandler(articlesGetHandler, articleGetHandler)

	// OpenAPI生成コードのHandlerのラッパーを作成
	// middlewareはルーティング時にパス毎に個別に設定する
	siw := &openapi.ServerInterfaceWrapper{
		Handler: sih,
		ErrorHandler: func(c *gin.Context, err error, statusCode int) {
			c.AbortWithStatusJSON(statusCode, &openapi.Error{
				Code:    openapi.InternalServerError,
				Message: err.Error(),
			})
		},
	}

	return siw, nil
}
