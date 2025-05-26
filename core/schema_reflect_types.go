package elemental

import (
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FieldType interface {
	String() string
}

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
var Uint32 = reflect.Uint32
var Uint64 = reflect.Uint64
var Float32 = reflect.Float32
var Float64 = reflect.Float64
var String = reflect.String

// Custom types

var Time = reflect.TypeOf(time.Time{})
var ObjectID = reflect.TypeOf(primitive.NilObjectID)
var ObjectIDSlice = reflect.TypeOf([]primitive.ObjectID{})
var StringSlice = reflect.TypeOf([]string{})
var StringMap = reflect.TypeOf(map[string]string{})
var IntSlice = reflect.TypeOf([]int{})
var IntMap = reflect.TypeOf(map[string]int{})
var BoolSlice = reflect.TypeOf([]bool{})
var BoolMap = reflect.TypeOf(map[string]bool{})
var Int32Slice = reflect.TypeOf([]int32{})
var Int32Map = reflect.TypeOf(map[string]int32{})
var Int64Slice = reflect.TypeOf([]int64{})
var Int64Map = reflect.TypeOf(map[string]int64{})
var UintSlice = reflect.TypeOf([]uint{})
var UintMap = reflect.TypeOf(map[string]uint{})
var Uint32Slice = reflect.TypeOf([]uint32{})
var Uint32Map = reflect.TypeOf(map[string]uint32{})
var Uint64Slice = reflect.TypeOf([]uint64{})
var Uint64Map = reflect.TypeOf(map[string]uint64{})
var Float32Slice = reflect.TypeOf([]float32{})
var Float64Slice = reflect.TypeOf([]float64{})
var Float32Map = reflect.TypeOf(map[string]float32{})
var Float64Map = reflect.TypeOf(map[string]float64{})
