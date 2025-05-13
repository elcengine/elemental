package elemental

import (
	"context"
	"reflect"
	"time"

	e_utils "github.com/elcengine/elemental/utils"
	"github.com/samber/lo"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (m Model[T]) FindOneAndDelete(query ...primitive.M) Model[T] {
	q := e_utils.MergedQueryOrDefault(query)
	if m.softDeleteEnabled {
		m = m.UpdateOne(&q, m.softDeletePayload())
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

func (m Model[T]) FindByIdAndDelete(id primitive.ObjectID) Model[T] {
	if m.softDeleteEnabled {
		return m.FindOneAndUpdate(lo.ToPtr(primitive.M{"_id": id}), m.softDeletePayload())
	} else {
		return m.FindOneAndDelete(primitive.M{"_id": id})
	}
}

func (m Model[T]) DeleteOne(query ...primitive.M) Model[T] {
	q := e_utils.MergedQueryOrDefault(query)
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

func (m Model[T]) DeleteByID(id primitive.ObjectID) Model[T] {
	if m.softDeleteEnabled {
		return m.UpdateOne(lo.ToPtr(primitive.M{"_id": id}), m.softDeletePayload())
	} else {
		return m.DeleteOne(primitive.M{"_id": id})
	}
}

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

func (m Model[T]) DeleteMany(query ...primitive.M) Model[T] {
	q := e_utils.MergedQueryOrDefault(query)
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

func (m Model[T]) EnableSoftDelete() Model[T] {
	m.softDeleteEnabled = true
	return m
}

func (m Model[T]) softDeletePayload() primitive.M {
	m.deletedAtFieldName = "deleted_at"
	return primitive.M{m.deletedAtFieldName: time.Now().Format(time.RFC3339)}
}
