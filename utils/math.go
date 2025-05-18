//nolint:exhaustive
package utils

import (
	"fmt"
	"reflect"
)

func parseReflectValue(ar, br reflect.Value) (reflect.Value, reflect.Value) {
	if ar.Kind() == reflect.Int && br.Kind() == reflect.Float64 {
		ar = reflect.ValueOf(float64(ar.Int()))
	}
	if br.Kind() == reflect.Int && ar.Kind() == reflect.Float64 {
		br = reflect.ValueOf(float64(br.Int()))
	}
	if ar.Kind() == reflect.Int32 && br.Kind() == reflect.Float64 {
		ar = reflect.ValueOf(float64(ar.Int()))
	}
	if br.Kind() == reflect.Int32 && ar.Kind() == reflect.Float64 {
		br = reflect.ValueOf(float64(br.Int()))
	}
	if ar.Kind() == reflect.Int64 && br.Kind() == reflect.Float64 {
		ar = reflect.ValueOf(float64(ar.Int()))
	}
	if br.Kind() == reflect.Int64 && ar.Kind() == reflect.Float64 {
		br = reflect.ValueOf(float64(br.Int()))
	}
	if ar.Kind() == reflect.Int32 && br.Kind() == reflect.Int {
		ar = reflect.ValueOf(int(ar.Int()))
	}
	if br.Kind() == reflect.Int32 && ar.Kind() == reflect.Int {
		br = reflect.ValueOf(int(br.Int()))
	}
	if ar.Kind() == reflect.Int64 && br.Kind() == reflect.Int {
		ar = reflect.ValueOf(int(ar.Int()))
	}
	if br.Kind() == reflect.Int64 && ar.Kind() == reflect.Int {
		br = reflect.ValueOf(int(br.Int()))
	}
	if ar.Kind() == reflect.String && br.Kind() == reflect.Int {
		br = reflect.ValueOf(fmt.Sprintf("%d", br.Int()))
	}
	if ar.Kind() == reflect.Int && br.Kind() == reflect.String {
		ar = reflect.ValueOf(fmt.Sprintf("%d", ar.Int()))
	}
	if ar.Kind() != br.Kind() {
		panic(fmt.Errorf("type mismatch: cannot compare %s with %s", ar.Kind(), br.Kind()))
	}
	return ar, br
}

func LT(a, b any) bool {
	ar := reflect.ValueOf(a)
	br := reflect.ValueOf(b)
	ar, br = parseReflectValue(ar, br)
	switch ar.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return ar.Int() < br.Int()
	case reflect.Float32, reflect.Float64:
		return ar.Float() < br.Float()
	case reflect.String:
		return ar.String() < br.String()
	default:
		panic(fmt.Errorf("unsupported type: %s", ar.Kind()))
	}
}

func LTE(a, b any) bool {
	ar := reflect.ValueOf(a)
	br := reflect.ValueOf(b)
	ar, br = parseReflectValue(ar, br)
	switch ar.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return ar.Int() <= br.Int()
	case reflect.Float32, reflect.Float64:
		return ar.Float() <= br.Float()
	case reflect.String:
		return ar.String() <= br.String()
	default:
		panic(fmt.Errorf("unsupported type: %s", ar.Kind()))
	}
}

func GT(a, b any) bool {
	ar := reflect.ValueOf(a)
	br := reflect.ValueOf(b)
	ar, br = parseReflectValue(ar, br)
	switch ar.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return ar.Int() > br.Int()
	case reflect.Float32, reflect.Float64:
		return ar.Float() > br.Float()
	case reflect.String:
		return ar.String() > br.String()
	default:
		panic(fmt.Errorf("unsupported type: %s", ar.Kind()))
	}
}

func GTE(a, b any) bool {
	ar := reflect.ValueOf(a)
	br := reflect.ValueOf(b)
	ar, br = parseReflectValue(ar, br)
	switch ar.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return ar.Int() >= br.Int()
	case reflect.Float32, reflect.Float64:
		return ar.Float() >= br.Float()
	case reflect.String:
		return ar.String() >= br.String()
	default:
		panic(fmt.Errorf("unsupported type: %s", ar.Kind()))
	}
}

func EQ(a, b any) bool {
	ar := reflect.ValueOf(a)
	br := reflect.ValueOf(b)
	ar, br = parseReflectValue(ar, br)
	switch ar.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return ar.Int() == br.Int()
	case reflect.Float32, reflect.Float64:
		return ar.Float() == br.Float()
	case reflect.String:
		return ar.String() == br.String()
	case reflect.Bool:
		return ar.Bool() == br.Bool()
	default:
		panic(fmt.Errorf("unsupported type: %s", ar.Kind()))
	}
}
