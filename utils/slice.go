package utils

import (
	"github.com/samber/lo"
)

func CastBSONSlice[T any](slice []any) []T {
	return lo.Map(slice, func(doc any, _ int) T {
		return CastBSON[T](doc)
	})
}
