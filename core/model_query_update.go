package elemental

import (
	"context"
	"reflect"

	"maps"

	"github.com/elcengine/elemental/utils"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Extends the query with an update operation matching the given query(s)
// If multiple queries are provided, they are merged into a single from left to right.
// It updates only the first document that matches the query.
func (m Model[T]) FindOneAndUpdate(query *primitive.M, doc any, opts ...*options.FindOneAndUpdateOptions) Model[T] {
	m.executor = func(m Model[T], ctx context.Context) any {
		m.middleware.pre.findOneAndUpdate.run(doc)
		return (func() any {
			var resultDoc T
			filters := lo.FromPtr(query)
			maps.Copy(filters, m.findMatchStage())
			result := m.Collection().FindOneAndUpdate(ctx, filters, primitive.M{"$set": m.parseDocument(doc)}, parseUpdateOptions(m, opts)...)
			m.middleware.post.findOneAndUpdate.run(&resultDoc)
			m.checkConditionsAndPanicForSingleResult(result)
			lo.Must0(result.Decode(&resultDoc))
			return resultDoc
		})()
	}
	return m
}

// Extends the query with an update operation matching the given id
// The id can be a string or an ObjectID.
func (m Model[T]) FindByIDAndUpdate(id any, doc any, opts ...*options.FindOneAndUpdateOptions) Model[T] {
	m.executor = func(m Model[T], ctx context.Context) any {
		var resultDoc T
		result := m.Collection().FindOneAndUpdate(ctx, primitive.M{"_id": e_utils.EnsureObjectID(id)}, primitive.M{"$set": m.parseDocument(doc)}, parseUpdateOptions(m, opts)...)
		m.checkConditionsAndPanicForSingleResult(result)
		lo.Must0(result.Decode(&resultDoc))
		return resultDoc
	}
	return m
}

// Extends the query with an update operation matching the given query(s)
// If multiple queries are provided, they are merged into a single from left to right.
// It updates only the first document that matches the query.
func (m Model[T]) UpdateOne(query *primitive.M, doc any, opts ...*options.UpdateOptions) Model[T] {
	m.executor = func(m Model[T], ctx context.Context) any {
		filters := make(primitive.M)
		if query != nil {
			filters = lo.FromPtr(query)
		}
		maps.Copy(filters, m.findMatchStage())
		if m.upsert {
			if len(opts) == 0 {
				opts = append(opts, &options.UpdateOptions{Upsert: lo.ToPtr(true)})
			} else {
				opts[0].SetUpsert(true)
			}
		}
		m.middleware.pre.updateOne.run(doc)
		result, err := m.Collection().UpdateOne(ctx, filters,
			primitive.M{"$set": m.parseDocument(doc)}, opts...)
		m.middleware.post.updateOne.run(result, err)
		m.checkConditionsAndPanicForErr(err)
		return result
	}
	return m
}

// Extends the query with an update operation matching the given id
// It updates only the first document that matches the id.
// The id can be a string or an ObjectID.
func (m Model[T]) UpdateByID(id any, doc any, opts ...*options.UpdateOptions) Model[T] {
	m.executor = func(m Model[T], ctx context.Context) any {
		result, err := m.Collection().UpdateOne(ctx, primitive.M{"_id": e_utils.EnsureObjectID(id)},
			primitive.M{"$set": m.parseDocument(doc)}, parseUpdateOptions(m, opts)...)
		m.checkConditionsAndPanicForErr(err)
		return result
	}
	return m
}

// Extends the query with an update operation matching the id of the given document
func (m Model[T]) Save(doc T) Model[T] {
	m.executor = func(m Model[T], ctx context.Context) any {
		m.middleware.pre.save.run(doc)
		m.UpdateByID(reflect.ValueOf(doc).FieldByName("ID").Interface().(primitive.ObjectID), doc).Exec(ctx) //nolint:contextcheck
		m.middleware.post.save.run(doc)
		return doc
	}
	return m
}

// Extends the query with an update operation matching the given query(s)
// If multiple queries are provided, they are merged into a single from left to right.
// It updates all documents that match the query.
func (m Model[T]) UpdateMany(query *primitive.M, doc any, opts ...*options.UpdateOptions) Model[T] {
	m.executor = func(m Model[T], ctx context.Context) any {
		filters := make(primitive.M)
		if query != nil {
			filters = lo.FromPtr(query)
		}
		maps.Copy(filters, m.findMatchStage())
		result, err := m.Collection().UpdateMany(ctx, filters, primitive.M{"$set": m.parseDocument(doc)}, parseUpdateOptions(m, opts)...)
		m.checkConditionsAndPanicForErr(err)
		return result
	}
	return m
}

// Extends the query with a replace operation matching the given query(s)
// If multiple queries are provided, they are merged into a single from left to right.
// It replaces only the first document that matches the query.
func (m Model[T]) ReplaceOne(query *primitive.M, doc any, opts ...*options.ReplaceOptions) Model[T] {
	m.executor = func(m Model[T], ctx context.Context) any {
		filters := make(primitive.M)
		if query != nil {
			filters = lo.FromPtr(query)
		}
		maps.Copy(filters, m.findMatchStage())
		result, err := m.Collection().ReplaceOne(ctx, filters, m.parseDocument(doc),
			parseUpdateOptions(m, opts)...)
		m.checkConditionsAndPanicForErr(err)
		return result
	}
	return m
}

// Extends the query with a replace operation matching the given id
// The id can be a string or an ObjectID.
func (m Model[T]) ReplaceByID(id any, doc any, opts ...*options.ReplaceOptions) Model[T] {
	m.executor = func(m Model[T], ctx context.Context) any {
		result, err := m.Collection().ReplaceOne(ctx, primitive.M{"_id": e_utils.EnsureObjectID(id)},
			m.parseDocument(doc), parseUpdateOptions(m, opts)...)
		m.checkConditionsAndPanicForErr(err)
		return result
	}
	return m
}

// Extends the query with a replace operation matching the given query(s)
// If multiple queries are provided, they are merged into a single from left to right.
// It replaces only the first document that matches the query.
func (m Model[T]) FindOneAndReplace(query *primitive.M, doc any, opts ...*options.FindOneAndReplaceOptions) Model[T] {
	m.executor = func(m Model[T], ctx context.Context) any {
		var resultDoc T
		filters := make(primitive.M)
		if query != nil {
			filters = lo.FromPtr(query)
		}
		maps.Copy(filters, m.findMatchStage())
		res := m.Collection().FindOneAndReplace(ctx, filters, m.parseDocument(doc), opts...)
		m.middleware.post.findOneAndReplace.run(&resultDoc)
		m.checkConditionsAndPanicForSingleResult(res)
		lo.Must0(res.Decode(&resultDoc))
		return resultDoc
	}
	return m
}

// Extends the query with a replace operation matching the given id
// This method will return the replaced document.
// The id can be a string or an ObjectID.
func (m Model[T]) FindByIDAndReplace(id any, doc any, opts ...*options.FindOneAndReplaceOptions) Model[T] {
	m.executor = func(m Model[T], ctx context.Context) any {
		var resultDoc T
		res := m.Collection().FindOneAndReplace(ctx, primitive.M{"_id": e_utils.EnsureObjectID(id)},
			m.parseDocument(doc), parseUpdateOptions(m, opts)...)
		m.checkConditionsAndPanicForSingleResult(res)
		lo.Must0(res.Decode(&resultDoc))
		return resultDoc
	}
	return m
}

func (m Model[T]) Set(doc any) Model[T] {
	return m.setUpdateOperator("$set", doc)
}

func (m Model[T]) Unset(doc any) Model[T] {
	if reflect.TypeOf(doc).Kind() == reflect.String {
		doc = primitive.M{doc.(string): ""}
	}
	return m.setUpdateOperator("$unset", doc)
}

// Extends the query with an increment operation matching the given field
func (m Model[T]) Inc(field string, value int) Model[T] {
	return m.setUpdateOperator("$inc", primitive.M{field: value})
}

// Extends the query with a decrement operation matching the given field
func (m Model[T]) Dec(field string, value int) Model[T] {
	return m.setUpdateOperator("$inc", primitive.M{field: -value})
}

// Extends the query with a multiplication operation matching the given field
func (m Model[T]) Mul(field string, value int) Model[T] {
	return m.setUpdateOperator("$mul", primitive.M{field: value})
}

// Extends the query with a division operation matching the given field
func (m Model[T]) Div(field string, value int) Model[T] {
	return m.setUpdateOperator("$mul", primitive.M{field: (float64(1) / float64(value))})
}

// Extends the query to update the name of the given field
func (m Model[T]) Rename(field string, newField string) Model[T] {
	return m.setUpdateOperator("$rename", primitive.M{field: newField})
}

// Extends the query to update the value of the given field if the given value is less than the current value
func (m Model[T]) Min(field string, value int) Model[T] {
	return m.setUpdateOperator("$min", primitive.M{field: value})
}

// Extends the query to update the value of the given field if the given value is greater than the current value
func (m Model[T]) Max(field string, value int) Model[T] {
	return m.setUpdateOperator("$max", primitive.M{field: value})
}

// Extends the query to set the value of the given field to the current date
func (m Model[T]) CurrentDate(field string) Model[T] {
	return m.setUpdateOperator("$currentDate", primitive.M{field: true})
}

// Extends the query to add the given values to the set of values for the given field if they are not already present
func (m Model[T]) AddToSet(field string, values ...any) Model[T] {
	if len(values) == 1 {
		return m.setUpdateOperator("$addToSet", primitive.M{field: values[0]})
	}
	return m.setUpdateOperator("$addToSet", primitive.M{field: primitive.M{"$each": values}})
}

// Extends the query to remove the last element from the array of the given field
func (m Model[T]) Pop(field string, value ...int) Model[T] {
	if len(value) == 0 {
		return m.setUpdateOperator("$pop", primitive.M{field: 1})
	}
	return m.setUpdateOperator("$pop", primitive.M{field: value[0]})
}

// Extends the query to remove the first element from the array of the given field
func (m Model[T]) Shift(field string) Model[T] {
	return m.setUpdateOperator("$pop", primitive.M{field: -1})
}

// Extends the query to remove all elements from the array of the given field where the value is equal to the given value
func (m Model[T]) Pull(field string, value any) Model[T] {
	return m.setUpdateOperator("$pull", primitive.M{field: value})
}

// Extends the query to remove all elements from the array of the given field where the value is equal to any of the given values
func (m Model[T]) PullAll(field string, values ...any) Model[T] {
	return m.setUpdateOperator("$pullAll", primitive.M{field: values})
}

// Extends the query to add the given values to the array of the given field
func (m Model[T]) Push(field string, values ...any) Model[T] {
	if len(values) == 1 {
		return m.setUpdateOperator("$push", primitive.M{field: values[0]})
	}
	return m.setUpdateOperator("$push", primitive.M{field: primitive.M{"$each": values}})
}

// Signals the query to insert a new document if no documents match the query
func (m Model[T]) Upsert() Model[T] {
	m.upsert = true
	return m
}

// Signals the query to return the new document instead of the original document after an update
func (m Model[T]) New() Model[T] {
	m.returnNew = true
	return m
}
