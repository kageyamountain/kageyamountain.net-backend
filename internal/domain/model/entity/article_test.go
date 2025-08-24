package entity

import (
	"testing"
	"time"

	"github.com/kageyamountain/kageyamountain.net-backend/internal/domain/model/value"
	"github.com/kageyamountain/kageyamountain.net-backend/internal/domain/model/value/enum"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewArticle(t *testing.T) {
	t.Parallel()

	t.Run("正常系", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name       string
			input      *NewArticleInput
			wantStatus enum.Status
			wantTags   []enum.Tag
		}{
			{
				name: "全項目が正常値な場合はArticleが生成される",
				input: &NewArticleInput{
					ID:            value.GenerateArticleID().Value(),
					Status:        enum.StatusPublish.String(),
					CreatedAt:     time.Date(2024, 12, 14, 10, 14, 14, 0, time.UTC),
					PublishedAt:   time.Date(2024, 12, 15, 11, 15, 15, 0, time.UTC),
					PublishedYear: "2024",
					Title:         "タイトル",
					Contents:      "本文",
					Tags:          []string{enum.TagGo.String(), enum.TagAWS.String(), enum.TagGin.String()},
				},
				wantStatus: enum.StatusPublish,
				wantTags:   []enum.Tag{enum.TagGo, enum.TagAWS, enum.TagGin},
			},
			{
				name: "タグが空スライスでもエラーにならない",
				input: &NewArticleInput{
					ID:            value.GenerateArticleID().Value(),
					Status:        enum.StatusDraft.String(),
					CreatedAt:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					PublishedAt:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					PublishedYear: "2024",
					Title:         "タイトル",
					Contents:      "本文",
					Tags:          []string{},
				},
				wantStatus: enum.StatusDraft,
				wantTags:   nil,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				// Arrange

				// Act
				gotArticle, err := NewArticle(tt.input)

				// Assert
				r := require.New(t)
				a := assert.New(t)
				a.NoError(err)
				a.NotNil(gotArticle)

				wantID, err := value.NewArticleID(tt.input.ID)
				r.NoError(err)
				want := &Article{
					ID:            *wantID,
					Status:        tt.wantStatus,
					CreatedAt:     tt.input.CreatedAt,
					PublishedAt:   tt.input.PublishedAt,
					PublishedYear: tt.input.PublishedYear,
					Title:         tt.input.Title,
					Contents:      tt.input.Contents,
					Tags:          tt.wantTags,
				}
				a.Equal(want, gotArticle)
			})
		}
	})

	t.Run("異常系", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name       string
			input      *NewArticleInput
			wantErrMsg string
		}{
			{
				name: "IDが空文字",
				input: &NewArticleInput{
					ID:            "",
					Status:        enum.StatusPublish.String(),
					CreatedAt:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					PublishedAt:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					PublishedYear: "2024",
					Title:         "タイトル",
					Contents:      "本文",
					Tags:          []string{enum.TagGo.String()},
				},
				wantErrMsg: "id is required",
			},
			{
				name: "IDにハイフンが含まれている",
				input: &NewArticleInput{
					ID:            "12345678-1234-4abc-8def-1234567890ab",
					Status:        enum.StatusPublish.String(),
					CreatedAt:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					PublishedAt:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					PublishedYear: "2024",
					Title:         "タイトル",
					Contents:      "本文",
					Tags:          []string{enum.TagGo.String()},
				},
				wantErrMsg: "article_id cannot contain hyphen. article_id:",
			},
			{
				name: "IDがUUID形式でない",
				input: &NewArticleInput{
					ID:            "not-a-uuid",
					Status:        enum.StatusPublish.String(),
					CreatedAt:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					PublishedAt:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					PublishedYear: "2024",
					Title:         "タイトル",
					Contents:      "本文",
					Tags:          []string{enum.TagGo.String()},
				},
				wantErrMsg: "invalid format of article_id. article_id:",
			},
			{
				name: "不正なステータス",
				input: &NewArticleInput{
					ID:            value.GenerateArticleID().Value(),
					Status:        "invalid status",
					CreatedAt:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					PublishedAt:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					PublishedYear: "2024",
					Title:         "タイトル",
					Contents:      "本文",
					Tags:          []string{enum.TagGo.String()},
				},
				wantErrMsg: "invalid status:",
			},
			{
				name: "不正なタグが含まれている",
				input: &NewArticleInput{
					ID:            value.GenerateArticleID().Value(),
					Status:        enum.StatusPublish.String(),
					CreatedAt:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					PublishedAt:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					PublishedYear: "2024",
					Title:         "タイトル",
					Contents:      "本文",
					Tags:          []string{enum.TagGo.String(), "UnknownTag", enum.TagAWS.String()},
				},
				wantErrMsg: "invalid tag:",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				// Arrange

				// Act
				gotArticle, err := NewArticle(tt.input)

				// Assert
				a := assert.New(t)
				a.Error(err)
				a.Contains(err.Error(), tt.wantErrMsg)
				a.Nil(gotArticle)
			})
		}
	})
}
