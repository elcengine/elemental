package elemental

import (
	"fmt"
	"reflect"

	"regexp"
	"time"

	"github.com/clubpay/qlubkit-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func enforceSchema[T any](schema Schema, doc *T, defaults ...bool) T {
	reflectedEntity := reflect.ValueOf(doc)
	if reflectedEntity.Kind() == reflect.Ptr {
		reflectedEntity = reflect.Indirect(reflectedEntity)
	}
	for field, definition := range schema.Definitions {
		reflectedField := reflectedEntity.FieldByName(field)
		if !reflectedField.IsValid() || reflectedField.IsZero() {
			if definition.Required {
				panic(fmt.Sprintf("Field %s is required", field))
			}
			if definition.Default != nil {
				reflectedField.Set(reflect.ValueOf(definition.Default))
			}
		}
		if definition.Type != reflect.Invalid && reflectedField.Kind() != definition.Type {
			panic(fmt.Sprintf("Field %s has an invalid type. It must be of type %s", field, definition.Type.String()))
		}
		if definition.Min != 0 && reflectedField.Float() < definition.Min {
			panic(fmt.Sprintf("Field %s must be greater than %f", field, definition.Min))
		}
		if definition.Max != 0 && reflectedField.Float() > definition.Max {
			panic(fmt.Sprintf("Field %s must be less than %f", field, definition.Max))
		}
		if definition.Length != 0 && int64(len(reflectedField.String())) > definition.Length {
			panic(fmt.Sprintf("Field %s must be less than %d characters", field, definition.Length))
		}
		if definition.Regex != "" && qkit.Must(regexp.Match(definition.Regex, reflectedField.Bytes())) {
			panic(fmt.Sprintf("Field %s must match the regex pattern %s", field, definition.Regex))
		}
	}
	if (len(defaults) == 0 || defaults[0]) {
		SetDefault(&reflectedEntity, "ID", primitive.NewObjectID())
		SetDefault(&reflectedEntity, "CreatedAt", time.Now())
		SetDefault(&reflectedEntity, "UpdatedAt", time.Now())
	}
	return qkit.Cast[T](reflectedEntity.Interface())
}

func SetDefault[T any](reflectedEntity *reflect.Value, fieldName string, defaultValue T) {
	field := reflectedEntity.FieldByName(fieldName)
	if field.IsValid() && field.IsZero() {
		if (field.Kind() == reflect.Ptr) {
			field.Set(reflect.ValueOf(qkit.ValPtr(defaultValue)))
		} else {
			field.Set(reflect.ValueOf(defaultValue))
		}
	}
}
