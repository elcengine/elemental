package utils

import (
	"reflect"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func IsEmpty(value any) bool {
	if value == nil {
		return true
	}
	if lo.IsEmpty(value) {
		return true
	}
	if dt, ok := value.(primitive.DateTime); ok {
		return dt.Time().IsZero()
	}
	if oid, ok := value.(primitive.ObjectID); ok {
		return oid.IsZero()
	}
	reflectedValue := reflect.ValueOf(value)
	if !reflectedValue.IsValid() || reflectedValue.IsZero() {
		return true
	}
	return false
}
