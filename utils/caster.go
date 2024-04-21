package e_utils

import (
	"github.com/clubpay/qlubkit-go"
	"go.mongodb.org/mongo-driver/bson"
)

// Converts any type to a given type based on their bson representation. It partially fills the target in case they are not directly compatible.
func CastBSON[T any](val any) T {
	return FromBSON[T](ToBSON(val))
}

// Converts a given value to a byte array.
func ToBSON(val any) []byte {
	return qkit.Ok(bson.Marshal(val))
}

// Converts a byte array to a given type.
func FromBSON[T any](bytes []byte) T {
	var v T
	bson.Unmarshal(bytes, &v)
	return v
}

// Converts an interface to a bson document
func ToBSONDoc(v interface{}) (doc *bson.M) {
    data, err := bson.Marshal(v)
    if err != nil {
        return nil
    }
	bson.Unmarshal(data, &doc)
	return doc
}