package router

import "github.com/kageyamountain/kageyamountain.net-backend/internal/presentation/handler"

type ServerInterfaceHandler struct {
	*handler.ArticlesGetHandler
	*handler.ArticleGetHandler
}

func NewServerInterfaceHandler(
	articlesGetHandler *handler.ArticlesGetHandler,
	articleGetHandler *handler.ArticleGetHandler,
) *ServerInterfaceHandler {
	return &ServerInterfaceHandler{
		ArticlesGetHandler: articlesGetHandler,
		ArticleGetHandler:  articleGetHandler,
	}
}
