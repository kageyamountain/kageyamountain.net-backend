package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/application/usecase"
	openapi "github.com/kageyamountain/kageyamountain.net-backend/internal/presentation/openapi/generate"
)

type ArticlesGetHandler struct {
	useCase usecase.ArticlesGetUseCase
}

func NewArticlesGetHandler(useCase usecase.ArticlesGetUseCase) *ArticlesGetHandler {
	return &ArticlesGetHandler{
		useCase: useCase,
	}
}

// GET /articles エンドポイント
func (a *ArticlesGetHandler) ArticlesGet(c *gin.Context, params openapi.ArticlesGetParams) {
	ctx := c.Request.Context()

	useCaseOutput, err := a.useCase.Execute(ctx)
	if err != nil {
		slog.Error("failed to useCaseOutput use case", slog.Any("err", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	fmt.Println(useCaseOutput)

	var articles []openapi.Article
	for _, article := range useCaseOutput.Articles {
		articles = append(articles, openapi.Article{
			Id:          article.PK,
			PublishedAt: article.PublishedAt,
			Title:       article.Title,
			Tags:        article.Tags,
		})
	}

	c.JSON(http.StatusOK, openapi.ArticlesGetResponseBody{
		Articles: articles,
	})
}
