package utils

import (
	"encoding/json"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
)

// Converts any type to a given type. If conversion fails, it returns the zero value of the given type.
func Cast[T any](val any) T {
	if val, ok := val.(T); ok {
		return val
	}
	var zero T
	return zero
}

// Converts any type to a map[string]any.
func ToMap(s any) map[string]any {
	m := make(map[string]any)
	json.Unmarshal(lo.Must(json.Marshal(s)), &m)
	return m
}

// Converts any type to a given type based on their bson representation. It partially fills the target in case they are not directly compatible.
func CastBSON[T any](val any) T {
	return FromBSON[T](ToBSON(val))
}

// Converts a given value to a byte array.
func ToBSON(val any) []byte {
	bytes, _ := bson.Marshal(val)
	return bytes
}

// Converts a byte array to a given type.
func FromBSON[T any](bytes []byte) T {
	var v T
	bson.Unmarshal(bytes, &v)
	return v
}
