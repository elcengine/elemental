package elemental

import (
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Extends the query to return only documents that equal the given value.
func (m Model[T]) Equals(value any) Model[T] {
	return m.addToFilters("$eq", value)
}

// Extends the query to return only documents that do not equal the given value.
func (m Model[T]) NotEquals(value any) Model[T] {
	return m.addToFilters("$ne", value)
}

// Extends the query to return only documents that are less than the given value.
func (m Model[T]) LessThan(value any) Model[T] {
	return m.addToFilters("$lt", value)
}

// Extends the query to return only documents that are greater than the given value.
func (m Model[T]) GreaterThan(value any) Model[T] {
	return m.addToFilters("$gt", value)
}

// Extends the query to return only documents that are less than or equal to the given value.
func (m Model[T]) LessThanOrEquals(value any) Model[T] {
	return m.addToFilters("$lte", value)
}

// Extends the query to return only documents that are greater than or equal to the given value.
func (m Model[T]) GreaterThanOrEquals(value any) Model[T] {
	return m.addToFilters("$gte", value)
}

// Extends the query to return only documents that are between the given values.
func (m Model[T]) Between(minimum, maximum any) Model[T] {
	return m.addToFilters("$gte", minimum).addToFilters("$lte", maximum)
}

// Extends the query to return only documents where the value of a field divided by a divisor is equal to the given remainder.
func (m Model[T]) Mod(divisor, remainder int) Model[T] {
	return m.addToFilters("$mod", []int{divisor, remainder})
}

// Extends the query to return only documents that match the given regular expression.
// It optionaly accepts options for the regex, such as "i" for case-insensitive matching.
func (m Model[T]) Regex(pattern string, options ...string) Model[T] {
	return m.addToFilters("$regex", primitive.Regex{Pattern: pattern, Options: lo.FirstOrEmpty(options)})
}

// Extends the query to return only documents where the given field exists.
func (m Model[T]) Exists(value bool) Model[T] {
	return m.addToFilters("$exists", value)
}

// Extends the query to return only documents where the given field has a value equal to any of the given values.
func (m Model[T]) In(values ...any) Model[T] {
	return m.addToFilters("$in", values)
}

// Extends the query to return only documents where the given field does not have any of the given values.
func (m Model[T]) NotIn(values ...any) Model[T] {
	return m.addToFilters("$nin", values)
}

// Extends the query to return only documents where the given field is an array and has an element that matches the given query.
func (m Model[T]) ElementMatches(query primitive.M) Model[T] {
	return m.addToFilters("$elemMatch", query)
}

// Extends the query to return only documents where the given field is an array and has an element that matches the given value.
// This is a shorthand for ElementMatches with an equality operator.
func (m Model[T]) Has(value string) Model[T] {
	query := primitive.M{"$eq": value}
	return m.addToFilters("$elemMatch", query)
}

// Extends the query to return only documents where the given field is an array and has a size equal to the given value.
func (m Model[T]) Size(value int) Model[T] {
	return m.addToFilters("$size", value)
}

// Extends the query with an "or" condition.
func (m Model[T]) Or() Model[T] {
	m.orConditionActive = true
	return m
}
