package elemental

import (
	"context"
	"reflect"
	"strings"

	e_connection "github.com/elcengine/elemental/connection"
	e_constants "github.com/elcengine/elemental/constants"
	e_utils "github.com/elcengine/elemental/utils"

	"github.com/gertd/go-pluralize"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
)

type Model[T any] struct {
	Name                string
	Schema              Schema
	pipeline            mongo.Pipeline
	executor            func(m Model[T], ctx context.Context) any
	whereField          string
	failWith            *error
	orConditionActive   bool
	upsert              bool
	returnNew           bool
	middleware          *middleware[T]
	clusterOps             *ClusterOp[T]
	temporaryConnection *string
	temporaryDatabase   *string
	temporaryCollection *string
	schedule            *string
	softDeleteEnabled   bool
	deletedAtFieldName  string
}

var Models = make(map[string]any)

// A Predefined model with type map[string]any which can be used to access the Elemental APIs if you don't want to define a schema.
var NativeModel = NewModel[map[string]any]("ElementalNativeModel", NewSchema(map[string]Field{}))

func NewModel[T any](name string, schema Schema) Model[T] {
	if _, ok := Models[name]; ok {
		return e_utils.Cast[Model[T]](Models[name])
	}
	schema.Options.Collection, _ = lo.Coalesce(schema.Options.Collection, pluralize.NewClient().Plural(strings.ToLower(name)))
	middleware := newMiddleware[T]()
	model := Model[T]{
		Name:       name,
		Schema:     schema,
		middleware: &middleware,
	}
	Models[name] = model
	connectionReady := func() {
		model.CreateCollection()
		model.SyncIndexes()
		if model.Schema.Options.Auditing {
			model.EnableAuditing()
		}
	}
	if model.Ping() != nil {
		e_connection.On(event.ConnectionReady, connectionReady)
	} else {
		connectionReady()
	}
	return model
}

func (m Model[T]) Create(doc T) Model[T] {
	m.executor = func(m Model[T], ctx context.Context) any {
		documentToInsert, detailedDocument := enforceSchema(m.Schema, &doc, nil)
		detailedDocumentEntity := e_utils.CastBSON[T](detailedDocument)
		m.middleware.pre.save.run(detailedDocumentEntity)
		_, err := m.Collection().InsertOne(ctx, documentToInsert)
		if err != nil {
			panic(err)
		}
		m.middleware.post.save.run(detailedDocumentEntity)
		return detailedDocumentEntity
	}
	return m
}

func (m Model[T]) CreateMany(docs []T) Model[T] {
	return m.InsertMany(docs)
}

func (m Model[T]) InsertMany(docs []T) Model[T] {
	m.executor = func(m Model[T], ctx context.Context) any {
		var documentsToInsert, detailedDocuments []interface{}
		for _, doc := range docs {
			documentToInsert, detailedDocument := enforceSchema(m.Schema, &doc, nil)
			documentsToInsert = append(documentsToInsert, documentToInsert)
			detailedDocuments = append(detailedDocuments, detailedDocument)
		}
		_, err := m.Collection().InsertMany(ctx, documentsToInsert)
		if err != nil {
			panic(err)
		}
		return e_utils.CastBSONSlice[T](detailedDocuments)
	}
	return m
}

func (m Model[T]) Find(query ...primitive.M) Model[T] {
	m.executor = func(m Model[T], ctx context.Context) any {
		var results []T
		cursor, err := m.Collection().Aggregate(ctx, m.pipeline)
		if err != nil {
			panic(err)
		}
		err = cursor.All(ctx, &results)
		if err != nil {
			panic(err)
		}
		m.middleware.post.find.run(results)
		m.checkConditionsAndPanic(results)
		return results
	}
	q := e_utils.DefaultQuery(query...)
	if m.softDeleteEnabled {
		q[m.deletedAtFieldName] = primitive.M{"$exists": false}
	}
	m.pipeline = append(m.pipeline, bson.D{{Key: "$match", Value: q}})
	return m
}

