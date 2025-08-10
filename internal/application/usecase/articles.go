package usecase

import (
	"context"

	"github.com/kageyamountain/kageyamountain.net-backend/internal/domain/model/entity"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/domain/repository"
)

type ArticlesUseCase interface {
	Execute(ctx context.Context) (*ArticlesUseCaseOutput, error)
}

type articles struct {
	articleRepository repository.ArticleRepository
}

func NewArticles(articleRepository repository.ArticleRepository) ArticlesUseCase {
	return &articles{
		articleRepository: articleRepository,
	}
}

type ArticlesUseCaseOutput struct {
	Articles []*ArticlesUseCaseOutputRow
}

type ArticlesUseCaseOutputRow struct {
	PK          string
	PublishedAt string
	Title       string
	Tags        []string
}

func (a *articles) Execute(ctx context.Context) (*ArticlesUseCaseOutput, error) {
	articlesEntity, err := a.articleRepository.FindAllForList(ctx)
	if err != nil {
		return nil, err
	}

	return a.convertToOutput(articlesEntity), nil
}

func (a *articles) convertToOutput(articles []*entity.Article) *ArticlesUseCaseOutput {
	var output ArticlesUseCaseOutput
	for _, article := range articles {
		outputRow := &ArticlesUseCaseOutputRow{
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
