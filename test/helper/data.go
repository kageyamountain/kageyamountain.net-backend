package helper

import (
	"testing"
	"time"

	"github.com/kageyamountain/kageyamountain.net-backend/internal/domain/model/entity"
	"github.com/stretchr/testify/require"
)

func NewTestArticle(
	t *testing.T,
	id string,
	status string,
	createdAt time.Time,
	publishedAt time.Time,
	publishedYear string,
	title string,
	contents string,
	tags []string,
) *entity.Article {
	t.Helper()

	input := &entity.NewArticleInput{
		ID:            id,
		Status:        status,
		CreatedAt:     createdAt,
		PublishedAt:   publishedAt,
		PublishedYear: publishedYear,
		Title:         title,
		Contents:      contents,
		Tags:          tags,
	}

	article, err := entity.NewArticle(input)
	require.NoError(t, err)

	return article
}
