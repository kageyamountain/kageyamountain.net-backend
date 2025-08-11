package value

import "fmt"

type Status struct {
	value string
}

func NewStatus(value string) (Status, error) {
	// 有効な値かをチェック
	_, exists := validStatuses[value]
	if !exists {
		return Status{}, fmt.Errorf("invalid status: %s", value)
	}

	return Status{value: value}, nil
}

func (s Status) Value() string {
	return s.value
}

const (
	StatusDraft   = "draft"
	StatusPublish = "publish"
)

var validStatuses = map[string]struct{}{
	StatusDraft:   {},
	StatusPublish: {},
}
