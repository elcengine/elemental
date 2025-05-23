package fq

import (
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
)

type Result struct {
	Filters          bson.M   // Primary filters which are usually evaluated as the first stage of an aggregation query
	SecondaryFilters bson.M   // Secondary filters which are usually evaluated after any lookups
	Sorts            bson.M   // Fields to sort by, with 1 or 'asc' for ascending and -1 or 'desc' for descending
	Include          []string // Fields to populate/lookup in the result set
	Select           bson.M   // Fields to select in the result set
	Prepaginate      bool     // Whether to paginate the results before any lookups
	Page             int64    // The page number for pagination
	Limit            int64    // The page size for pagination
}

// Type alias for the Result struct
//
// Deprecated: Use Result instead.
type FilterQueryResult = Result

// Parses the given query string into a Elemental FilterQueryResult.
func Parse(qs string) Result {
	result := Result{}
	result.Filters = bson.M{}
	result.SecondaryFilters = bson.M{}
	result.Sorts = bson.M{}
	result.Select = bson.M{}
	queries := strings.Split(qs, "&")
	for _, query := range queries {
		if query == "" {
			continue
		}
		pair := strings.Split(query, "=")
		key := pair[0]
		value := strings.Join(pair[1:], "=")
		if strings.Contains(key, "filter") {
			if filterKey := extractFieldName(key); filterKey != "" {
				result.Filters[filterKey] = value
			}
		}
		if strings.Contains(key, "secondaryFilter") {
			if filterKey := extractFieldName(key); filterKey != "" {
				result.SecondaryFilters[filterKey] = value
			}
		}
		if strings.Contains(key, "sort") {
			if sortKey := extractFieldName(key); sortKey != "" {
				if value == "asc" || value == "1" {
					result.Sorts[sortKey] = 1
				} else {
					result.Sorts[sortKey] = -1
				}
			}
		}
		if key == "include" {
			result.Include = append(result.Include, strings.Split(value, ",")...)
		}
		if key == "select" {
			for _, field := range strings.Split(value, ",") {
				if strings.HasPrefix(field, "-") {
					result.Select[field[1:]] = 0
				} else {
					result.Select[field] = 1
				}
			}
		}
		if key == "prepaginate" {
			result.Prepaginate = value == "true"
		}
		if key == "page" {
			result.Page = max(cast.ToInt64(value), 0)
		}
		if key == "limit" {
			result.Limit = max(cast.ToInt64(value), 0)
		}
	}
	result.Filters = mapFilters(result.Filters)
	result.SecondaryFilters = mapFilters(result.SecondaryFilters)
	return result
}
