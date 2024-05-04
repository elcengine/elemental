package queryfilter

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

var complexOperators = []string{"and", "or"}

func ReplaceOperator(value, operator string) string {
    return strings.TrimSuffix(strings.Replace(value, operator+"(", "", 1), ")")
}
 
func ParseOperatorValue(value, operator string) interface{} {
	value = ReplaceOperator(value, operator)
	if _, err := strconv.ParseFloat(value, 64); err != nil {
		if _, err := time.Parse(time.RFC3339, value); err == nil {
			return value
		} else if matched, _ := regexp.MatchString("^[0-9a-fA-F]{24}$", value); matched {
			return struct{ ID string }{value}
		}
	} else {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return value
}

func MapValue(value string) interface{} {
	if strings.HasPrefix(value, "eq(") {
		value = ParseOperatorValue(value, "eq").(string)
		if value == "true" || value == "false" {
			parsedvalue, _ := strconv.ParseBool(value)
			return parsedvalue
		}
		return value
	} else if strings.HasPrefix(value, "ne(") {
		return map[string]interface{}{"$ne": ParseOperatorValue(value, "ne")}
	} else if strings.HasPrefix(value, "gt(") {
		return map[string]interface{}{"$gt": ParseOperatorValue(value, "gt")}
	} else if strings.HasPrefix(value, "gte(") {
		return map[string]interface{}{"$gte": ParseOperatorValue(value, "gte")}
	} else if strings.HasPrefix(value, "lt(") {
		return map[string]interface{}{"$lt": ParseOperatorValue(value, "lt")}
	} else if strings.HasPrefix(value, "lte(") {
		return map[string]interface{}{"$lte": ParseOperatorValue(value, "lte")}
	} else if strings.HasPrefix(value, "in(") {
		return map[string]interface{}{"$in": strings.Split(ParseOperatorValue(value, "in").(string), ",")}
	} else if strings.HasPrefix(value, "nin(") {
		return map[string]interface{}{"$nin": strings.Split(ParseOperatorValue(value, "nin").(string), ",")}
	} else if strings.HasPrefix(value, "reg(") {
		return map[string]interface{}{"$regex": regexp.MustCompile(ReplaceOperator(value, "reg"))}
	} else if strings.HasPrefix(value, "exists(") {
		b, _ := strconv.ParseBool(ParseOperatorValue(value, "exists").(string))
		return map[string]interface{}{"$exists": b}
	}
	if value == "true" || value == "false" {
		parsedvalue, _ := strconv.ParseBool(value)
		return parsedvalue
	}
	return value
}

func MapFilters(filter map[string]string) map[string]interface{} {
	result := make(map[string]interface{})
	if filter != nil {
		for key, value := range filter {
			if contains(complexOperators, key) {
				subFilters := make([]map[string]interface{}, 0)
				kvPairs := strings.Split(value, ",")
				for _, kv := range kvPairs {
					subKey, subValue := strings.Split(kv, "=")[0], strings.Split(kv, "=")[1]
					subFilters = append(subFilters, map[string]interface{}{subKey: MapValue(subValue)})
				}
				result["$"+key] = subFilters
			}  else {
				for _, op := range complexOperators {
					if strings.HasPrefix(value, op+"(") {
						values := strings.Split(ParseOperatorValue(value, op).(string), ",")
						subFilters := make([]map[string]interface{}, 0)
						for _, subValue := range values {
							subFilters = append(subFilters, map[string]interface{}{key: MapValue(subValue)})
						}
						result["$"+op] = subFilters
					} else {
					result[key] = filter[key]}
				}
			}
		}
	}
	return result
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
