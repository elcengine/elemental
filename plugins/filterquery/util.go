//nolint:goconst
package fq

import (
	"regexp"
	"slices"
	"strings"

	"github.com/samber/lo"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var complexOperators = []string{"and", "or"}

func extractFieldName(input string) string {
	re := regexp.MustCompile(`^[^\[]+\[(.+?)\]$`)
	matches := re.FindStringSubmatch(input)
	return lo.NthOrEmpty(matches, 1)
}

func replaceOperator(value string, operator string) string {
	return value[len(operator)+1 : len(value)-1]
}

func parseOperatorValue(value any, operator string) any {
	strVal := cast.ToString(value)
	if operator != "" {
		strVal = replaceOperator(strVal, operator)
	}
	if f, err := cast.ToFloat64E(strVal); err == nil {
		return f
	}
	if oid, err := primitive.ObjectIDFromHex(strVal); err == nil {
		return oid
	}
	if t, err := cast.ToTimeE(strVal); err == nil {
		return t
	}
	return strVal
}

func mapValue(value any) any {
	switch {
	case strings.HasPrefix(cast.ToString(value), "eq("):
		value = parseOperatorValue(value, "eq")
		if value == "true" || value == "false" {
			return bson.M{"$eq": value == "true"}
		}
		return bson.M{"$eq": value}
	case strings.HasPrefix(cast.ToString(value), "ne("):
		return bson.M{"$ne": parseOperatorValue(value, "ne")}
	case strings.HasPrefix(cast.ToString(value), "gt("):
		return bson.M{"$gt": parseOperatorValue(value, "gt")}
	case strings.HasPrefix(cast.ToString(value), "gte("):
		return bson.M{"$gte": parseOperatorValue(value, "gte")}
	case strings.HasPrefix(cast.ToString(value), "lt("):
		return bson.M{"$lt": parseOperatorValue(value, "lt")}
	case strings.HasPrefix(cast.ToString(value), "lte("):
		return bson.M{"$lte": parseOperatorValue(value, "lte")}
	case strings.HasPrefix(cast.ToString(value), "in("):
		return bson.M{"$in": lo.Map(strings.Split(replaceOperator(cast.ToString(value), "in"), ","), func(value string, index int) any {
			return parseOperatorValue(value, "")
		})}
	case strings.HasPrefix(cast.ToString(value), "nin("):
		return bson.M{"$nin": lo.Map(strings.Split(replaceOperator(cast.ToString(value), "nin"), ","), func(value string, index int) any {
			return parseOperatorValue(value, "")
		})}
	case strings.HasPrefix(cast.ToString(value), "reg("):
		result := strings.Split(replaceOperator(cast.ToString(value), "reg"), "...[")
		regex := primitive.Regex{Pattern: result[0]}
		if len(result) > 1 {
			modifiers := result[1]
			regex.Options = modifiers[:len(modifiers)-1]
		}
		return bson.M{"$regex": regex}
	case strings.HasPrefix(cast.ToString(value), "exists("):
		return bson.M{"$exists": parseOperatorValue(value, "exists") == "true"}
	case value == "true" || value == "false":
		return value == "true"
	default:
		return value
	}
}

func mapFilters(filter bson.M) bson.M {
	for key, value := range filter {
		if slices.Contains(complexOperators, key) {
			filter["$"+key] = lo.Map(strings.Split(cast.ToString(value), ","), func(kv string, index int) any {
				key, value := strings.Split(kv, "=")[0], strings.Split(kv, "=")[1]
				return bson.M{key: mapValue(value)}
			})
			delete(filter, key)
		} else {
			complexOp, found := lo.Find(complexOperators, func(op string) bool {
				return strings.HasPrefix(cast.ToString(value), op+"(")
			})
			if complexOp != "" && found {
				values := strings.Split(cast.ToString(parseOperatorValue(value, complexOp)), ",")
				filter["$"+complexOp] = lo.Map(values, func(subValue string, index int) any {
					return bson.M{key: mapValue(subValue)}
				})
				delete(filter, key)
			} else {
				filter[key] = mapValue(value)
			}
		}
	}
	return filter
}
