package elemental

import (
	"reflect"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UTILITY FUNCTIONS. MOVE LATER
func structToMap(obj any) map[string]any {
	result := make(map[string]any)
	v := reflect.ValueOf(obj)

	// Check if the input is a struct or a pointer to a struct
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Iterate through the fields of the struct
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		fieldName := field.Name
		fieldValue := v.Field(i).Interface()
		result[fieldName] = fieldValue
	}
	return result
}

// Function to convert map[string]Field to map[string]any
func convertMap(original map[string]Field) map[string]any {
	result := make(map[string]any)
	for key, field := range original {
		result[key] = field
	}
	return result
}

// Function to convert primitive.D to map[string]any
func convertDToMap(doc primitive.D) map[string]any {
	result := make(map[string]any)
	for _, elem := range doc {
		result[elem.Key] = elem.Value
	}
	return result
}
