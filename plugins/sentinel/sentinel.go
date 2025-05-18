package sentinel

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/elcengine/elemental/core"
	"github.com/elcengine/elemental/utils"
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

	inputMap := utils.ToMap(input)

	for i := range v.NumField() {
		field := v.Type().Field(i)
		value := v.Field(i).Interface()
		tags := strings.Split(field.Tag.Get("augmented_validate"), ";")
		for _, t := range tags {
			if t == "" {
				continue
			}
			tagSections := strings.Split(t, "=")
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
				return lo.CoalesceOrEmpty(reference, fieldName)
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
					return NewFieldError(v.Type().Name(), field.Name, tag)
				}
			case "exists":
				doc := augmentedQuery(elemental.NativeModel.FindOne(primitive.M{fieldName: value})).Exec()
				if doc == nil {
					return NewFieldError(v.Type().Name(), field.Name, tag)
				}
			case "greater_than":
				doc := augmentedQuery(elemental.NativeModel.FindOne(primitive.M{getReferenceField(): getReferenceFieldValue()})).Exec()
				if doc == nil || utils.LTE(value, utils.Cast[map[string]any](doc)[fieldName]) {
					return NewFieldError(v.Type().Name(), field.Name, tag)
				}
			case "greater_than_or_equal_to":
				doc := augmentedQuery(elemental.NativeModel.FindOne(primitive.M{getReferenceField(): getReferenceFieldValue()})).Exec()
				if doc == nil || utils.LT(value, utils.Cast[map[string]any](doc)[fieldName]) {
					return NewFieldError(v.Type().Name(), field.Name, tag)
				}
			case "less_than":
				doc := augmentedQuery(elemental.NativeModel.FindOne(primitive.M{getReferenceField(): getReferenceFieldValue()})).Exec()
				if doc == nil || utils.GTE(value, utils.Cast[map[string]any](doc)[fieldName]) {
					return NewFieldError(v.Type().Name(), field.Name, tag)
				}
			case "less_than_or_equal_to":
				doc := augmentedQuery(elemental.NativeModel.FindOne(primitive.M{getReferenceField(): getReferenceFieldValue()})).Exec()
				if doc == nil || utils.GT(value, utils.Cast[map[string]any](doc)[fieldName]) {
					return NewFieldError(v.Type().Name(), field.Name, tag)
				}
			case "equals":
				doc := augmentedQuery(elemental.NativeModel.FindOne(primitive.M{getReferenceField(): getReferenceFieldValue()})).Exec()
				if doc == nil || !utils.EQ(value, utils.Cast[map[string]any](doc)[fieldName]) {
					return NewFieldError(v.Type().Name(), field.Name, tag)
				}
			case "not_equals":
				doc := augmentedQuery(elemental.NativeModel.FindOne(primitive.M{getReferenceField(): getReferenceFieldValue()})).Exec()
				if doc != nil && utils.EQ(value, utils.Cast[map[string]any](doc)[fieldName]) {
					return NewFieldError(v.Type().Name(), field.Name, tag)
				}
			default:
				return fmt.Errorf("unknown augmented validation tag: %s", tag)
			}
		}
	}
	return nil
}

// NewFieldError creates a field error message for the given namespace, field, and tag.
// It returns an error with a formatted message indicating the validation failure.
func NewFieldError(namespace, field, tag string) error {
	return fmt.Errorf("Key: '%s.%s' Error:Field validation for '%s' failed on the '%s' tag", namespace, field, field, tag)
}
