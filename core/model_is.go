package elemental

import (
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

func (m Model[T]) IsType(value bsontype.Type) Model[T] {
	return m.addToPipeline("$match", "$type", value)
}

