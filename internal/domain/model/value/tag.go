package value

import "fmt"

type Tag struct {
	value string
}

func NewTag(value string) (Tag, error) {
	// 有効な値かをチェック
	_, exists := validTags[value]
	if !exists {
		return Tag{}, fmt.Errorf("invalid tag: %s", value)
	}

	return Tag{value: value}, nil
}

func (t Tag) Value() string {
	return t.value
}

const (
	TagGo  = "Go"
	TagGin = "Gin"

	TagJavaScript = "JavaScript"
	TagTypeScript = "TypeScript"
	TagVueJS      = "Vue.js"

	TagAWS      = "AWS"
	TagDynamoDB = "DynamoDB"
)

var validTags = map[string]struct{}{
	TagGo:         {},
	TagGin:        {},
	TagJavaScript: {},
	TagTypeScript: {},
	TagVueJS:      {},
	TagAWS:        {},
	TagDynamoDB:   {},
}
