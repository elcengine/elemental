package utils

import (
	"github.com/samber/lo"
)

func CastSlice[T any](slice []any) []T {
	return lo.Map(slice, func(doc any, _ int) T {
		return Cast[T](doc)
	})
}

func CastBSONSlice[T any](slice []any) []T {
	return lo.Map(slice, func(doc any, _ int) T {
		return CastBSON[T](doc)
	})
}
