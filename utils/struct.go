//nolint:exhaustive
package e_utils

import (
	"fmt"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"strconv"
)

func setField(field reflect.Value, val string) error {
	if !field.CanSet() {
		return fmt.Errorf("can't set value")
	}
	switch field.Kind() {
	case reflect.Int:
		if val, err := strconv.ParseInt(val, 10, 64); err == nil {
			field.Set(reflect.ValueOf(int(val)).Convert(field.Type()))
		}
	case reflect.String:
		field.Set(reflect.ValueOf(val).Convert(field.Type()))
	}
	return nil
}

func SetDefaults(ptr any) error {
	if reflect.TypeOf(ptr).Kind() != reflect.Ptr {
		return fmt.Errorf("not a pointer")
	}
	v := reflect.ValueOf(ptr).Elem()
	t := v.Type()
	for i := range t.NumField() {
		if defaultVal := t.Field(i).Tag.Get("default"); defaultVal != "-" {
			if err := setField(v.Field(i), defaultVal); err != nil {
				return err
			}
		}
	}
	return nil
}

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
	var objectId primitive.ObjectID
	if reflectedValueType == reflect.TypeOf(&objectId) || reflectedValueType == reflect.TypeOf(objectId) {
		return value.(primitive.ObjectID).IsZero()
	}
	return false
}
