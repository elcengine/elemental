package elemental

import (
	"context"
	"reflect"
	"strings"

	"github.com/elcengine/elemental/utils"
	"github.com/spf13/cast"

	"github.com/gertd/go-pluralize"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ModelInterface[T any] interface {
	Exec(ctx ...context.Context) any
	Connection() *mongo.Client
}

type Model[T any] struct {
	Name                string // The name of the model
	Schema              Schema // The underlying schema used by this model
	Cloned              bool   // Indicates if this model has been cloned at least once
	pipeline            mongo.Pipeline
	executor            func(m Model[T], ctx context.Context) any
	whereField          string
	failWith            *error
	orConditionActive   bool
	upsert              bool
	returnNew           bool
	middleware          *middleware[T]
	temporaryConnection *string
	temporaryDatabase   *string
	temporaryCollection *string
	schedule            *string
	onScheduleExecError *func(any)
	softDeleteEnabled   bool
	deletedAtFieldName  string
	triggerExit         chan bool
}

var pluralizeClient = pluralize.NewClient()

// Static map of all created models. You can use this map to access models by name.
// The key is the name of the model and the value is the model itself.
var Models = make(map[string]any)

// A Predefined model with type map[string]any which can be used to access the Elemental APIs if you don't want to define a schema.
var NativeModel = NewModel[map[string]any]("ElementalNativeModel", NewSchema(map[string]Field{}))

// NewModel creates a new model with the given name and schema. If a model with the same name already exists, it will return the existing model.
func NewModel[T any](name string, schema Schema) Model[T] {
	if _, ok := Models[name]; ok {
		return utils.Cast[Model[T]](Models[name])
	}
	schema.Options.Collection = lo.CoalesceOrEmpty(schema.Options.Collection, pluralizeClient.Plural(strings.ToLower(name)))
	middleware := newMiddleware[T]()
	model := Model[T]{
		Name:        name,
		Schema:      schema,
		middleware:  &middleware,
		triggerExit: make(chan bool, 1),
	}
	Models[name] = model
	onConnectionComplete := func() {
		model.CreateCollection()
		model.SyncIndexes()
		if model.Schema.Options.Auditing {
			model.EnableAuditing()
		}
	}
	if model.Ping() != nil {
		OnConnectionEvent(EventDeploymentDiscovered, onConnectionComplete)
	} else {
		onConnectionComplete()
	}
	return model
}

// Extends the query to insert a single document into the collection.
// This method validates the document against the model schema and panics if any errors are found.
func (m Model[T]) Create(doc T) Model[T] {
	m.executor = func(m Model[T], ctx context.Context) any {
		documentToInsert, detailedDocument := enforceSchema(m.Schema, &doc, nil)
		detailedDocumentEntity := utils.CastBSON[T](detailedDocument)
		m.middleware.pre.save.run(detailedDocumentEntity)
		lo.Must(m.Collection().InsertOne(ctx, documentToInsert))
		m.middleware.post.save.run(detailedDocumentEntity)
		return detailedDocumentEntity
	}
	return m
}

// Extends the query to insert multiple documents into the collection.
// This method is an alias for InsertMany and is provided for convenience.
func (m Model[T]) CreateMany(docs []T) Model[T] {
	return m.InsertMany(docs)
}

// Extends the query to insert multiple documents into the collection.
// This method validates the document against the model schema and panics if any errors are found.
func (m Model[T]) InsertMany(docs []T) Model[T] {
	m.executor = func(m Model[T], ctx context.Context) any {
		var documentsToInsert, detailedDocuments []any
		for _, doc := range docs {
			documentToInsert, detailedDocument := enforceSchema(m.Schema, &doc, nil)
			documentsToInsert = append(documentsToInsert, documentToInsert)
			detailedDocuments = append(detailedDocuments, detailedDocument)
		}
		lo.Must(m.Collection().InsertMany(ctx, documentsToInsert))
		return utils.CastBSONSlice[T](detailedDocuments)
	}
	return m
}

// Extends the query with a match stage to find multiple documents in the collection.
// Optionally accepts one or more queries to filter the results.
// If multiple queries are provided, they are merged into a single from left to right.
func (m Model[T]) Find(query ...primitive.M) Model[T] {
	m.executor = func(m Model[T], ctx context.Context) any {
		var results []T
		cursor := lo.Must(m.Collection().Aggregate(ctx, m.pipeline))
		m.checkConditionsAndPanicForErr(cursor.All(ctx, &results))
		m.middleware.post.find.run(results)
		m.checkConditionsAndPanic(results)
		return results
	}
	q := utils.MergedQueryOrDefault(query)
	if m.softDeleteEnabled {
		q[m.deletedAtFieldName] = primitive.M{"$exists": false}
	}
	m.pipeline = append(m.pipeline, bson.D{{Key: "$match", Value: q}})
	return m
}

// Extends the query with a limit stage to find a single document.
// It optionally accepts one or more queries to filter the results before limiting.
// If multiple queries are provided, they are merged into a single from left to right.
func (m Model[T]) FindOne(query ...primitive.M) Model[T] {
	q := utils.MergedQueryOrDefault(query)
	if m.softDeleteEnabled {
		q[m.deletedAtFieldName] = primitive.M{"$exists": false}
	}
	m.pipeline = append(m.pipeline,
		bson.D{{Key: "$match", Value: q}},
		bson.D{{Key: "$limit", Value: 1}},
	)
	m.executor = func(m Model[T], ctx context.Context) any {
		var results []T
		cursor := lo.Must(m.Collection().Aggregate(ctx, m.pipeline))
		m.checkConditionsAndPanicForErr(cursor.All(ctx, &results))
		m.checkConditionsAndPanic(results)
		if len(results) == 0 {
			return nil
		}
		return results[0]
	}
	return m
}

// Extends the query with a match stage to find a document by its ID.
// The id can be a string or an ObjectID.
func (m Model[T]) FindByID(id any) Model[T] {
	q := primitive.M{"_id": utils.EnsureObjectID(id)}
	if m.softDeleteEnabled {
		q[m.deletedAtFieldName] = primitive.M{"$exists": false}
	}
	return m.FindOne(q)
}

// Extends the query with a count stage to count the number of documents in the collection.
// It optionally accepts one or more queries to filter the results before counting.
// If multiple queries are provided, they are merged into a single from left to right.
func (m Model[T]) CountDocuments(query ...primitive.M) Model[T] {
	q := utils.MergedQueryOrDefault(query)
	if m.softDeleteEnabled {
		q[m.deletedAtFieldName] = primitive.M{"$exists": false}
	}
	m.pipeline = append(m.pipeline, bson.D{{Key: "$match", Value: q}}, bson.D{{Key: "$count", Value: "count"}})
	m.executor = func(m Model[T], ctx context.Context) any {
		var results []map[string]any
		cursor := lo.Must(m.Collection().Aggregate(ctx, m.pipeline))
		m.checkConditionsAndPanicForErr(cursor.All(ctx, &results))
		if len(results) == 0 {
			return 0
		}
		return cast.ToInt64(results[0]["count"])
	}
	return m
}

// Distinct returns a list of distinct values for the given field.
// It optionally accepts one or more queries to filter the results before getting the distinct values.
// If multiple queries are provided, they are merged into a single from left to right.
func (m Model[T]) Distinct(field string, query ...primitive.M) Model[T] {
	q := utils.MergedQueryOrDefault(query)
	if m.softDeleteEnabled {
		q[m.deletedAtFieldName] = primitive.M{"$exists": false}
	}
	m.pipeline = append(m.pipeline, bson.D{{Key: "$match", Value: q}}, bson.D{{Key: "$group", Value: primitive.M{"_id": "$" + field}}})
	m.executor = func(m Model[T], ctx context.Context) any {
		var results []map[string]any
		cursor := lo.Must(m.Collection().Aggregate(ctx, m.pipeline))
		m.checkConditionsAndPanicForErr(cursor.All(ctx, &results))
		var distinct = make([]string, 0, len(results))
		for _, result := range results {
			distinct = append(distinct, utils.Cast[string](result["_id"]))
		}
		return distinct
	}
	return m
}

// Extends the query with a sort stage.
// The sort stage is used to specify the order in which the results should be returned.
func (m Model[T]) Sort(args ...any) Model[T] {
	if len(args) == 1 {
		for field, order := range utils.Cast[primitive.M](args[0]) {
			m = m.addToPipeline("$sort", field, order)
		}
	} else {
		if (len(args) % 2) != 0 {
			panic(ErrMustPairSortArguments)
		}
		for i := 0; i < len(args); i += 2 {
			m = m.addToPipeline("$sort", utils.Cast[string](args[i]), args[i+1])
		}
	}
	return m
}

// Extends the query with a projection stage.
// The projection stage is used to specify which fields to include or exclude from the results.
func (m Model[T]) Select(fields ...any) Model[T] {
	inputType := reflect.TypeOf(fields[0]).Kind()
	if inputType == reflect.Map {
		m.pipeline = append(m.pipeline, bson.D{{Key: "$project", Value: fields[0]}})
		return m
	}
	var selection []string
	switch {
	case inputType == reflect.Slice:
		selection = fields[0].([]string)
	case len(fields) == 1 && inputType == reflect.String:
		selection = strings.FieldsFunc(fields[0].(string), func(r rune) bool {
			return r == ',' || r == ' '
		})
	case len(fields) > 1:
		selection = cast.ToStringSlice(fields)
	}
	for _, field := range selection {
		if strings.HasPrefix(field, "-") {
			m = m.addToPipeline("$project", field[1:], 0)
		} else {
			m = m.addToPipeline("$project", field, 1)
		}
	}
	return m
}

// Creates a clone of the current model with the same query pipeline and options as has been set on the current model.
// Invoking this method will set the Cloned flag of the newly created model to true.
func (m Model[T]) Clone() Model[T] {
	return Model[T]{
		Name:                m.Name,
		Schema:              m.Schema,
		Cloned:              true,
		pipeline:            m.pipeline,
		executor:            m.executor,
		whereField:          m.whereField,
		failWith:            m.failWith,
		orConditionActive:   m.orConditionActive,
		upsert:              m.upsert,
		returnNew:           m.returnNew,
		middleware:          m.middleware,
		temporaryConnection: m.temporaryConnection,
		temporaryDatabase:   m.temporaryDatabase,
		temporaryCollection: m.temporaryCollection,
		schedule:            m.schedule,
		softDeleteEnabled:   m.softDeleteEnabled,
		deletedAtFieldName:  m.deletedAtFieldName,
	}
}
