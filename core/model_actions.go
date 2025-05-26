package elemental

import (
	"context"
	"github.com/elcengine/elemental/utils"
	"github.com/samber/lo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Creates the collection used by this model. This method will only create the collection if it does not exist.
// This will happen automatically when the model is created, so you most likely won't need to call this method manually.
func (m Model[T]) CreateCollection(ctx ...context.Context) *mongo.Collection {
	UseDatabase(m.Schema.Options.Database, m.Schema.Options.Connection).
		CreateCollection(utils.CtxOrDefault(ctx), m.Schema.Options.Collection, &m.Schema.Options.CollectionOptions)
	return m.Collection()
}

// Drops the collection used by this model.
func (m Model[T]) Drop(ctx ...context.Context) {
	lo.Must0(m.Collection().Drop(utils.CtxOrDefault(ctx)))
}

// Sends out a ping to the underlying client connection used by this model.
func (m Model[T]) Ping(ctx ...context.Context) error {
	return m.Database().Client().Ping(utils.CtxOrDefault(ctx), nil)
}

// Creates or updates the indexes for this model. This method will only create the indexes if they do not exist.
func (m Model[T]) SyncIndexes(ctx ...context.Context) {
	m.Schema.syncIndexes(m.docReflectType, lo.FromPtr(m.temporaryDatabase), lo.FromPtr(m.temporaryConnection), lo.FromPtr(m.temporaryCollection), ctx...)
}

// Drops all indexes for this model except the default `_id` index.
func (m Model[T]) DropIndexes(ctx ...context.Context) (bson.Raw, error) {
	return m.Collection().Indexes().DropAll(utils.CtxOrDefault(ctx))
}

// Drops a specific index for this model.
func (m Model[T]) DropIndex(indexName string, ctx ...context.Context) (bson.Raw, error) {
	return m.Collection().Indexes().DropOne(utils.CtxOrDefault(ctx), indexName)
}

// Validates a document against the model schema. This method will panic if any errors are found.
// This is the method being called when a new document is inserted.
func (m Model[T]) Validate(doc T) {
	enforceSchema(m.Schema, &doc, nil, false)
}

// Sets a temporary connection for this model. This connection will be used for the next operation only.
func (m Model[T]) SetConnection(connection string) Model[T] {
	m.temporaryConnection = &connection
	return m
}

// Sets a temporary database for this model. This database will be used for the next operation only.
func (m Model[T]) SetDatabase(database string) Model[T] {
	m.temporaryDatabase = &database
	return m
}

// Sets a temporary collection for this model. This collection will be used for the next operation only.
func (m Model[T]) SetCollection(collection string) Model[T] {
	m.temporaryCollection = &collection
	return m
}
