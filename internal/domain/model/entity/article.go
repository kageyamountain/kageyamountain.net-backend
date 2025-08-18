package entity

import (
	"errors"
	"time"

	"github.com/kageyamountain/kageyamountain.net-backend/internal/domain/model/value/enum"
)

type Article struct {
	ID            string
	Status        enum.Status
	CreatedAt     time.Time
	PublishedAt   time.Time
	PublishedYear string
	Title         string
	Contents      string
	Tags          []enum.Tag
}

type NewArticleInput struct {
	ID            string
	Status        string
	CreatedAt     time.Time
	PublishedAt   time.Time
	PublishedYear string
	Title         string
	Contents      string
	Tags          []string
}

func NewArticle(input *NewArticleInput) (*Article, error) {
	if input.ID == "" {
		return nil, errors.New("partition key is required")
	}

	status, err := enum.NewStatus(input.Status)
	if err != nil {
		return nil, err
	}

	var tags []enum.Tag
	for _, tag := range input.Tags {
		t, err := enum.NewTag(tag)
		if err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}

	return &Article{
		ID:            input.ID,
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
	return a.Status == enum.StatusDraft
}

func (a *Article) IsPublish() bool {
	return a.Status == enum.StatusPublish
}
