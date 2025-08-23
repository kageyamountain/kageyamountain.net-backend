package value

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewArticleID(t *testing.T) {
	t.Parallel()

	t.Run("正常系", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name  string
			input string
		}{
			{name: "ハイフン無しのUUID", input: "1234567812344abc8def1234567890ab"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				// Arrange

				// Act
				gotArticleID, err := NewArticleID(tt.input)

				// Assert
				a := assert.New(t)
				a.NoError(err)
				a.Equal(tt.input, gotArticleID.Value())
			})
		}
	})

	t.Run("異常系", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name       string
			input      string
			wantErrMsg string
		}{
			{name: "ハイフンあり", input: "12345678-1234-4abc-8def-1234567890ab", wantErrMsg: "article_id cannot contain hyphen. article_id: 12345678-1234-4abc-8def-1234567890ab"},
			{name: "不正な形式", input: "not-a-uuid", wantErrMsg: "invalid format of article_id. article_id: not-a-uuid"},
			{name: "空文字", input: "", wantErrMsg: "invalid format of article_id. article_id: "},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				// Arrange

				// Act
				gotArticleID, err := NewArticleID(tt.input)

				// Assert
				a := assert.New(t)
				a.Error(err)
				a.Contains(err.Error(), tt.wantErrMsg)
				a.Nil(gotArticleID)
			})
		}
	})
}

func TestGenerateArticleID(t *testing.T) {
	t.Parallel()

	t.Run("正常系", func(t *testing.T) {
		t.Parallel()
		t.Run("生成されたIDは32文字の小文字16進数でハイフン無し", func(t *testing.T) {
			t.Parallel()
			// Arrange

			// Act
			gotArticleID := GenerateArticleID()

			// Assert
			a := assert.New(t)
			a.NotNil(gotArticleID)
			a.Len(gotArticleID.Value(), 32)
			a.NotContains(gotArticleID.Value(), "-")
			a.Regexp(regexp.MustCompile(`^[0-9a-f]{32}$`), gotArticleID.Value())
		})

		t.Run("100回連続生成して重複しない", func(t *testing.T) {
			t.Parallel()
			// Arrange
			ids := make(map[string]int)
			exists := false

			// Act
			for i := 0; i < 100; i++ {
				id := GenerateArticleID().Value()
				_, exists = ids[id]
				if exists {
					break
				}
				ids[id] = i
			}

			// Assert
			a := assert.New(t)
			a.False(exists)
			a.Len(ids, 100)
		})
	})
}
