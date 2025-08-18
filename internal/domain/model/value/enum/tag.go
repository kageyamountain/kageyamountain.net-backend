package enum

import "fmt"

type Tag string

const (
	TagGo  Tag = "Go"
	TagGin Tag = "Gin"

	TagJavaScript Tag = "JavaScript"
	TagTypeScript Tag = "TypeScript"
	TagVueJS      Tag = "Vue.js"

	TagAWS      Tag = "AWS"
	TagDynamoDB Tag = "DynamoDB"
)

var validTags = map[Tag]bool{
	TagGo:         true,
	TagGin:        true,
	TagJavaScript: true,
	TagTypeScript: true,
	TagVueJS:      true,
	TagAWS:        true,
	TagDynamoDB:   true,
}

func NewTag(value string) (Tag, error) {
	// 有効な値かをチェック
	tag := Tag(value)
	if !validTags[tag] {
		return "", fmt.Errorf("invalid tag: %s", value)
	}

	return tag, nil
}

func (t Tag) String() string {
	return string(t)
}
