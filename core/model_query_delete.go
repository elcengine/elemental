package elemental

import (
	"context"
	"reflect"
	"time"

	"github.com/elcengine/elemental/utils"
	"github.com/samber/lo"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Extends the query with a delete operation matching the given query(s)
// If multiple queries are provided, they are merged into a single from left to right.
// It deletes only the first document that matches the query.
// This method will return the deleted document.
// If the model has soft delete enabled, it will update the document with a deleted_at field instead of deleting it.
func (m Model[T]) FindOneAndDelete(query ...primitive.M) Model[T] {
	q := utils.MergedQueryOrDefault(query)
	if m.softDeleteEnabled {
		m = m.FindOneAndUpdate(&q, m.softDeletePayload())
	} else {
		m.executor = func(m Model[T], ctx context.Context) any {
			var doc T
			m.middleware.pre.findOneAndDelete.run(q)
			result := m.Collection().FindOneAndDelete(ctx, q)
			m.middleware.post.findOneAndDelete.run(&doc)
			m.checkConditionsAndPanicForSingleResult(result)
			lo.Must0(result.Decode(&doc))
			return doc
		}
	}
	return m
}

// Extends the query with a delete operation matching the given id
// It deletes only the first document that matches the id.
// This method will return the deleted document.
// If the model has soft delete enabled, it will update the document with a deleted_at field instead of deleting it.
// The id can be a string or an ObjectID.
func (m Model[T]) FindByIDAndDelete(id any) Model[T] {
	if m.softDeleteEnabled {
		return m.FindOneAndUpdate(lo.ToPtr(primitive.M{"_id": utils.EnsureObjectID(id)}),
			m.softDeletePayload())
	} else {
		return m.FindOneAndDelete(primitive.M{"_id": utils.EnsureObjectID(id)})
	}
}

// Extends the query with a delete operation matching the given query(s).
// If multiple queries are provided, they are merged into a single from left to right.
// It deletes only the first document that matches the query.
// This method will not return the deleted document.
// If the model has soft delete enabled, it will update the document with a deleted_at field instead of deleting it.
func (m Model[T]) DeleteOne(query ...primitive.M) Model[T] {
	q := utils.MergedQueryOrDefault(query)
	if m.softDeleteEnabled {
		m = m.UpdateOne(&q, m.softDeletePayload())
	} else {
		m.executor = func(m Model[T], ctx context.Context) any {
			m.middleware.pre.deleteOne.run(q)
			result, err := m.Collection().DeleteOne(ctx, q)
			m.middleware.post.deleteOne.run(result, err)
			m.checkConditionsAndPanicForErr(err)
			return result
		}
	}
	return m
}

// Extends the query with a delete operation matching the given id.
// It deletes only the first document that matches the id.
// This method will not return the deleted document.
// If the model has soft delete enabled, it will update the document with a deleted_at field instead of deleting it.
// The id can be a string or an ObjectID.
func (m Model[T]) DeleteByID(id any) Model[T] {
	if m.softDeleteEnabled {
		return m.UpdateOne(lo.ToPtr(primitive.M{"_id": utils.EnsureObjectID(id)}), m.softDeletePayload())
	} else {
		return m.DeleteOne(primitive.M{"_id": utils.EnsureObjectID(id)})
	}
}

// Extends the query with a delete operation matching the given document.
// If the model has soft delete enabled, it will update the document with a deleted_at field instead of deleting it.
func (m Model[T]) Delete(doc T) Model[T] {
	if m.softDeleteEnabled {
		m.executor = func(m Model[T], ctx context.Context) any {
			m.UpdateByID(reflect.ValueOf(doc).FieldByName("ID").Interface().(primitive.ObjectID), m.softDeletePayload()).Exec(ctx) //nolint:contextcheck
			return nil
		}
	} else {
		m.executor = func(m Model[T], ctx context.Context) any {
			m.DeleteByID(reflect.ValueOf(doc).FieldByName("ID").Interface().(primitive.ObjectID)).Exec(ctx) //nolint:contextcheck
			return nil
		}
	}
	return m
}

// Extends the query with a delete operation matching the given query(s)
// If multiple queries are provided, they are merged into a single from left to right.
// It deletes all documents that match the query.
// This method will not return the deleted documents.
// If the model has soft delete enabled, it will update the documents with a deleted_at field instead of deleting them.
func (m Model[T]) DeleteMany(query ...primitive.M) Model[T] {
	q := utils.MergedQueryOrDefault(query)
	if m.softDeleteEnabled {
		m = m.UpdateMany(&q, m.softDeletePayload())
	} else {
		m.executor = func(m Model[T], ctx context.Context) any {
			m.middleware.pre.deleteMany.run(q)
			result, err := m.Collection().DeleteMany(ctx, q)
			m.checkConditionsAndPanicForErr(err)
			m.middleware.post.deleteMany.run(result, err)
			return result
		}
	}
	return m
}

// Enables soft delete for the model.
func (m *Model[T]) EnableSoftDelete() {
	m.deletedAtFieldName = "deleted_at"
	m.softDeleteEnabled = true
}

// Disables soft delete for the model.
func (m *Model[T]) DisableSoftDelete() {
	m.softDeleteEnabled = false
}

func (m Model[T]) softDeletePayload() primitive.M {
	return primitive.M{m.deletedAtFieldName: time.Now().Format(time.RFC3339)}
}
