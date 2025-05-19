package utils

import (
	"reflect"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var dateTime primitive.DateTime
var objectID primitive.ObjectID

var reflectTypeDateTimePtr = reflect.TypeOf(&dateTime)
var reflectTypeDateTime = reflect.TypeOf(dateTime)

var reflectTypeObjectIDPtr = reflect.TypeOf(&objectID)
var reflectTypeObjectID = reflect.TypeOf(objectID)

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
	if reflectedValueType == reflectTypeDateTimePtr || reflectedValueType == reflectTypeDateTime {
		return value.(primitive.DateTime).Time().IsZero()
	}
	if reflectedValueType == reflectTypeObjectIDPtr || reflectedValueType == reflectTypeObjectID {
		return value.(primitive.ObjectID).IsZero()
	}
	return false
}
