package elemental

type middlewareFunc func(...interface{}) bool

type listener[T any] struct {
	functions []middlewareFunc
}

type pre[T any] struct {
	save listener[T]
}

type post[T any] struct {
	save listener[T]
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
