package e_utils

import (
	"fmt"
	"reflect"
	"strconv"
	"go.mongodb.org/mongo-driver/bson"
)

func setField(field reflect.Value, defaultVal string) error {
	if !field.CanSet() {
		return fmt.Errorf("Can't set value\n")
	}
	switch field.Kind() {
	case reflect.Int:
		if val, err := strconv.ParseInt(defaultVal, 10, 64); err == nil {
			field.Set(reflect.ValueOf(int(val)).Convert(field.Type()))
		}
	case reflect.String:
		field.Set(reflect.ValueOf(defaultVal).Convert(field.Type()))
	}
	return nil
}

func SetDefaults(ptr interface{},) error {
	if reflect.TypeOf(ptr).Kind() != reflect.Ptr {
		return fmt.Errorf("Not a pointer")
	}
	v := reflect.ValueOf(ptr).Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		if defaultVal := t.Field(i).Tag.Get("default"); defaultVal != "-" {
			if err := setField(v.Field(i), defaultVal); err != nil {
				return err
			}

		}
	}
	return nil
}

func ToBSON(v interface{}) (doc *bson.M) {
    data, err := bson.Marshal(v)
    if err != nil {
        return nil
    }
	bson.Unmarshal(data, &doc)
	return doc
}