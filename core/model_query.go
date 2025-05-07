package elemental

import (
	"context"
	"errors"
	"github.com/elcengine/elemental/utils"

	"github.com/samber/lo"
	"github.com/spf13/cast"
)

func (m Model[T]) Where(field string, equals ...any) Model[T] {
	m.whereField = field
	if len(equals) > 0 {
		m = m.Equals(equals[0])
	}
	return m
}

func (m Model[T]) OrWhere(field string, equals ...any) Model[T] {
	m.whereField = field
	m.orConditionActive = true
	if len(equals) > 0 {
		m = m.Equals(equals[0])
	}
	return m
}

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
			m.executor(m, e_utils.DefaultCTX(ctx))
		})
		if err != nil {
			panic(errors.New("failed to schedule query"))
		}
		cron.Start()
		return cast.ToInt(id)
	}
	return m.executor(m, e_utils.DefaultCTX(ctx))
}

// ExecT is a convenience method that executes the query and returns the first result.
// It is a type safe method, so you don't need to cast the result. If the query returns nothing
// it will return the zero value of the type.
func (m Model[T]) ExecT(ctx ...context.Context) T {
	result := m.Exec(ctx...)
	return e_utils.Cast[T](result)
}

// ExecP is a convenience method that executes the query and returns the first result as a pointer.
// It is a type safe method, so you don't need to cast the result. If the query returns nothing
// it will return nil.
func (m Model[T]) ExecP(ctx ...context.Context) *T {
	result := m.Exec(ctx...)
	return e_utils.Cast[*T](result)
}

// ExecSlice is a convenience method that executes the query and returns the results as a slice.
// It is a type safe method, so you don't need to cast the result. If the query returns nothing
// it will return an empty slice.
func (m Model[T]) ExecSlice(ctx ...context.Context) []T {
	result := m.Exec(ctx...)
	return e_utils.Cast[[]T](result)
}
