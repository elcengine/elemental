package e_utils

import (
	qkit "github.com/clubpay/qlubkit-go"
)

func ElementAtIndex[T any](slice []T, index int) T {
	if index < 0 || index >= len(slice) {
		var zero T
		return zero
	}
	return slice[index]
}

func FirstPtr[T any](slice []T) *T {
	if len(slice) == 0 {
		return nil
	}
	return &slice[0]
}

func First[T any](slice []T) T {
	return ElementAtIndex(slice, 0)
}

func CastSlice[T any](slice []any) []T {
	return qkit.Map(func(doc any) T {
		return qkit.Cast[T](doc)
	}, slice)
}

func CastSliceFromMaps[T any](slice []map[string]any) []T {
	var result []T
	for _, doc := range slice {
		result = append(result, qkit.CastJSON[T](doc))
	}
	return result
}
