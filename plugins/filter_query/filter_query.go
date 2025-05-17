package fq

import (
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
)

type FilterQueryResult struct {
	Filters          bson.M   // Primary filters which are usually evaluated as the first stage of an aggregation query
	SecondaryFilters bson.M   // Secondary filters which are usually evaluated after any lookups
	Sorts            bson.M   // Fields to sort by, with 1 or 'asc' for ascending and -1 or 'desc' for descending
	Include          []string // Fields to populate/lookup in the result set
	Select           bson.M   // Fields to select in the result set
	Prepaginate      bool     // Whether to paginate the results before any lookups
	Page             int      // The page number for pagination
	Limit            int      // The page size for pagination
}

// Parses the given query string into a Elemental FilterQueryResult.
func Parse(queryString string) FilterQueryResult {
	result := FilterQueryResult{}
	result.Filters = bson.M{}
	result.SecondaryFilters = bson.M{}
	result.Sorts = bson.M{}
	result.Select = bson.M{}
	queries := strings.Split(queryString, "&")
	for _, query := range queries {
		if query == "" {
			continue
		}
		pair := strings.Split(query, "=")
		key := pair[0]
		value := pair[1]
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
		if strings.Contains(key, "include") {
			result.Include = append(result.Include, strings.Split(value, ",")...)
		}
		if strings.Contains(key, "select") {
			for _, field := range strings.Split(value, ",") {
				if strings.HasPrefix(field, "-") {
					result.Select[field[1:]] = 0
				} else {
					result.Select[field] = 1
				}
			}
		}
		if strings.Contains(key, "prepaginate") {
			result.Prepaginate = value == "true"
		}
		if strings.Contains(key, "page") {
			result.Page = cast.ToInt(value)
			if result.Page < 0 {
				result.Page = 0
			}
		}
		if strings.Contains(key, "limit") {
			result.Limit = cast.ToInt(value)
			if result.Limit < 0 {
				result.Limit = 0
			}
		}
	}
	result.Filters = mapFilters(result.Filters)
	result.SecondaryFilters = mapFilters(result.SecondaryFilters)
	return result
}
