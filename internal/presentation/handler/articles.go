package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/application/usecase"
)

type ArticlesHandler struct {
	useCase usecase.ArticlesUseCase
}

func NewArticles(useCase usecase.ArticlesUseCase) *ArticlesHandler {
	return &ArticlesHandler{
		useCase: useCase,
	}
}

func (a *ArticlesHandler) Execute(c *gin.Context) {
	ctx := c.Request.Context()

	articles, err := a.useCase.Execute(ctx)
	if err != nil {
		slog.Error("failed to articles use case", slog.Any("err", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	fmt.Println(articles)

	c.JSON(http.StatusOK, gin.H{"body": "OK"})
}
