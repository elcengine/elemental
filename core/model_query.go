package elemental

import (
	"context"
	"errors"

	e_utils "github.com/elcengine/elemental/utils"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

// Extends the query with a where clause. The value of the clause if specified within this method itself
// will function as an equals clause. If you want to use a different operator, you can use the Where method
// and chain it with another operator method.
func (m Model[T]) Where(field string, equals ...any) Model[T] {
	m.whereField = field
	if len(equals) > 0 {
		m = m.Equals(equals[0])
	}
	return m
}

// Extends the query with an or where clause. The value of the clause if specified within this method itself
// will function as an equals clause. If you want to use a different operator, you can use the OrWhere method
// and chain it with another operator method.
func (m Model[T]) OrWhere(field string, equals ...any) Model[T] {
	m.whereField = field
	m.orConditionActive = true
	if len(equals) > 0 {
		m = m.Equals(equals[0])
	}
	return m
}

// Instructs a query to panic if no results are found matching the given query.
// It optionally accepts a custom error to panic with. If no error is provided, it will panic with a default error message.
func (m Model[T]) OrFail(err ...error) Model[T] {
	if len(err) > 0 {
		m.failWith = &err[0]
	} else {
		m.failWith = lo.ToPtr(errors.New("no results found matching the given query"))
	}
	return m
}

// Exec is the final step in the query builder chain. It executes the query and returns the results.
// The result of this method is not type safe, so you need to cast it to the expected type.
func (m Model[T]) Exec(ctx ...context.Context) any {
	if m.executor == nil {
		m.executor = func(m Model[T], ctx context.Context) any {
			var results []T
			cursor, err := m.Collection().Aggregate(ctx, m.pipeline)
			if err != nil {
				panic(err)
			}
			err = cursor.All(ctx, &results)
			if err != nil {
				panic(err)
			}
			m.checkConditionsAndPanic(results)
			return results
		}
	}
	if m.schedule != nil {
		id, err := cron.AddFunc(*m.schedule, func() {
			m.executor(m, e_utils.CtxOrDefault(ctx))
		})
		if err != nil {
			panic(err)
		}
		cron.Start()
		return cast.ToInt(id)
	}
	return m.executor(m, e_utils.CtxOrDefault(ctx))
}

// ExecT is a convenience method that executes the query and returns the first result.
// It is a type safe method, so you don't need to cast the result. If the query returns nothing
// it will return the zero value of the type.
func (m Model[T]) ExecT(ctx ...context.Context) T {
	result := m.Exec(ctx...)
	return e_utils.Cast[T](result)
}

// ExecPtr is a convenience method that executes the query and returns the first result as a pointer.
// It is a type safe method, so you don't need to cast the result. If the query returns nothing
// it will return nil.
func (m Model[T]) ExecPtr(ctx ...context.Context) *T {
	result := m.Exec(ctx...)
	if result == nil {
		return nil
	}
	if val, ok := result.(*T); ok {
		return val
	}
	return lo.ToPtr(e_utils.Cast[T](result))
}

// ExecTT is a convenience method that executes the query and returns the results as a slice.
// It is a type safe method, so you don't need to cast the result. If the query returns nothing
// it will return an empty slice.
func (m Model[T]) ExecTT(ctx ...context.Context) []T {
	result := m.Exec(ctx...)
	return e_utils.Cast[[]T](result)
}

// ExecTP is a convenience method that executes the query and returns the results as a PaginateResult.
// It is a type safe method, so you don't need to cast the result. If the query returns nothing, it will return an empty PaginateResult.
// This method is useful for pagination queries.
func (m Model[T]) ExecTP(ctx ...context.Context) PaginateResult[T] {
	result := m.Exec(ctx...)
	return e_utils.Cast[PaginateResult[T]](result)
}

// ExecInt is a convenience method that executes the query and returns the first result as an int.
// It is a type safe method, so you don't need to cast the result. If the query returns nothing
// it will return 0.
// This method is useful for queries that return a single integer value, such as count queries
// or schedule queries which return a cron entry ID.
func (m Model[T]) ExecInt(ctx ...context.Context) int {
	result := m.Exec(ctx...)
	return cast.ToInt(result)
}

// ExecStringSlice is a convenience method that executes the query and returns the first result as a slice of strings.
// It is a type safe method, so you don't need to cast the result. If the query returns nothing
// it will return an empty slice.
// This method is useful for queries that return an array of strings, such as distinct queries.
func (m Model[T]) ExecStringSlice(ctx ...context.Context) []string {
	result := m.Exec(ctx...)
	return cast.ToStringSlice(result)
}
