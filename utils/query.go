package e_utils

import (
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func MergedQueryOrDefault(query []primitive.M) primitive.M {
	if len(query) == 0 {
		return primitive.M{}
	}
	return lo.Assign(query...)
}
