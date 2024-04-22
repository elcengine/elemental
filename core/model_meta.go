package elemental

import (
	"context"
	"elemental/connection"
	"elemental/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (m Model[T]) Collection() *mongo.Collection {
	return e_connection.Use(m.schema.Options.Database, m.schema.Options.Connection).Collection(m.schema.Options.Collection)
}

func (m Model[T]) Database() *mongo.Database {
	return m.Collection().Database()
}

func (m Model[T]) EstimatedDocumentCount(ctx ...context.Context) int64 {
	count, _ := m.Collection().EstimatedDocumentCount(e_utils.DefaultCTX(ctx))
	return count
}

func (m Model[T]) Stats(ctx ...context.Context) CollectionStats {
	result := m.Database().RunCommand(e_utils.DefaultCTX(ctx), bson.M{"collStats": m.schema.Options.Collection})
	var stats CollectionStats
	e_utils.Must(result.Decode(&stats))
	return stats
}

func (m Model[T]) TotalSize(ctx ...context.Context) int64 {
	return m.Stats(e_utils.DefaultCTX(ctx)).Size
}

func (m Model[T]) StorageSize(ctx ...context.Context) int64 {
	return m.Stats(e_utils.DefaultCTX(ctx)).StorageSize
}

func (m Model[T]) TotalIndexSize(ctx ...context.Context) int64 {
	return m.Stats(e_utils.DefaultCTX(ctx)).TotalIndexSize
}

func (m Model[T]) AvgObjSize(ctx ...context.Context) int64 {
	return m.Stats(e_utils.DefaultCTX(ctx)).AvgObjSize
}