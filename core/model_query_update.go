package elemental

import (
	"context"
	"elemental/utils"
	"reflect"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (m Model[T]) FindOneAndUpdate(query *primitive.M, doc any, opts ...*options.FindOneAndUpdateOptions) Model[T] {
	m.executor = func(m Model[T], ctx context.Context) any {
		return (func() any {
			var resultDoc T
			filters := lo.FromPtr(query)
			for k, v := range m.findMatchStage().Map() {
				filters[k] = v
			}
			result := m.Collection().FindOneAndUpdate(ctx, filters, primitive.M{"$set": m.parseDocument(doc)}, parseUpdateOptions(m, opts)...)
			m.checkConditionsAndPanicForSingleResult(result)
			e_utils.Must(result.Decode(&resultDoc))
			return resultDoc
		})()
	}
	return m
}

func (m Model[T]) FindByIDAndUpdate(id primitive.ObjectID, doc any, opts ...*options.FindOneAndUpdateOptions) Model[T] {
	m.executor = func(m Model[T], ctx context.Context) any {
		var resultDoc T
		result := m.Collection().FindOneAndUpdate(ctx, primitive.M{"_id": id}, primitive.M{"$set": m.parseDocument(doc)}, parseUpdateOptions(m, opts)...)
		m.checkConditionsAndPanicForSingleResult(result)
		e_utils.Must(result.Decode(&resultDoc))
		return resultDoc
	}
	return m
}

func (m Model[T]) UpdateOne(query *primitive.M, doc any, opts ...*options.UpdateOptions) Model[T] {
	m.executor = func(m Model[T], ctx context.Context) any {
		filters := lo.FromPtr(query)
		for k, v := range m.findMatchStage().Map() {
			filters[k] = v
		}
		if m.upsert {
			if len(opts) == 0 {
				opts = append(opts, &options.UpdateOptions{Upsert: lo.ToPtr(true)})
			} else {
				opts[0].SetUpsert(true)
			}
		}
		result, err := m.Collection().UpdateOne(ctx, filters, primitive.M{"$set": m.parseDocument(doc)}, opts...)
		m.checkConditionsAndPanicForErr(err)
		return result
	}
	return m
}

func (m Model[T]) UpdateByID(id primitive.ObjectID, doc any, opts ...*options.UpdateOptions) Model[T] {
	m.executor = func(m Model[T], ctx context.Context) any {
		result, err := m.Collection().UpdateOne(ctx, primitive.M{"_id": id}, primitive.M{"$set": m.parseDocument(doc)}, parseUpdateOptions(m, opts)...)
		m.checkConditionsAndPanicForErr(err)
		return result
	}
	return m
}

func (m Model[T]) Save(doc T) {
	m.UpdateByID(reflect.ValueOf(doc).FieldByName("ID").Interface().(primitive.ObjectID), doc).Exec()
}

func (m Model[T]) UpdateMany(query *primitive.M, doc any, opts ...*options.UpdateOptions) Model[T] {
	m.executor = func(m Model[T], ctx context.Context) any {
		filters := lo.FromPtr(query)
		for k, v := range m.findMatchStage().Map() {
			filters[k] = v
		}
		result, err := m.Collection().UpdateMany(ctx, filters, primitive.M{"$set": m.parseDocument(doc)}, parseUpdateOptions(m, opts)...)
		m.checkConditionsAndPanicForErr(err)
		return result
	}
	return m
}

func (m Model[T]) ReplaceOne(query *primitive.M, doc any, opts ...*options.ReplaceOptions) Model[T] {
	m.executor = func(m Model[T], ctx context.Context) any {
		filters := lo.FromPtr(query)
		for k, v := range m.findMatchStage().Map() {
			filters[k] = v
		}
		result, err := m.Collection().ReplaceOne(ctx, filters, m.parseDocument(doc), parseUpdateOptions(m, opts)...)
		m.checkConditionsAndPanicForErr(err)
		return result
	}
	return m
}

func (m Model[T]) ReplaceByID(id primitive.ObjectID, doc any, opts ...*options.ReplaceOptions) Model[T] {
	m.executor = func(m Model[T], ctx context.Context) any {
		result, err := m.Collection().ReplaceOne(ctx, primitive.M{"_id": id}, m.parseDocument(doc), parseUpdateOptions(m, opts)...)
		m.checkConditionsAndPanicForErr(err)
		return result
	}
	return m
}

func (m Model[T]) FindOneAndReplace(query *primitive.M, doc any, opts ...*options.FindOneAndReplaceOptions) Model[T] {
	m.executor = func(m Model[T], ctx context.Context) any {
		var resultDoc T
		filters := lo.FromPtr(query)
		for k, v := range m.findMatchStage().Map() {
			filters[k] = v
		}
		res := m.Collection().FindOneAndReplace(ctx, filters, m.parseDocument(doc), opts...)
		m.checkConditionsAndPanicForSingleResult(res)
		e_utils.Must(res.Decode(&resultDoc))
		return resultDoc
	}
	return m
}

func (m Model[T]) FindByIDAndReplace(id primitive.ObjectID, doc any, opts ...*options.FindOneAndReplaceOptions) Model[T] {
	m.executor = func(m Model[T], ctx context.Context) any {
		var resultDoc T
		res := m.Collection().FindOneAndReplace(ctx, primitive.M{"_id": id}, m.parseDocument(doc), parseUpdateOptions(m, opts)...)
		m.checkConditionsAndPanicForSingleResult(res)
		e_utils.Must(res.Decode(&resultDoc))
		return resultDoc
	}
	return m
}

// Insert a new document if no documents match the query
func (m Model[T]) Upsert() Model[T] {
	m.upsert = true
	return m
}

// Return the new document instead of the original document after an update
func (m Model[T]) New() Model[T] {
	m.returnNew = true
	return m
}
