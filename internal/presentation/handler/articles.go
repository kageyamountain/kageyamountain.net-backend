package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/application/usecase"
)

type Articles struct {
	useCase usecase.ArticlesUseCase
}

func NewArticles(useCase usecase.ArticlesUseCase) *Articles {
	return &Articles{
		useCase: useCase,
	}
}

func (a *Articles) Execute(c *gin.Context) {
	ctx := c.Request.Context()

	articles, err := a.useCase.Execute(ctx)
	if err != nil {
		log.Println("Error executing Articles use case:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	fmt.Println(articles)

	c.JSON(http.StatusOK, gin.H{"body": "OK"})
}
