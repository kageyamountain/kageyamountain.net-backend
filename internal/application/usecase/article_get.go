package usecase

import (
	"context"
	"time"

	"github.com/kageyamountain/kageyamountain.net-backend/internal/domain/model/entity"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/domain/repository"
)

type ArticleGetUseCase interface {
	Execute(ctx context.Context, articleID string) (*ArticleGetUseCaseOutput, error)
}

type articleGetUseCase struct {
	articleRepository repository.ArticleRepository
}

func NewArticleUseCase(articleRepository repository.ArticleRepository) ArticleGetUseCase {
	return &articleGetUseCase{
		articleRepository: articleRepository,
	}
}

type ArticleGetUseCaseOutput struct {
	ID          string
	PublishedAt time.Time
	Title       string
	Contents    string
	Tags        []string
}

func (a *articleGetUseCase) Execute(ctx context.Context, articleID string) (*ArticleGetUseCaseOutput, error) {
	articleEntity, err := a.articleRepository.FindByID(ctx, articleID)
	if err != nil {
		return nil, err
	}

	// 指定IDの記事が存在しない場合はnilを返す
	if articleEntity == nil {
		return nil, nil
	}

	// 公開中の記事でない場合はnilを返す
	if !articleEntity.IsPublish() {
		return nil, nil
	}

	return a.convertToOutput(articleEntity), nil
}

func (a *articleGetUseCase) convertToOutput(article *entity.Article) *ArticleGetUseCaseOutput {
	output := &ArticleGetUseCaseOutput{
		ID:          article.ID,
		PublishedAt: article.PublishedAt,
		Title:       article.Title,
		Contents:    article.Contents,
		Tags:        make([]string, len(article.Tags)),
	}

	for i, tag := range article.Tags {
		output.Tags[i] = tag.String()
	}

	return output
}
