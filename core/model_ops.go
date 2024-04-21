package elemental

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (m Model[T]) Equals(value any) Model[T] {
	return m.addToPipeline("$match", "$eq", value)
}

func (m Model[T]) NotEquals(value any) Model[T] {
	return m.addToPipeline("$match", "$ne", value)
}

func (m Model[T]) LessThan(value any) Model[T] {
	return m.addToPipeline("$match", "$lt", value)
}

func (m Model[T]) GreaterThan(value any) Model[T] {
	return m.addToPipeline("$match", "$gt", value)
}

func (m Model[T]) LessThanOrEquals(value any) Model[T] {
	return m.addToPipeline("$match", "$lte", value)
}

func (m Model[T]) GreaterThanOrEquals(value any) Model[T] {
	return m.addToPipeline("$match", "$gte", value)
}

func (m Model[T]) Between(min, max any) Model[T] {
	return m.addToPipeline("$match", "$gte", min).addToPipeline("$match", "$lte", max)
}

func (m Model[T]) Exists(value bool) Model[T] {
	return m.addToPipeline("$match", "$exists", value)
}

func (m Model[T]) In(values ...any) Model[T] {
	return m.addToPipeline("$match", "$in", values)
}

func (m Model[T]) NotIn(values ...any) Model[T] {
	return m.addToPipeline("$match", "$nin", values)
}

func (m Model[T]) ElementMatches(query primitive.M) Model[T] {
	return m.addToPipeline("$match", "$elemMatch", query)
}

func (m Model[T]) Size(value int) Model[T] {
	return m.addToPipeline("$match", "$size", value)
}