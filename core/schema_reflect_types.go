package elemental

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
)

// Inbuilt types

var Slice = reflect.Slice
var Map = reflect.Map
var Struct = reflect.Struct
var Interface = reflect.Interface
var Array = reflect.Array
var Bool = reflect.Bool
var Int = reflect.Int
var Int32 = reflect.Int32
var Int64 = reflect.Int64
var Uint = reflect.Uint
var Uint8 = reflect.Uint8
var Uint16 = reflect.Uint16
var Uint32 = reflect.Uint32
var Uint64 = reflect.Uint64
var Float32 = reflect.Float32
var Float64 = reflect.Float64
var String = reflect.String

// Custom types

var ObjectID = reflect.TypeOf(primitive.NilObjectID).Kind()
