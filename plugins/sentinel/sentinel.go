package sentinel

import (
	"fmt"
	"reflect"

	elemental "github.com/elcengine/elemental/core"
	"github.com/go-playground/validator/v10"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"

	// "strconv"
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
			enhancedQuery := func(q elemental.Model[any]) elemental.Model[any] {
				database := field.Tag.Get("database")
				if database != "" {
					q = q.SetDatabase(database)
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
				doc := enhancedQuery(elemental.NativeModel.FindOne(primitive.M{fieldName: value})).Exec()
				if doc != nil {
					return fmt.Errorf("Key: '%s.%s' Error:Field validation for '%s' failed on the '%s' tag", v.Type().Name(), field.Name, field.Name, tag)
				}
			default:
				return fmt.Errorf("Unknown augmented validation tag: %s", tag)
			}

		}
		// 	switch tag[0] {
		// 		case "exists":
		// 			params := strings.Split(tag[1], ",")
		// 			table := params[0]
		// 			column := params[1]
		// 			// Perform database lookup to check existence

		// 			// if input != nil {
		// 			// 	fieldValue := reflect.ValueOf(input).FieldByName(field.Name)
		// 			// 	if !checkExistence( column, fieldValue.Interface()) {
		// 			// 		return fmt.Errorf("%s does not exist in %s", column, table)
		// 			// 	}
		// 			// }
		// 			// fmt.Println("Checking if the %v exists in the database")
		// 			fmt.Printf("Checking if the %s with value %s exists in the %s table, column %s\n", field.Name, value, table, column)
		// 			// return nil

		// 		case "IsGreater":
		// 			params := strings.Split(tag[1], ",")
		// 			threshold, _ := strconv.Atoi(params[0])
		// 			//chek if the dataset has a larger value

		// 			// if !checkIsGreater(threshold, field.Name) {
		// 			// 	return fmt.Errorf("%s is not greater than %d", field.Name, threshold)
		// 			// }
		// 			fmt.Printf("Checking if the %s is greater than %d\n", field.Name, threshold)
		// 			// return nil

		// 		case "isTrue":
		// 			//pass the value of the field to checkIsTrue function

		// 			// key, _ := strconv.Atoi(v.Field(0).Interface().(string))
		// 			// if !checkIsTrue(field.Name, key) {
		// 			// 	return fmt.Errorf("%s is not true", field.Name)
		// 			// }
		// 			fmt.Printf("Checking if the %s is true\n", field.Name)
		// 			// return nil
		// 		}

	}
	return nil
}

// var users = data.Users

// func checkExistence( column string, value interface{}) bool {
// 	for _, user := range users {
// 		if reflect.ValueOf(user).FieldByName(column).Interface() == value {
// 			return true
// 		}
// 	}
// 	return false
// }

// func checkIsGreater(threshold int, feildName string) bool {
// 	for _, user := range users {
// 		if reflect.ValueOf(user).FieldByName(feildName).Interface().(int) > threshold {
// 			return true
// 		}
// 	}
// 	return false
// }

// func checkIsTrue(feildName string, id int) bool {
// 	for _, user := range users {
// 		if user.ID == id && reflect.ValueOf(user).FieldByName(feildName).Interface().(bool){
// 			return true
// 		}
// 	}
// 	return false
// }
