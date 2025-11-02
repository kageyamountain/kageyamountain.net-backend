//go:generate sh -c "go tool moq -out ./$(basename ${GOFILE} .go)_mock.go . UseCase"
package articles_get

import (
	"context"
	"time"

	"github.com/kageyamountain/kageyamountain.net-backend/internal/domain/model/entity"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/domain/repository"
)

type UseCase interface {
	Execute(ctx context.Context) (*UseCaseOutput, error)
}

type useCase struct {
	articleRepository repository.ArticleRepository
}

func NewUseCase(articleRepository repository.ArticleRepository) UseCase {
	return &useCase{
		articleRepository: articleRepository,
	}
}

type UseCaseOutput struct {
	Articles []*UseCaseOutputRow
}

type UseCaseOutputRow struct {
	ID          string
	PublishedAt time.Time
	Title       string
	Tags        []string
}

func (u *useCase) Execute(ctx context.Context) (*UseCaseOutput, error) {
	articlesEntity, err := u.articleRepository.FindAllForList(ctx)
	if err != nil {
		return nil, err
	}

	return u.convertToOutput(articlesEntity), nil
}

func (u *useCase) convertToOutput(articles []*entity.Article) *UseCaseOutput {
	var output UseCaseOutput
	for _, article := range articles {
		outputRow := &UseCaseOutputRow{
			ID:          article.ID.Value(),
			PublishedAt: article.PublishedAt,
			Title:       article.Title,
			Tags:        make([]string, len(article.Tags)),
		}
		for i, tag := range article.Tags {
			outputRow.Tags[i] = tag.String()
		}
		output.Articles = append(output.Articles, outputRow)
	}

	return &output
}
