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
	var subpipeline any
	switch v := value.(type) {
	case primitive.M:
		path = utils.Cast[string](v["path"])
		selectField = v["select"]
		subpipeline = v["pipeline"]
	default:
		path = utils.Cast[string](value)
	}
	if path != "" {
		for i := range m.docReflectType.NumField() {
			field := m.docReflectType.Field(i)
			tag := cleanTag(field.Tag.Get("bson"))
			if path == field.Name {
				path = tag
			}
			if path == tag {
				fieldname = field.Name
				break
			}
		}
		if fieldname != "" {
			schemaField := m.Schema.Field(fieldname)
			if schemaField != nil {
				collection := schemaField.Collection
				if collection == "" {
					if schemaField.Ref != "" {
						model := reflect.ValueOf(Models[schemaField.Ref])
						collection = model.FieldByName("Schema").FieldByName("Options").FieldByName("Collection").String()
					}
				}
				if collection != "" {
					lookup := primitive.M{
						"from":         collection,
						"localField":   path,
						"foreignField": "_id",
						"as":           path,
					}
					if subpipeline != nil {
						lookup["pipeline"] = subpipeline
					} else if selectField != nil {
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
//
// It can accept a single string, a slice of strings, or a map with 'path' and optionally a 'select' or a 'pipeline' key.
func (m Model[T]) Populate(values ...any) Model[T] {
	m.setResult([]bson.M{})
	m.executor = func(m Model[T], ctx context.Context) any {
		cursor := lo.Must(m.Collection().Aggregate(ctx, m.pipeline))
		lo.Must0(cursor.All(ctx, m.result))
		m.checkConditionsAndPanic(m.result)
		return m.result
	}
	if len(values) == 1 {
		var parts []string
		switch v := values[0].(type) {
		case []string:
			parts = v
		case string:
			if strings.Contains(v, ",") || strings.Contains(v, " ") {
				parts = strings.FieldsFunc(v, func(r rune) bool {
					return r == ',' || r == ' '
				})
			}
		}
		if len(parts) > 0 {
			for _, value := range parts {
				m = m.populate(value)
			}
			return m
		}
	}
	for _, value := range values {
		m = m.populate(value)
	}
	return m
}
