package e_utils

import (
	"context"
)

// Extracts and returns the context from an optional slice of contexts. If the slice is empty, it returns a new context.
func CtxOrDefault(slice []context.Context) context.Context {
	if len(slice) == 0 {
		return context.TODO()
	}
	return slice[0]
}
