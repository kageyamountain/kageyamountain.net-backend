package entity

import (
	"errors"

	"github.com/kageyamountain/kageyamountain.net-backend/internal/domain/model/value"
)

type Article struct {
	PK            string
	Status        value.Status
	CreatedAt     string
	PublishedAt   string
	PublishedYear string
	Title         string
	Contents      string
	Tags          []value.Tag
}

type NewArticleInput struct {
	ID            string
	Status        string
	CreatedAt     string
	PublishedAt   string
	PublishedYear string
	Title         string
	Contents      string
	Tags          []string
}

func NewArticle(input *NewArticleInput) (*Article, error) {
	if input.ID == "" {
		return nil, errors.New("partition key is required")
	}

	status, err := value.NewStatus(input.Status)
	if err != nil {
		return nil, err
	}

	var tags []value.Tag
	for _, tag := range input.Tags {
		t, err := value.NewTag(tag)
		if err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}

	return &Article{
		PK:            input.ID,
		Status:        status,
		CreatedAt:     input.CreatedAt,
		PublishedAt:   input.PublishedAt,
		PublishedYear: input.PublishedYear,
		Title:         input.Title,
		Contents:      input.Contents,
		Tags:          tags,
	}, nil
}

func (a *Article) IsDraft() bool {
	return a.Status.Value() == value.StatusDraft
}

func (a *Article) IsPublish() bool {
	return a.Status.Value() == value.StatusPublish
}
