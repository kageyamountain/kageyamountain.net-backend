package article_get

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/common/apperror"
	openapi "github.com/kageyamountain/kageyamountain.net-backend/internal/presentation/openapi/v1"
)

type ArticleGetHandler struct {
	useCase UseCase
}

func NewHandler(useCase UseCase) *ArticleGetHandler {
	return &ArticleGetHandler{
		useCase: useCase,
	}
}

// GET /articles エンドポイント
func (a *ArticleGetHandler) ArticleGet(c *gin.Context, articleId string) {
	ctx := c.Request.Context()

	useCaseOutput, err := a.useCase.Execute(ctx, articleId)
	if err != nil {
		if apperror.IsOneOf(err, apperror.ErrorNoData, apperror.ErrorUnpublishedArticle) {
			slog.InfoContext(ctx, "failed to use case", slog.Any("err", err))
			c.AbortWithStatusJSON(http.StatusNotFound, openapi.Error{
				Code:    openapi.NotFound,
				Message: "article not found",
			})
			return
		}

		slog.ErrorContext(ctx, "failed to use case", slog.Any("err", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, openapi.Error{
			Code:    openapi.InternalServerError,
			Message: "server error",
		})
		return
	}

	c.JSON(http.StatusOK, &openapi.ArticleGetResponseBody{
		Article: openapi.ArticleDetail{
			Id:          useCaseOutput.ID,
			UpdatedAt:   useCaseOutput.UpdatedAt,
			PublishedAt: useCaseOutput.PublishedAt,
			Title:       useCaseOutput.Title,
			Contents:    useCaseOutput.Contents,
			Tags:        useCaseOutput.Tags,
		},
	})
}
