package elemental

import (
	"elemental/utils"
	"fmt"
	"reflect"

	"regexp"
	"time"

	"github.com/clubpay/qlubkit-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func enforceSchema[T any](schema Schema, doc *T) primitive.M {
	reflectedEntity := reflect.ValueOf(doc)
	if reflectedEntity.Kind() == reflect.Ptr {
		reflectedEntity = reflect.Indirect(reflectedEntity)
	}
	for field, definition := range schema.Definitions {
		reflectedField := reflect.Indirect(reflectedEntity.FieldByName("Base")).FieldByName(field)
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
	extractBaseValueOrSetDefault(&reflectedEntity, "ID", primitive.NewObjectID())
	extractBaseValueOrSetDefault(&reflectedEntity, "CreatedAt", time.Now())
	extractBaseValueOrSetDefault(&reflectedEntity, "UpdatedAt", time.Now())
	result := qkit.PtrVal(e_utils.ToBSON(reflectedEntity.FieldByName("Base").Interface()))
	for k, v := range qkit.PtrVal(e_utils.ToBSON(reflectedEntity.Interface())) {
		if (k != "base") {
			result[k] = v
		}
	}
	return result
}

func extractBaseValueOrSetDefault[T any](reflectedEntity *reflect.Value, fieldName string, defaultValue T) {
	field := reflect.Indirect(reflectedEntity.FieldByName("Base")).FieldByName(fieldName)
	if field.IsValid() && !field.IsZero() {
		reflectedEntity.FieldByName(fieldName).Set(reflect.ValueOf(field))
	}
	field = reflectedEntity.FieldByName(fieldName)
	if field.IsValid() && field.IsZero() {
		field.Set(reflect.ValueOf(qkit.ValPtr(defaultValue)))
	}
}
