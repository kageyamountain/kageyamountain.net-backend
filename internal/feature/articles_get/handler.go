package articles_get

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi "github.com/kageyamountain/kageyamountain.net-backend/internal/presentation/openapi/v1"
)

type ArticlesGetHandler struct {
	useCase UseCase
}

func NewHandler(useCase UseCase) *ArticlesGetHandler {
	return &ArticlesGetHandler{
		useCase: useCase,
	}
}

// GET /articles エンドポイント
func (a *ArticlesGetHandler) ArticlesGet(c *gin.Context, params openapi.ArticlesGetParams) {
	ctx := c.Request.Context()

	// TODO リクエストパラメータを利用した公開年とタグのフィルタリング機能実装

	useCaseOutput, err := a.useCase.Execute(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to ArticlesGet use case", slog.Any("err", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, openapi.Error{
			Code:    openapi.InternalServerError,
			Message: "server error",
		})
		return
	}

	// レスポンスボディの型へ変換
	articles := make([]openapi.ArticleSummary, 0, len(useCaseOutput.Articles))
	for _, article := range useCaseOutput.Articles {
		articles = append(articles, openapi.ArticleSummary{
			Id:          article.ID,
			PublishedAt: article.PublishedAt,
			Title:       article.Title,
			Tags:        article.Tags,
		})
	}

	c.JSON(http.StatusOK, openapi.ArticlesGetResponseBody{
		Articles: articles,
	})
}
