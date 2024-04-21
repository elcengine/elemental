package elemental

import (
	"context"
	"elemental/connection"
	"reflect"
	"github.com/clubpay/qlubkit-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
)

type ModelSkeleton[T any] interface {
	Schema() Schema
	Create() primitive.ObjectID
	FindOne(query primitive.M) *T
}

type Model[T any] struct {
	schema   Schema
	pipeline mongo.Pipeline
	returnSingleRecord bool
}

var models = make(map[string]Model[any])

func NewModel[T any](name string, schema Schema) Model[T] {
	var sample [0]T
	if _, ok := models[name]; ok {
		return qkit.Cast[Model[T]](models[name])
	}
	model := Model[T]{schema: schema}
	models[name] = qkit.Cast[Model[any]](model)
	e_connection.On(event.ConnectionReady, func() {
		schema.syncIndexes(reflect.TypeOf(sample).Elem())
	})
	return model
}

func (m Model[T]) Create(doc T) T {
	document := enforceSchema(m.schema, &doc)
	qkit.Must(m.Collection().InsertOne(context.TODO(), document))
	return document
}

func (m Model[T]) Find(query *primitive.M) Model[T] {
	m.pipeline = append(m.pipeline, bson.D{{"$match", query}})
	return m
}

func (m Model[T]) FindOne(query *primitive.M) *T {
	return nil
}

func (m Model[T]) Exec() any {
	cursor := qkit.Must(m.Collection().Aggregate(context.TODO(), m.pipeline))
	var results []T
	if err := cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	if m.returnSingleRecord {
		if len(results) == 0 {
			return nil
		}
		return results[0]
	}
	return results
}

func (m Model[T]) Validate() error {
	return nil
}

func (m Model[T]) ValidateField() error {
	return nil
}

func (m Model[T]) Collection() *mongo.Collection {
	return e_connection.Use(m.schema.Options.Database, m.schema.Options.Connection).Collection(m.schema.Options.Collection)
}
