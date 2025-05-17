package elemental

import (
	"fmt"
	"maps"
	"reflect"
	"strings"

	"github.com/elcengine/elemental/utils"
	"github.com/spf13/cast"

	"regexp"
	"time"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func enforceSchema[T any](schema Schema, doc *T, reflectedEntityType *reflect.Type, defaults ...bool) (bson.M, bson.M) {
	var entityToInsert bson.M

	// Fast return when bypass schema enforcement or value is not a struct
	if doc != nil && (reflect.TypeOf(doc).Elem().Kind() != reflect.Struct || schema.Options.BypassSchemaEnforcement) {
		entityToInsert = *e_utils.ToBSONDoc(doc)
		return entityToInsert, entityToInsert
	}

	if reflectedEntityType != nil {
		entityToInsert = e_utils.Cast[bson.M](doc)
		if entityToInsert == nil {
			entityToInsert = make(bson.M)
		}
	} else {
		entityToInsert = *e_utils.ToBSONDoc(doc)
		reflectedEntityType = lo.ToPtr(reflect.TypeOf(doc).Elem())
	}

	if len(defaults) == 0 || defaults[0] {
		for _, field := range []string{"ID", "CreatedAt", "UpdatedAt"} {
			if reflectedField, ok := (*reflectedEntityType).FieldByName(field); ok && reflectedField.Type != nil {
				key := cleanTag(reflectedField.Tag.Get("bson"))
				if e_utils.IsEmpty(entityToInsert[key]) {
					switch field {
					case "ID":
						entityToInsert[key] = primitive.NewObjectID()
					case "CreatedAt", "UpdatedAt":
						entityToInsert[key] = time.Now()
					}
				}
			}
		}
	}

	detailedEntity := make(bson.M)
	maps.Copy(detailedEntity, entityToInsert)

	for field, definition := range schema.Definitions {
		reflectedField, ok := (*reflectedEntityType).FieldByName(field)
		if !ok {
			continue
		}
		fieldBsonName := cleanTag(reflectedField.Tag.Get("bson"))
		val := entityToInsert[fieldBsonName]

		// Required and default checks
		if e_utils.IsEmpty(val) {
			if definition.Required {
				panic(fmt.Errorf("Field %s is required", field))
			}
			if definition.Default != nil {
				entityToInsert[fieldBsonName] = definition.Default
				detailedEntity[fieldBsonName] = definition.Default
				val = definition.Default
			}
		}

		// Type check
		if definition.Type != reflect.Invalid {
			actualKind := reflectedField.Type.Kind()
			if actualKind == reflect.Ptr {
				actualKind = reflectedField.Type.Elem().Kind()
			}
			if actualKind != definition.Type {
				panic(fmt.Errorf("Field %s has an invalid type. It must be of type %s", field, definition.Type.String()))
			}
		}

		// Nested schema validation
		if definition.Type == reflect.Struct && definition.Schema != nil {
			subdocumentField := reflectedField
			entityToInsert[fieldBsonName], detailedEntity[fieldBsonName] =
				enforceSchema(*definition.Schema, e_utils.Cast[*bson.M](val), &subdocumentField.Type, false)
			continue
		}

		// Reference/Collection
		if definition.Type == reflect.Struct && (definition.Ref != "" || definition.Collection != "") && val != nil {
			subdocumentField := reflectedField
			if subdocumentIdField, ok := subdocumentField.Type.FieldByName("ID"); ok {
				entityToInsert = lo.Assign(
					entityToInsert,
					bson.M{
						fieldBsonName: val.(primitive.M)[subdocumentIdField.Tag.Get("bson")],
					},
				)
			}
			continue
		}

		if definition.Min != 0 {
			if v := cast.ToFloat64(val); v < definition.Min {
				panic(fmt.Errorf("Field %s must be greater than or equal to %v", field, definition.Min))
			}
		}
		if definition.Max != 0 {
			if v := cast.ToFloat64(val); v > definition.Max {
				panic(fmt.Errorf("Field %s must be less than or equal to %v", field, definition.Max))
			}
		}
		if definition.Length != 0 {
			if s := cast.ToString(val); int64(len(s)) > definition.Length {
				panic(fmt.Errorf("Field %s must be less than or equal to %d characters", field, definition.Length))
			}
		}
		if definition.Regex != "" {
			if matched := lo.Must(regexp.MatchString(definition.Regex, cast.ToString(val))); !matched {
				panic(fmt.Errorf("Field %s must match the regex pattern %s", field, definition.Regex))
			}
		}
	}
	return entityToInsert, detailedEntity
}

func cleanTag(tag string) string {
	return strings.ReplaceAll(tag, ",omitempty", "")
}
