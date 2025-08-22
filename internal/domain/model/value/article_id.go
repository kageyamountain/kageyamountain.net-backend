package value

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type ArticleID struct {
	// ハイフン無しのUUIDv4形式
	value string
}

func NewArticleID(value string) (*ArticleID, error) {
	err := uuid.Validate(value)
	if err != nil {
		return nil, fmt.Errorf("invalid format of article_id. article_id: %s error: %w", value, err)
	}

	if strings.Contains(value, "-") {
		return nil, fmt.Errorf("article_id cannot contain hyphen. article_id: %s", value)
	}

	return &ArticleID{
		value: value,
	}, nil
}

func GenerateArticleID() *ArticleID {
	uuidWithHyphen := uuid.New().String()
	value := strings.ReplaceAll(uuidWithHyphen, "-", "")

	return &ArticleID{
		value: value,
	}
}

func (a *ArticleID) Value() string {
	return a.value
}
