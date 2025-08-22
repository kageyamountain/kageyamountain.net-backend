package enum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseStatus(t *testing.T) {
	t.Parallel()

	t.Run("正常系", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name       string
			input      string
			wantStatus Status
		}{
			{name: "draft", input: "draft", wantStatus: StatusDraft},
			{name: "publish", input: "publish", wantStatus: StatusPublish},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				// Arrange

				// Act
				gotStatus, err := ParseStatus(tt.input)

				// Assert
				a := assert.New(t)
				a.NoError(err)
				a.Equal(tt.wantStatus, gotStatus)
			})
		}
	})

	t.Run("異常系", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name       string
			input      string
			wantStatus Status
		}{
			{name: "空文字", input: "", wantStatus: ""},
			{name: "想定されていない値", input: "UnknownStatus", wantStatus: ""},
			{name: "大文字小文字の違い", input: "Draft", wantStatus: ""},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				// Arrange

				// Act
				gotStatus, err := ParseStatus(tt.input)

				// Assert
				a := assert.New(t)
				a.Error(err)
				a.Equal("invalid status: "+tt.input, err.Error())
				a.Equal(tt.wantStatus, gotStatus)
			})
		}
	})
}
