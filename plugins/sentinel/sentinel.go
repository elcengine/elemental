package sentinel

import (
	"fmt"
	"reflect"
	elemental "github.com/elcengine/elemental/core"
	"github.com/go-playground/validator/v10"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
)

var validate = validator.New()

// Legitimize validates the input data based on the given validation tags within it's type definition. Basic validations are inherited from the go-playground/validator package while the augmented validations are provided by the sentinel package.
func Legitimize(input interface{}) error {
	err := validate.Struct(input)
	if err != nil {
		return err
	}
	v := reflect.ValueOf(input)
	for i := 0; i < v.NumField(); i++ {
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
			fieldName, _ := lo.Coalesce(field.Tag.Get("json"), field.Tag.Get("bson"), "_id")
			if len(definitionSections) > 1 {
				fieldName = definitionSections[1]
			}
			augmentedQuery := func(q elemental.Model[any]) elemental.Model[any] {
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
			default:
				return fmt.Errorf("Unknown augmented validation tag: %s", tag)
			}

		}
	}
	return nil
}