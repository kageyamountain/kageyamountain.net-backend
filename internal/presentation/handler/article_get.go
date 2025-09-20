package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/application/usecase"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/presentation/openapi"
)

type ArticleGetHandler struct {
	useCase usecase.ArticleGetUseCase
}

func NewArticleGetHandler(useCase usecase.ArticleGetUseCase) *ArticleGetHandler {
	return &ArticleGetHandler{
		useCase: useCase,
	}
}

// GET /articles エンドポイント
func (a *ArticleGetHandler) ArticleGet(c *gin.Context, articleId string) {
	ctx := c.Request.Context()

	// ユースケースの実行
	useCaseOutput, err := a.useCase.Execute(ctx, articleId)
	if err != nil {
		slog.Error("failed to ArticleGet use case", slog.Any("err", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, openapi.Error{
			Code:    openapi.InternalServerError,
			Message: "server error",
		})
		return
	}

	// 指定IDの記事が存在しない場合は404
	if useCaseOutput == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, openapi.Error{
			Code:    openapi.NotFound,
			Message: "article not found",
		})
		return
	}

	c.JSON(http.StatusOK, &openapi.Article{
		Id:          useCaseOutput.ID,
		UpdatedAt:   &useCaseOutput.UpdatedAt,
		PublishedAt: useCaseOutput.PublishedAt,
		Title:       useCaseOutput.Title,
		Contents:    &useCaseOutput.Contents,
		Tags:        useCaseOutput.Tags,
	})
}
