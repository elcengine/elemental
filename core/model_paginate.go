package elemental

import (
	"context"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Limits the number of documents returned by the query.
func (m Model[T]) Limit(limit int64) Model[T] {
	m.pipeline = append(m.pipeline, bson.D{{Key: "$limit", Value: limit}})
	return m
}

// Skips the first n documents in the query result.
func (m Model[T]) Skip(skip int64) Model[T] {
	for i, stage := range m.pipeline {
		if stage[0].Key == "$limit" {
			newPipeline := make([]bson.D, len(m.pipeline)+1)
			copy(newPipeline, m.pipeline[:i])
			newPipeline[i] = bson.D{{Key: "$skip", Value: skip}}
			copy(newPipeline[i+1:], m.pipeline[i:])
			m.pipeline = newPipeline
			return m
		}
	}
	m.pipeline = append(m.pipeline, bson.D{{Key: "$skip", Value: skip}})
	return m
}

// Paginate the results of the query.
// It adds a $facet stage to the pipeline to get both the documents and the total count.
// The final result of the query will be a PaginateResult[T] struct containing the documents,
// total count, current page, limit, total pages, and next/previous page information.
func (m Model[T]) Paginate(page, limit int64) Model[T] {
	m = m.Skip((page - 1) * limit).Limit(limit)
	m.pipeline = []bson.D{{{Key: "$facet", Value: primitive.M{
		"docs": m.pipeline,
		"count": []primitive.M{
			{"$count": "count"},
		},
	}}}}
	m.executor = func(m Model[T], ctx context.Context) any {
		var results []facetResult[T]
		cursor, err := m.Collection().Aggregate(ctx, m.pipeline)
		if err != nil {
			panic(err)
		}
		err = cursor.All(ctx, &results)
		if err != nil {
			panic(err)
		}
		totalDocs := results[0].Count[0]["count"]
		totalPages := (totalDocs + limit - 1) / limit
		var prevPage, nextPage *int64
		prevPage = lo.ToPtr(page - 1)
		if lo.FromPtr(prevPage) < 1 {
			prevPage = nil
		}
		nextPage = lo.ToPtr(page + 1)
		if lo.FromPtr(nextPage) > totalPages {
			nextPage = nil
		}
		return PaginateResult[T]{
			Docs:       results[0].Docs,
			TotalDocs:  totalDocs,
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
			NextPage:   nextPage,
			PrevPage:   prevPage,
			HasPrev:    page > 1,
			HasNext:    page < totalPages,
		}
	}
	return m
}
