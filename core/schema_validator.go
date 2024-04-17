package elemental

import (
	"fmt"
	"reflect"
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
		
	}
}