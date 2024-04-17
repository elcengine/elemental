package elemental

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/clubpay/qlubkit-go"
)

func enforceSchema(schema Schema, entity any) {
	reflectedEntity := reflect.ValueOf(entity)
	for field, definition := range schema.Definitions {
		if (reflectedEntity.FieldByName(field).IsZero()) {
			if definition.Required {
				panic(fmt.Sprintf("Field %s is required", field))
			}
			if (definition.Default != nil) {
				reflectedEntity.FieldByName(field).Set(reflect.ValueOf(definition.Default))
			}
		}
		if (definition.Type != reflect.Invalid && reflectedEntity.FieldByName(field).Kind() != definition.Type) {
			panic(fmt.Sprintf("Field %s has an invalid type. It must be of type %s", field, definition.Type.String()))
		}
		if (definition.Min != 0 && reflectedEntity.FieldByName(field).Float() < definition.Min) {
			panic(fmt.Sprintf("Field %s must be greater than %f", field, definition.Min))
		}
		if (definition.Max != 0 && reflectedEntity.FieldByName(field).Float() > definition.Max) {
			panic(fmt.Sprintf("Field %s must be less than %f", field, definition.Max))
		}
		if (definition.Length != 0 && int64(len(reflectedEntity.FieldByName(field).String())) > definition.Length) {
			panic(fmt.Sprintf("Field %s must be less than %d characters", field, definition.Length))
		}
		if (definition.Regex != "" && qkit.Must(regexp.Match(definition.Regex, reflectedEntity.FieldByName(field).Bytes()))) {
			panic(fmt.Sprintf("Field %s must match the regex pattern %s", field, definition.Regex))
		}
	}
}