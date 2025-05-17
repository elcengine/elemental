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

func EnsureObjectID(id any) primitive.ObjectID {
	if idStr, ok := id.(string); ok {
		parsed, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			return primitive.NilObjectID
		}
		return parsed
	} else if idRaw, ok := id.(primitive.ObjectID); ok {
		return idRaw
	} else if idPtr, ok := id.(*primitive.ObjectID); ok {
		return *idPtr
	}
	return primitive.NilObjectID
}
