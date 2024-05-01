package elemental

import (
	e_utils "elemental/utils"

	"go.mongodb.org/mongo-driver/mongo"
)

type middlewareFunc func(...interface{}) bool

type listener[T any] struct {
	functions []middlewareFunc
}

type pre[T any] struct {
	save listener[T]
	updateOne listener[T]
}

type post[T any] struct {
	save listener[T]
	updateOne listener[T]
}

type middleware[T any] struct {
	pre  pre[T]
	post post[T]
}

func newMiddleware[T any]() middleware[T] {
	return middleware[T]{}
}

func (l listener[T]) run(args ...interface{}) {
	for _, middleware := range l.functions {
		next := middleware(args...)
		if !next {
			break
		}
	}
}

func (m Model[T]) PreSave(f func(doc T) bool) {
	m.middleware.pre.save.functions = append(m.middleware.pre.save.functions, func(args ...interface{}) bool {
		return f(args[0].(T))
	})
}

func (m Model[T]) PostSave(f func(doc T) bool) {
	m.middleware.post.save.functions = append(m.middleware.post.save.functions, func(args ...interface{}) bool {
		return f(args[0].(T))
	})
}

func (m Model[T]) PreUpdateOne(f func(doc any) bool) {
	m.middleware.pre.updateOne.functions = append(m.middleware.pre.updateOne.functions, func(args ...interface{}) bool {
		return f(args[0])
	})
}

func (m Model[T]) PostUpdateOne(f func(result *mongo.UpdateResult, err error) bool) {
	m.middleware.post.updateOne.functions = append(m.middleware.post.updateOne.functions, func(args ...interface{}) bool {
		return f(args[0].(*mongo.UpdateResult), e_utils.Cast[error](args[1]))
	})
}