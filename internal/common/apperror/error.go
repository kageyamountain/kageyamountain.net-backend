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

func NewErrorNoData(params ...string) error {
	return formatError(ErrorNoData, params...)
}

func NewErrorUnpublishedArticle(params ...string) error {
	return formatError(ErrorUnpublishedArticle, params...)
}

func formatError(err error, params ...string) error {
	if len(params) == 0 {
		return err
	}
	return fmt.Errorf("%w: %s", err, strings.Join(params, ", "))
}

func IsOneOf(err error, targets ...error) bool {
	for _, target := range targets {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}
