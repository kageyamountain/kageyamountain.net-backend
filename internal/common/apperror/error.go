package apperror

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrorNoData             = errors.New("error no data")
	ErrorUnpublishedArticle = errors.New("error unpublished article")
)

func NewErrorNodata(params ...string) error {
	if len(params) == 0 {
		return fmt.Errorf("%w", ErrorNoData)
	}
	return fmt.Errorf("%w: %s", ErrorNoData, strings.Join(params, ", "))
}

func NewErrorUnpublishedArticle(params ...string) error {
	if len(params) == 0 {
		return fmt.Errorf("%w", ErrorUnpublishedArticle)
	}
	return fmt.Errorf("%w: %s", ErrorUnpublishedArticle, strings.Join(params, ", "))
}

func IsOneOf(err error, targets ...error) bool {
	for _, target := range targets {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}
