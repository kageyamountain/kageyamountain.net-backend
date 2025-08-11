package usecase

import (
	"context"
	"time"

	"github.com/kageyamountain/kageyamountain.net-backend/internal/domain/model/entity"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/domain/repository"
)

type ArticlesGetUseCase interface {
	Execute(ctx context.Context) (*ArticlesGetUseCaseOutput, error)
}

type articlesGetUseCase struct {
	articleRepository repository.ArticleRepository
}

func NewArticlesUseCase(articleRepository repository.ArticleRepository) ArticlesGetUseCase {
	return &articlesGetUseCase{
		articleRepository: articleRepository,
	}
}

type ArticlesGetUseCaseOutput struct {
	Articles []*ArticlesGetUseCaseOutputRow
}

type ArticlesGetUseCaseOutputRow struct {
	PK          string
	PublishedAt time.Time
	Title       string
	Tags        []string
}

func (a *articlesGetUseCase) Execute(ctx context.Context) (*ArticlesGetUseCaseOutput, error) {
	articlesEntity, err := a.articleRepository.FindAllForList(ctx)
	if err != nil {
		return nil, err
	}

	return a.convertToOutput(articlesEntity), nil
}

func (a *articlesGetUseCase) convertToOutput(articles []*entity.Article) *ArticlesGetUseCaseOutput {
	var output ArticlesGetUseCaseOutput
	for _, article := range articles {
		outputRow := &ArticlesGetUseCaseOutputRow{
			PK:          article.PK,
			PublishedAt: article.PublishedAt,
			Title:       article.Title,
			Tags:        make([]string, len(article.Tags)),
		}
		for i, tag := range article.Tags {
			outputRow.Tags[i] = tag.Value()
		}
		output.Articles = append(output.Articles, outputRow)
	}

	return &output
}