func (m Model[T]) FindOne(query ...primitive.M) Model[T] {
	q := e_utils.DefaultQuery(query...)
	if m.softDeleteEnabled {
		q[m.deletedAtFieldName] = primitive.M{"$exists": false}
	}
	m.pipeline = append(m.pipeline,
		bson.D{{Key: "$match", Value: q}},
		bson.D{{Key: "$limit", Value: 1}},
	)
	m.executor = func(m Model[T], ctx context.Context) any {
		var results []T
		cursor, err := m.Collection().Aggregate(ctx, m.pipeline)
		if err != nil {
			panic(err)
		}
		err = cursor.All(ctx, &results)
		if err != nil {
			panic(err)
		}
		m.checkConditionsAndPanic(results)
		if len(results) == 0 {
			return nil
		}
		return results[0]
	}
	return m
}

func (m Model[T]) FindByID(id primitive.ObjectID) Model[T] {
	q := primitive.M{"_id": id}
	if m.softDeleteEnabled {
		q[m.deletedAtFieldName] = primitive.M{"$exists": false}
	}
	return m.FindOne(q)
}

func (m Model[T]) CountDocuments(query ...primitive.M) Model[T] {
	q := e_utils.DefaultQuery(query...)
	if m.softDeleteEnabled {
		q[m.deletedAtFieldName] = primitive.M{"$exists": false}
	}
	m.pipeline = append(m.pipeline, bson.D{{Key: "$match", Value: q}}, bson.D{{Key: "$count", Value: "count"}})
	m.executor = func(m Model[T], ctx context.Context) any {
		var results []map[string]any
		cursor, err := m.Collection().Aggregate(ctx, m.pipeline)
		if err != nil {
			panic(err)
		}
		err = cursor.All(ctx, &results)
		if err != nil {
			panic(err)
		}
		return int64(e_utils.Cast[int32](results[0]["count"]))
	}
	return m
}

func (m Model[T]) Distinct(field string, query ...primitive.M) Model[T] {
	q := e_utils.DefaultQuery(query...)
	if m.softDeleteEnabled {
		q[m.deletedAtFieldName] = primitive.M{"$exists": false}
	}
	m.pipeline = append(m.pipeline, bson.D{{Key: "$match", Value: q}}, bson.D{{Key: "$group", Value: primitive.M{"_id": "$" + field}}})
	m.executor = func(m Model[T], ctx context.Context) any {
		var results []map[string]any
		cursor, err := m.Collection().Aggregate(ctx, m.pipeline)
		if err != nil {
			panic(err)
		}
		err = cursor.All(ctx, &results)
		if err != nil {
			panic(err)
		}
		var distinct []string
		for _, result := range results {
			distinct = append(distinct, e_utils.Cast[string](result["_id"]))
		}
		return distinct
	}
	return m
}

func (m Model[T]) Sort(args ...any) Model[T] {
	if len(args) == 1 {
		for field, order := range e_utils.Cast[primitive.M](args[0]) {
			m = m.addToPipeline("$sort", field, order)
		}
	} else {
		if (len(args) % 2) != 0 {
			panic(e_constants.ErrMustPairSortArguments)
		}
		for i := 0; i < len(args); i += 2 {
			m = m.addToPipeline("$sort", e_utils.Cast[string](args[i]), args[i+1])
		}
	}
	return m
}

func (m Model[T]) Select(fields ...any) Model[T] {
	var selection []string
	if len(fields) == 1 && reflect.TypeOf(fields[0]).Kind() == reflect.String {
		selection = strings.FieldsFunc(fields[0].(string), func(r rune) bool {
			return r == ',' || r == ' '
		})
	} else if len(fields) > 1 {
		selection = e_utils.CastSlice[string](fields)
	} else if reflect.TypeOf(fields[0]).Kind() == reflect.Slice {
		selection = fields[0].([]string)
	}
	if len(selection) > 0 {
		for _, field := range selection {
			if strings.HasPrefix(field, "-") {
				m = m.addToPipeline("$project", field[1:], 0)
			} else {
				m = m.addToPipeline("$project", field, 1)
			}
		}
	} else if reflect.TypeOf(fields[0]).Kind() == reflect.Map {
		for field, value := range e_utils.Cast[primitive.M](fields[0]) {
			m = m.addToPipeline("$project", field, value)
		}
	}
	return m
}

// This feature is still experimental and not fully implemented.
func (m Model[T]) UseCluster(connection *string) ClusterOp[T] {
	return Cluster(&m, connection)
}
