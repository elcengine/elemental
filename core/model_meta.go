package elemental

import (
	"context"
	"elemental/connection"
	"elemental/utils"

	"go.mongodb.org/mongo-driver/mongo"
)


func (m Model[T]) Collection() *mongo.Collection {
	return e_connection.Use(m.schema.Options.Database, m.schema.Options.Connection).Collection(m.schema.Options.Collection)
}

func (m Model[T]) EstimatedDocumentCount(ctx ...context.Context) int64 {
	count, _ := m.Collection().EstimatedDocumentCount(e_utils.DefaultCTX(ctx))
	return count
}