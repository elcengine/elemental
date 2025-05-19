package utils

import (
	"fmt"
	"github.com/spf13/cast"
)

func LT(a, b any) bool {
	switch a.(type) {
	case int, int8, int16, int32, int64, float32, float64:
		return cast.ToFloat64(a) < cast.ToFloat64(b)
	case string:
		return cast.ToString(a) < cast.ToString(b)
	default:
		panic(fmt.Errorf("unsupported type: %T", a))
	}
}

func LTE(a, b any) bool {
	switch a.(type) {
	case int, int8, int16, int32, int64, float32, float64:
		return cast.ToFloat64(a) <= cast.ToFloat64(b)
	case string:
		return cast.ToString(a) <= cast.ToString(b)
	default:
		panic(fmt.Errorf("unsupported type: %T", a))
	}
}

func GT(a, b any) bool {
	switch a.(type) {
	case int, int8, int16, int32, int64, float32, float64:
		return cast.ToFloat64(a) > cast.ToFloat64(b)
	case string:
		return cast.ToString(a) > cast.ToString(b)
	default:
		panic(fmt.Errorf("unsupported type: %T", a))
	}
}

func GTE(a, b any) bool {
	switch a.(type) {
	case int, int8, int16, int32, int64, float32, float64:
		return cast.ToFloat64(a) >= cast.ToFloat64(b)
	case string:
		return cast.ToString(a) >= cast.ToString(b)
	default:
		panic(fmt.Errorf("unsupported type: %T", a))
	}
}

func EQ(a, b any) bool {
	switch a.(type) {
	case int, int8, int16, int32, int64, float32, float64:
		return cast.ToFloat64(a) == cast.ToFloat64(b)
	case string:
		return cast.ToString(a) == cast.ToString(b)
	case bool:
		return cast.ToBool(a) == cast.ToBool(b)
	default:
		panic(fmt.Errorf("unsupported type: %T", a))
	}
}
