package value

import (
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
				got, err := NewArticleID(tt.input)

				// Assert
				a := assert.New(t)
				a.NoError(err)
				a.Equal(tt.input, got.Value())
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
				got, err := NewArticleID(tt.input)

				// Assert
				a := assert.New(t)
				a.Error(err)
				a.Contains(err.Error(), tt.wantErrMsg)
				a.Nil(got)
			})
		}
	})
}
