package elemental

import (
	"context"
	"elemental/utils"
	"reflect"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (m Model[T]) FindOneAndDelete(query ...primitive.M) Model[T] {
	m.executor = func(ctx context.Context) any {
		var doc T
		result := m.Collection().FindOneAndDelete(ctx, e_utils.DefaultQuery(query...))
		if result.Err() != nil {
			if m.failWith != nil {
				panic(*m.failWith)
			}
			panic(result.Err())
		}
		e_utils.Must(result.Decode(&doc))
		return doc
	}
	return m
}

func (m Model[T]) FindByIdAndDelete(id primitive.ObjectID) Model[T] {
	return m.FindOneAndDelete(primitive.M{"_id": id})
}

func (m Model[T]) DeleteOne(query ...primitive.M) Model[T] {
	m.executor = func(ctx context.Context) any {
		_, err := m.Collection().DeleteOne(ctx, e_utils.DefaultQuery(query...))
		if err != nil {
			if m.failWith != nil {
				panic(*m.failWith)
			}
			panic(err)
		}
		return nil
	}
	return m
}

func (m Model[T]) DeleteById(id primitive.ObjectID) Model[T] {
	return m.DeleteOne(primitive.M{"_id": id})
}

func (m Model[T]) Delete(doc T) {
	m.DeleteById(reflect.ValueOf(doc).FieldByName("ID").Interface().(primitive.ObjectID)).Exec()
}

func (m Model[T]) DeleteMany(query ...primitive.M) Model[T] {
	m.executor = func(ctx context.Context) any {
		result, err := m.Collection().DeleteMany(ctx, e_utils.DefaultQuery(query...))
		if err != nil {
			if m.failWith != nil {
				panic(*m.failWith)
			}
			panic(err)
		}
		return result
	}
	return m
}
