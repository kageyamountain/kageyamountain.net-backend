//go:generate sh -c "go tool moq -out ./$(basename ${GOFILE} .go)_mock.go . UseCase"
package article_get

import (
	"context"
	"fmt"
	"time"

	"github.com/kageyamountain/kageyamountain.net-backend/internal/common/apperror"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/domain/model/entity"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/domain/model/value"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/domain/repository"
)

type UseCase interface {
	Execute(ctx context.Context, articleID string) (*UseCaseOutput, error)
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
	ID          string
	UpdatedAt   time.Time
	PublishedAt time.Time
	Title       string
	Contents    string
	Tags        []string
}

func (u *useCase) Execute(ctx context.Context, inputArticleID string) (*UseCaseOutput, error) {
	articleID, err := value.NewArticleID(inputArticleID)
	if err != nil {
		return nil, err
	}

	articleEntity, err := u.articleRepository.FindByID(ctx, articleID)
	if err != nil {
		return nil, err
	}

	// 非公開記事の場合はエラーを返す
	if !articleEntity.IsPublish() {
		return nil, apperror.NewErrorUnpublishedArticle(fmt.Sprintf("article_id: %s", articleID.Value()))
	}

	return u.convertToOutput(articleEntity), nil
}

func (u *useCase) convertToOutput(article *entity.Article) *UseCaseOutput {
	output := &UseCaseOutput{
		ID:          article.ID.Value(),
		UpdatedAt:   article.UpdatedAt,
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
