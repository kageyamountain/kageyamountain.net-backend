package repository

import (
	"context"

	"github.com/kageyamountain/kageyamountain.net-backend/internal/domain/model/entity"
)

type ArticleRepository interface {
	FindAllForList(ctx context.Context) ([]entity.Article, error)
}
