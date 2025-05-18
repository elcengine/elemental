package elemental

import "errors"

var (
	ErrURIRequired               = errors.New("URI is required")
	ErrInvalidConnectionArgument = errors.New("invalid connection argument")
	ErrMustPairSortArguments     = errors.New("sort arguments must be in pairs")
)
