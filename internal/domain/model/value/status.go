package value

import "fmt"

type Status string

const (
	StatusDraft   Status = "draft"
	StatusPublish Status = "publish"
)

var validStatuses = map[Status]bool{
	StatusDraft:   true,
	StatusPublish: true,
}

func NewStatus(value string) (Status, error) {
	// 有効な値かをチェック
	status := Status(value)
	if !validStatuses[status] {
		return "", fmt.Errorf("invalid status: %s", value)
	}

	return status, nil
}

func (s Status) String() string {
	return string(s)
}
