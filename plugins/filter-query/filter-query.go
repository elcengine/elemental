package filter_query

import (
	"strings"
)

type FilterQueryResult struct {
	Filters     map[string]interface{}
	Sorts       map[string]interface{}
	Include     []string
	Select      map[string]interface{}
	Prepaginate bool
}

func Parse(queryString string) FilterQueryResult {
	result := FilterQueryResult{}
	result.Filters = make(map[string]interface{})
	result.Sorts = make(map[string]interface{})
	result.Select = make(map[string]interface{})
	queries := strings.Split(queryString, "&")
	for _, query := range queries {
		if query == "" {
			continue
		}
		pair := strings.Split(query, "=")
		key := pair[0]
		value := pair[1]
		if strings.Contains(key, "filter") {
			filterKey := extractFieldName(key)
			result.Filters[filterKey] = value
		}
		if strings.Contains(key, "sort") {
			sortKey := extractFieldName(key)
			if value == "asc" || value == "1" {
				result.Sorts[sortKey] = 1
			} else {
				result.Sorts[sortKey] = -1
			}
		}
		if strings.Contains(key, "include") {
			result.Include = append(result.Include, strings.Split(value, ",")...)
		}
		if strings.Contains(key, "select") {
			for _, field := range strings.Split(value, ",") {
				if (strings.HasPrefix(field, "-")) {
					result.Select[field[1:]] = 0
				} else {
					result.Select[field] = 1
				}
			}
		}
		if strings.Contains(key, "prepaginate") {
			result.Prepaginate = value == "true"
		}
	}
	return result
}
