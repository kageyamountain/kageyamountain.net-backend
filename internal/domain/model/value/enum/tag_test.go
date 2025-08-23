package enum

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestParseTag(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Parallel()

	t.Run("正常系", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name    string
			input   string
			wantTag Tag
		}{
			{name: "Go", input: "Go", wantTag: TagGo},
			{name: "Gin", input: "Gin", wantTag: TagGin},
			{name: "JavaScript", input: "JavaScript", wantTag: TagJavaScript},
			{name: "TypeScript", input: "TypeScript", wantTag: TagTypeScript},
			{name: "Vue.js", input: "Vue.js", wantTag: TagVueJS},
			{name: "AWS", input: "AWS", wantTag: TagAWS},
			{name: "DynamoDB", input: "DynamoDB", wantTag: TagDynamoDB},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				// Arrange

				// Act
				gotTag, err := ParseTag(tt.input)

				// Assert
				a := assert.New(t)
				a.NoError(err)
				a.Equal(tt.wantTag, gotTag)
			})
		}
	})

	t.Run("異常系", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name    string
			input   string
			wantTag Tag
		}{
			{name: "空文字", input: "", wantTag: ""},
			{name: "想定されていない値", input: "UnknownTag", wantTag: ""},
			{name: "小文字・大文字の違い", input: "go", wantTag: ""},
			{name: "ドット無し", input: "Vuejs", wantTag: ""},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				// Arrange

				// Act
				gotTag, err := ParseTag(tt.input)

				// Assert
				a := assert.New(t)
				a.Error(err)
				a.Equal(err.Error(), "invalid tag: "+tt.input)
				a.Equal(tt.wantTag, gotTag)
			})
		}
	})
}
