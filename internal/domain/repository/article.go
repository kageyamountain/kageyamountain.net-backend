//go:generate sh -c "go tool moq -out ./$(basename ${GOFILE} .go)_mock.go . ArticleRepository"
package repository

import (
	"context"

	"github.com/kageyamountain/kageyamountain.net-backend/internal/domain/model/entity"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/domain/model/value"
)

type ArticleRepository interface {
	FindAllForList(ctx context.Context) ([]*entity.Article, error)
	FindByID(ctx context.Context, articleID *value.ArticleID) (*entity.Article, error)
}
