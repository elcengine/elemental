// Deprecated: This package is not being maintained anymore, use at your own risk.
package e_repository

import (
	"context"
	"errors"
	"log"

	elemental "github.com/elcengine/elemental/core"
	"github.com/elcengine/elemental/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// A lightweight repository powered by Elemental.
// This has no overhead of a model or schema definition. It's just a simple interface for CRUD operations.
type Repository[T any] struct {
	collection string
}

// NewRepository creates a new repository for the given collection.
// The type parameter T is the type of the document in the collection.
func NewRepository[T any](collection string) Repository[T] {
	return Repository[T]{collection: collection}
}

func (r Repository[T]) Create(payload T) primitive.ObjectID {
	result, err := elemental.UseDefaultDatabase().Collection(r.collection).InsertOne(context.TODO(), payload)
	if err != nil {
		panic(err)
	}
	return result.InsertedID.(primitive.ObjectID)
}

func (r Repository[T]) FindOne(query primitive.M) *T {
	model := new(T)
	doc := elemental.UseDefaultDatabase().Collection(r.collection).FindOne(context.Background(), query)
	if doc.Err() != nil {
		if errors.Is(doc.Err(), mongo.ErrNoDocuments) {
			log.Fatalf("%v %s", r, doc.Err().Error())
			return nil
		}
		panic(doc.Err())
	}
	doc.Decode(&model)
	return model
}

func (r Repository[T]) FindByID(id any) *T {
	return r.FindOne(primitive.M{"_id": utils.EnsureObjectID(id)})
}

func (r Repository[T]) FindAll() []T {
	var users []T
	cursor, err := elemental.UseDefaultDatabase().Collection(r.collection).Find(context.Background(), primitive.M{})
	if err != nil {
		panic(err)
	}
	cursor.All(context.Background(), &users)
	return users
}

func (r Repository[T]) Update(id any, payload T) {
	_, err := elemental.UseDefaultDatabase().Collection(r.collection).
		UpdateOne(context.Background(), primitive.M{"_id": utils.EnsureObjectID(id)}, primitive.M{"$set": payload})
	if err != nil {
		panic(err)
	}
}

func (r Repository[T]) Delete(id any) {
	_, err := elemental.UseDefaultDatabase().Collection(r.collection).
		DeleteOne(context.Background(), primitive.M{"_id": utils.EnsureObjectID(id)})
	if err != nil {
		panic(err)
	}
}
