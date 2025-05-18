package elemental

import (
	"context"

	"github.com/elcengine/elemental/utils"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Returns the underlying collection instance this model uses.
func (m Model[T]) Collection() *mongo.Collection {
	connection := lo.FromPtr(lo.CoalesceOrEmpty(m.temporaryConnection, &m.Schema.Options.Connection))
	database := lo.FromPtr(lo.CoalesceOrEmpty(m.temporaryDatabase, &m.Schema.Options.Database))
	collection := lo.FromPtr(lo.CoalesceOrEmpty(m.temporaryCollection, &m.Schema.Options.Collection))
	return UseDatabase(database, connection).Collection(collection)
}

// Returns the underlying client instance this model uses
func (m Model[T]) Connection() *mongo.Client {
	return GetConnection(lo.FromPtr(lo.CoalesceOrEmpty(m.temporaryConnection, &m.Schema.Options.Connection)))
}

// Returns the underlying client instance this model uses
// This method is an alias for Connection() and is kept for clarity to indicate that this retrieves a client instance
func (m Model[T]) Client() *mongo.Client {
	return m.Connection()
}

// Returns the underlying database instance this model uses
func (m Model[T]) Database() *mongo.Database {
	return m.Collection().Database()
}

// Returns the count of all documents in a collection or view
func (m Model[T]) EstimatedDocumentCount(ctx ...context.Context) int64 {
	count, _ := m.Collection().EstimatedDocumentCount(utils.CtxOrDefault(ctx))
	return count
}

// Returns statistics about the model collection
func (m Model[T]) Stats(ctx ...context.Context) CollectionStats {
	result := m.Database().RunCommand(utils.CtxOrDefault(ctx), bson.M{"collStats": m.Schema.Options.Collection})
	var stats CollectionStats
	lo.Must0(result.Decode(&stats))
	return stats
}

// The total amount of storage in bytes allocated to this collection for document storage
func (m Model[T]) StorageSize(ctx ...context.Context) int64 {
	return m.Stats(utils.CtxOrDefault(ctx)).StorageSize
}

// The total size in bytes of the data in the collection plus the size of every index on the collection
func (m Model[T]) TotalSize(ctx ...context.Context) int64 {
	return m.Stats(utils.CtxOrDefault(ctx)).Size
}

// The total size of all indexes for the collection
func (m Model[T]) TotalIndexSize(ctx ...context.Context) int64 {
	return m.Stats(utils.CtxOrDefault(ctx)).TotalIndexSize
}

// The average size of each document in the collection
func (m Model[T]) AvgObjSize(ctx ...context.Context) int64 {
	return m.Stats(utils.CtxOrDefault(ctx)).AvgObjSize
}

// Returns true if the model collection is a capped collection, otherwise returns false
func (m Model[T]) IsCapped(ctx ...context.Context) bool {
	return m.Stats(utils.CtxOrDefault(ctx)).Capped
}

// Returns the number of indexes in the model collection
func (m Model[T]) NumberOfIndexes(ctx ...context.Context) int64 {
	return m.Stats(utils.CtxOrDefault(ctx)).NumberOfIndexes
}
