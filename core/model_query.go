package elemental

import (
	"context"
	"elemental/utils"
	"errors"

	"github.com/samber/lo"
)

func (m Model[T]) Where(field string) Model[T] {
	m.whereField = field
	return m
}

func (m Model[T]) OrWhere(field string) Model[T] {
	m.whereField = field
	m.orConditionActive = true
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

func (m Model[T]) Exec(ctx ...context.Context) any {
	if m.executor == nil {
		m.executor = func(ctx context.Context) any {
			var results []T
			e_utils.Must(lo.Must(m.Collection().Aggregate(ctx, m.pipeline)).All(ctx, &results))
			m.checkConditionsAndPanic(results)
			return results
		}
	}
	return m.executor(e_utils.DefaultCTX(ctx))
}
