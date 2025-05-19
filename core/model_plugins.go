package elemental

import (
	"github.com/elcengine/elemental/plugins/filterquery"
)

// QS allows you to construct an Elemental query directly from a request's query string.
//
// It uses the filterquery plugin to parse the query string and apply filters, sorting, lookups, and projections to the final query.
//
// Usage:
//
//	UserModel.QS("filter[name]=John&sort[name]=asc&include=field1&select=field1").ExecTT()
func (m Model[T]) QS(query string) Model[T] {
	return m.QSR(fq.Parse(query))
}

// QSR allows you to construct an Elemental query directly from a FilterQueryResult.
//
// Usage:
//
//	UserModel.QSR(fq.Parse("filter[name]=John&sort[name]=asc&include=field1&select=field1")).ExecTT()
func (m Model[T]) QSR(result fq.Result) Model[T] {
	if len(result.Filters) > 0 {
		m = m.Find(result.Filters)
	}
	if len(result.Include) > 0 {
		for _, field := range result.Include {
			m = m.Populate(field)
		}
	}
	if len(result.SecondaryFilters) > 0 {
		m = m.Find(result.SecondaryFilters)
	}
	if len(result.Sorts) > 0 {
		m = m.Sort(result.Sorts)
	}
	if len(result.Select) > 0 {
		m = m.Select(result.Select)
	}
	if result.Page > 0 && result.Limit > 0 {
		m = m.Paginate(result.Page, result.Limit)
	}
	return m
}
