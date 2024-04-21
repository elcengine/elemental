package elemental

import (
	"context"
	"elemental/connection"
	"elemental/utils"
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
	Name               string
	schema             Schema
	pipeline           mongo.Pipeline
	executor           func(ctx context.Context) any
}

var models = make(map[string]Model[any])

func NewModel[T any](name string, schema Schema) Model[T] {
	var sample [0]T
	if _, ok := models[name]; ok {
		return qkit.Cast[Model[T]](models[name])
	}
	model := Model[T]{
		Name:   name,
		schema: schema,
	}
	models[name] = qkit.Cast[Model[any]](model)
	e_connection.On(event.ConnectionReady, func() {
		schema.syncIndexes(reflect.TypeOf(sample).Elem())
	})
	return model
}

func (m Model[T]) Create(doc T, ctx ...context.Context) T {
	document := enforceSchema(m.schema, &doc)
	qkit.Must(m.Collection().InsertOne(e_utils.DefaultCTX(ctx), document))
	return document
}

func (m Model[T]) InsertMany(docs []T, ctx ...context.Context) []T {
	var documents []interface{}
	for _, doc := range docs {
		documents = append(documents, enforceSchema(m.schema, &doc))
	}
	qkit.Must(m.Collection().InsertMany(e_utils.DefaultCTX(ctx), documents))
	return e_utils.CastSlice[T](documents)
}

func (m Model[T]) Find(query ...primitive.M) Model[T] {
	m.pipeline = append(m.pipeline, bson.D{{Key: "$match", Value: e_utils.DefaultQuery(query...)}})
	m.executor = func(ctx context.Context) any {
		var results []T
		e_utils.Must(qkit.Must(m.Collection().Aggregate(ctx, m.pipeline)).All(ctx, &results))
		return results
	}
	return m
}

func (m Model[T]) FindOne(query ...primitive.M) Model[T] {
	m.pipeline = append(m.pipeline, 
		bson.D{{Key: "$match", Value: e_utils.DefaultQuery(query...)}},
		bson.D{{Key: "$limit", Value: 1}},
	)
	m.executor = func(ctx context.Context) any {
		var results []T
		e_utils.Must(qkit.Must(m.Collection().Aggregate(ctx, m.pipeline)).All(ctx, &results))
		return results[0]
	}
	return m
}
func (m Model[T]) CountDocuments(query ...primitive.M) Model[T] {
	m.pipeline = append(m.pipeline, bson.D{{Key: "$count", Value: "count"}})
	m.executor = func(ctx context.Context) any {
		var results []map[string]any
		e_utils.Must(qkit.Must(m.Collection().Aggregate(ctx, m.pipeline)).All(ctx, &results))
		return int64(qkit.Cast[int32](results[0]["count"]))
	}
	m.pipeline = append(m.pipeline, bson.D{{Key: "$match", Value: e_utils.DefaultQuery(query...)}})
	return m
}

func (m Model[T]) Exec(ctx ...context.Context) any {
	return m.executor(e_utils.DefaultCTX(ctx))
}

func (m Model[T]) Collection() *mongo.Collection {
	return e_connection.Use(m.schema.Options.Database, m.schema.Options.Connection).Collection(m.schema.Options.Collection)
}

func (m Model[T]) CreateCollection(ctx ...context.Context) *mongo.Collection {
	e_connection.Use(m.schema.Options.Database, m.schema.Options.Connection).CreateCollection(e_utils.DefaultCTX(ctx), m.schema.Options.Collection)
	return m.Collection()
}

func (m Model[T]) Validate() error {
	return nil
}
