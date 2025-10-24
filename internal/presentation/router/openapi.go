package router

import (
	"github.com/kageyamountain/kageyamountain.net-backend/internal/module/articles_get"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/presentation/handler"
)

type ServerInterfaceHandler struct {
	*articles_get.Handler
	*handler.ArticleGetHandler
}

func NewServerInterfaceHandler(
	articlesGetHandler *articles_get.Handler,
	articleGetHandler *handler.ArticleGetHandler,
) *ServerInterfaceHandler {
	return &ServerInterfaceHandler{
		Handler:           articlesGetHandler,
		ArticleGetHandler: articleGetHandler,
	}
}
