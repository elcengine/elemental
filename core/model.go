package elemental

import (
	"context"
	"elemental/connection"
	"time"

	"github.com/clubpay/qlubkit-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ModelSkeleton[T any] interface {
	Schema() Schema
	Create() primitive.ObjectID
	FindOne(query primitive.M) *T
}

type Model[T any] struct {
	schema Schema
}

type Document[T any] struct {
	Base      T                   `json:",omitempty" bson:",omitempty"`
	ID        *primitive.ObjectID `json:"_id" bson:"_id"`
	CreatedAt *time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt *time.Time          `json:"updated_at" bson:"updated_at"`
}

var models = make(map[string]Model[any])

func NewModel[T any](name string, schema Schema) Model[T] {
	if _, ok := models[name]; ok {
		return qkit.Cast[Model[T]](models[name])
	}
	model := Model[T]{schema: schema}
	models[name] = qkit.Cast[Model[any]](model)
	return model
}

func (m Model[T]) Create(doc T) primitive.ObjectID {
	document := enforceSchema(m.schema, &Document[T]{
		Base:      doc,
		CreatedAt: qkit.ValPtr(time.Now()),
	})
	result := qkit.Must(m.Collection().InsertOne(context.TODO(), document))
	return result.InsertedID.(primitive.ObjectID)
}

func (m Model[T]) FindOne(query primitive.M) *T {
	return nil
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
