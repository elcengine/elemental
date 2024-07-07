package filter_query

import "regexp"

func extractFieldName(input string) (string) {
	re := regexp.MustCompile(`sort\[(.+?)\]`)
	matches := re.FindStringSubmatch(input)
	if len(matches) > 1 {
		return matches[1]
	}
	panic("Invalid field name")
}