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
	reflectedValue := reflect.ValueOf(value)
	if !reflectedValue.IsValid() || reflectedValue.IsZero() {
		return true
	}
	reflectedValueType := reflect.TypeOf(value)
	var dateTime primitive.DateTime
	if reflectedValueType == reflect.TypeOf(&dateTime) || reflectedValueType == reflect.TypeOf(dateTime) {
		return value.(primitive.DateTime).Time().IsZero()
	}
	var objectID primitive.ObjectID
	if reflectedValueType == reflect.TypeOf(&objectID) || reflectedValueType == reflect.TypeOf(objectID) {
		return value.(primitive.ObjectID).IsZero()
	}
	return false
}
