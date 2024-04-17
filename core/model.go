package elemental

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ModelSkeleton[T any] interface {
	Schema() Schema
	Create() primitive.ObjectID
	FindOne(query primitive.M) *T
}

type Model[T any] struct {
	ModelSkeleton[T]
	CreatedAt primitive.DateTime
	UpdatedAt primitive.DateTime
}

func (u Model[T]) Create() primitive.ObjectID {
	enforceSchema(u.Schema(), u)
	return primitive.ObjectID{}
}

func (u Model[T]) FindOne(query primitive.M) *T {
	return nil
}

func (u Model[T]) Validate() error {
	return nil
}

func (u Model[T]) ValidateField() error {
	return nil
}