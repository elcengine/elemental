package sentinel

import (
	"fmt"
	"reflect"
	"strings"

	elemental "github.com/elcengine/elemental/core"
	e_utils "github.com/elcengine/elemental/utils"
	"github.com/go-playground/validator/v10"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var validate = validator.New()

// Legitimize validates the input data based on the given validation tags within it's type definition.
// Basic validations are inherited from the go-playground/validator package while the augmented validations
// are provided by the sentinel package.
func Legitimize(input any) error {
	err := validate.Struct(input)
	if err != nil {
		return err
	}
	v := reflect.ValueOf(input)

	inputMap := e_utils.ToMap(input)

	for i := range v.NumField() {
		field := v.Type().Field(i)
		value := v.Field(i).Interface()
		vts := strings.Split(field.Tag.Get("augmented_validate"), ";")
		for _, vt := range vts {
			if vt == "" {
				continue
			}
			tagSections := strings.Split(vt, "=")
			tag := tagSections[0]
			definition := tagSections[1]
			definitionSections := strings.Split(definition, "->")
			fieldName := lo.CoalesceOrEmpty(field.Tag.Get("json"), field.Tag.Get("bson"), "_id")
			if len(definitionSections) > 1 {
				fieldName = definitionSections[1]
			}
			reference := field.Tag.Get("ref")

			augmentedQuery := func(q elemental.Model[map[string]any]) elemental.Model[map[string]any] {
				database := field.Tag.Get("database")
				if database != "" {
					q = q.SetDatabase(database)
				}
				connection := field.Tag.Get("connection")
				if connection != "" {
					q = q.SetConnection(connection)
				}
				modelOrCollection := definitionSections[0]
				collection := modelOrCollection
				modelFromCache := elemental.Models[modelOrCollection]
				if modelFromCache != nil {
					collection = reflect.ValueOf(modelFromCache).FieldByName("Schema").FieldByName("Options").FieldByName("Collection").String()
				}
				q = q.SetCollection(collection)
				return q
			}

			getReferenceField := func() string {
				if reference != "" {
					return reference
				}
				return fieldName
			}

			getReferenceFieldValue := func() any {
				if reference != "" {
					return inputMap[reference]
				}
				return value
			}

			switch tag {
			case "unique":
				doc := augmentedQuery(elemental.NativeModel.FindOne(primitive.M{fieldName: value})).Exec()
				if doc != nil {
					return fmt.Errorf("Key: '%s.%s' Error:Field validation for '%s' failed on the '%s' tag", v.Type().Name(), field.Name, field.Name, tag)
				}
			case "exists":
				doc := augmentedQuery(elemental.NativeModel.FindOne(primitive.M{fieldName: value})).Exec()
				if doc == nil {
					return fmt.Errorf("Key: '%s.%s' Error:Field validation for '%s' failed on the '%s' tag", v.Type().Name(), field.Name, field.Name, tag)
				}
			case "greater_than":
				doc := augmentedQuery(elemental.NativeModel.FindOne(primitive.M{getReferenceField(): getReferenceFieldValue()})).Exec()
				if doc == nil || e_utils.LTE(value, e_utils.Cast[map[string]any](doc)[fieldName]) {
					return fmt.Errorf("Key: '%s.%s' Error:Field validation for '%s' failed on the '%s' tag", v.Type().Name(), field.Name, field.Name, tag)
				}
			case "greater_than_or_equal_to":
				doc := augmentedQuery(elemental.NativeModel.FindOne(primitive.M{getReferenceField(): getReferenceFieldValue()})).Exec()
				if doc == nil || e_utils.LT(value, e_utils.Cast[map[string]any](doc)[fieldName]) {
					return fmt.Errorf("Key: '%s.%s' Error:Field validation for '%s' failed on the '%s' tag", v.Type().Name(), field.Name, field.Name, tag)
				}
			case "less_than":
				doc := augmentedQuery(elemental.NativeModel.FindOne(primitive.M{getReferenceField(): getReferenceFieldValue()})).Exec()
				if doc == nil || e_utils.GTE(value, e_utils.Cast[map[string]any](doc)[fieldName]) {
					return fmt.Errorf("Key: '%s.%s' Error:Field validation for '%s' failed on the '%s' tag", v.Type().Name(), field.Name, field.Name, tag)
				}
			case "less_than_or_equal_to":
				doc := augmentedQuery(elemental.NativeModel.FindOne(primitive.M{getReferenceField(): getReferenceFieldValue()})).Exec()
				if doc == nil || e_utils.GT(value, e_utils.Cast[map[string]any](doc)[fieldName]) {
					return fmt.Errorf("Key: '%s.%s' Error:Field validation for '%s' failed on the '%s' tag", v.Type().Name(), field.Name, field.Name, tag)
				}
			case "equals":
				doc := augmentedQuery(elemental.NativeModel.FindOne(primitive.M{getReferenceField(): getReferenceFieldValue()})).Exec()
				if doc == nil || !e_utils.EQ(value, e_utils.Cast[map[string]any](doc)[fieldName]) {
					return fmt.Errorf("Key: '%s.%s' Error:Field validation for '%s' failed on the '%s' tag", v.Type().Name(), field.Name, field.Name, tag)
				}
			case "not_equals":
				doc := augmentedQuery(elemental.NativeModel.FindOne(primitive.M{getReferenceField(): getReferenceFieldValue()})).Exec()
				if doc != nil && e_utils.EQ(value, e_utils.Cast[map[string]any](doc)[fieldName]) {
					return fmt.Errorf("Key: '%s.%s' Error:Field validation for '%s' failed on the '%s' tag", v.Type().Name(), field.Name, field.Name, tag)
				}
			default:
				return fmt.Errorf("Unknown augmented validation tag: %s", tag)
			}
		}
	}
	return nil
}
