package router

import (
	"github.com/kageyamountain/kageyamountain.net-backend/internal/module/article_get"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/module/articles_get"
)

type ServerInterfaceHandler struct {
	*articles_get.ArticlesGetHandler
	*article_get.ArticleGetHandler
}

func NewServerInterfaceHandler(
	articlesGetHandler *articles_get.ArticlesGetHandler,
	articleGetHandler *article_get.ArticleGetHandler,
) *ServerInterfaceHandler {
	return &ServerInterfaceHandler{
		ArticlesGetHandler: articlesGetHandler,
		ArticleGetHandler:  articleGetHandler,
	}
}
