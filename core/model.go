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
	returnSingleRecord bool
	resultExtractor    func(docs []map[string]any) any
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
	m.resultExtractor = func(docs []map[string]any) any {
		return e_utils.CastSliceFromMaps[T](docs)
	}
	return m
}

func (m Model[T]) FindOne(query ...primitive.M) Model[T] {
	m.returnSingleRecord = true
	m.resultExtractor = func(docs []map[string]any) any {
		if len(docs) == 0 {
			return nil
		}
		return qkit.CastJSON[T](docs[0])
	}
	m.pipeline = append(m.pipeline, bson.D{{Key: "$match", Value: e_utils.DefaultQuery(query...)}})
	return m
}
func (m Model[T]) CountDocuments(query ...primitive.M) Model[T] {
	m.pipeline = append(m.pipeline, bson.D{{Key: "$count", Value: "count"}})
	m.resultExtractor = func(docs []map[string]any) any {
		if len(docs) == 0 {
			return 0
		}
		return int64(qkit.Cast[int32](docs[0]["count"]))
	}
	m.pipeline = append(m.pipeline, bson.D{{Key: "$match", Value: e_utils.DefaultQuery(query...)}})
	return m
}

func (m Model[T]) Exec(ctx ...context.Context) any {
	cursor := qkit.Must(m.Collection().Aggregate(e_utils.DefaultCTX(ctx), m.pipeline))
	var results []map[string]any
	if err := cursor.All(e_utils.DefaultCTX(ctx), &results); err != nil {
		panic(err)
	}
	if m.resultExtractor != nil {
		return m.resultExtractor(results)
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

func (m Model[T]) CreateCollection(ctx ...context.Context) *mongo.Collection {
	e_connection.Use(m.schema.Options.Database, m.schema.Options.Connection).CreateCollection(e_utils.DefaultCTX(ctx), m.schema.Options.Collection)
	return m.Collection()
}
