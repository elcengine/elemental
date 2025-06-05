package elemental

import (
	"context"
	"reflect"
	"strings"

	"github.com/elcengine/elemental/utils"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (m Model[T]) populate(value any) Model[T] {
	var path, fieldname string
	var selectField any
	switch v := value.(type) {
	case primitive.M:
		path = utils.Cast[string](v["path"])
		selectField = v["select"]
	default:
		path = utils.Cast[string](value)
	}
	if path != "" {
		for i := range m.docReflectType.NumField() {
			field := m.docReflectType.Field(i)
			if cleanTag(field.Tag.Get("bson")) == path {
				fieldname = field.Name
				break
			}
		}
		if fieldname != "" {
			schemaField := m.Schema.Field(fieldname)
			if schemaField != nil {
				collection := schemaField.Collection
				if lo.IsEmpty(collection) {
					if !lo.IsEmpty(schemaField.Ref) {
						model := reflect.ValueOf(Models[schemaField.Ref])
						collection = model.FieldByName("Schema").FieldByName("Options").FieldByName("Collection").String()
					}
				}
				if !lo.IsEmpty(collection) {
					lookup := primitive.M{
						"from":         collection,
						"localField":   path,
						"foreignField": "_id",
						"as":           path,
					}
					if selectField != nil {
						lookup["pipeline"] = []primitive.M{
							{"$project": selectField},
						}
					}
					m.pipeline = append(m.pipeline, bson.D{{Key: "$lookup", Value: lookup}})
					if schemaField.Type != reflect.Slice {
						unwind := primitive.M{
							"path":                       "$" + path,
							"preserveNullAndEmptyArrays": true,
						}
						m.pipeline = append(m.pipeline, bson.D{{Key: "$unwind", Value: unwind}})
					} else {
						m.pipeline = append(m.pipeline, bson.D{{Key: "$lookup", Value: lookup}})
					}
				}
			}
		}
	}
	return m
}

// Finds and attaches the referenced documents to the main document returned by the query.
// The fields to populate must have a 'Collection' or 'Ref' property in their schema definition.
func (m Model[T]) Populate(values ...any) Model[T] {
	if len(values) == 1 {
		if str, ok := values[0].(string); ok && (strings.Contains(str, ",") || strings.Contains(str, " ")) {
			parts := strings.FieldsFunc(str, func(r rune) bool {
				return r == ',' || r == ' '
			})
			for _, value := range parts {
				m = m.populate(value)
			}
			return m
		}
	}
	for _, value := range values {
		m = m.populate(value)
	}
	m.executor = func(m Model[T], ctx context.Context) any {
		var results []bson.M
		cursor := lo.Must(m.Collection().Aggregate(ctx, m.pipeline))
		lo.Must0(cursor.All(ctx, &results))
		m.checkConditionsAndPanic(results)
		return results
	}
	return m
}
