package e_utils

func ElementAtIndex[T any](slice []T, index int) T {
	if index < 0 || index >= len(slice) {
		var zero T
		return zero
	}
	return slice[index]
}