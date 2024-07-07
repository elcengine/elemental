package filter_query

import (
	"strings"
)

type FilterQueryResult struct {
	Filters     map[string]interface{}
	Sorts       map[string]interface{}
	Include     []string
	Select      []string
	Prepaginate bool
}

func Parse(queryString string) FilterQueryResult {
	result := FilterQueryResult{}
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
			result.Select = append(result.Select, strings.Split(value, ",")...)
		}
		if strings.Contains(key, "prepaginate") {
			result.Prepaginate = value == "true"
		}
	}
	return result
}
