package elemental

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"strings"
)

type ClusterOp[T any] struct {
	model      *Model[T]
	result     **map[string]any
	operations []Operation
}

type Operation func()

func Cluster[T any](model *Model[T], connection *string) ClusterOp[T] {
	if connection != nil {
		model.SetConnection(*connection)
	}
	var b *map[string]any

	c := ClusterOp[T]{
		model:      model,
		result:     &b,
		operations: []Operation{},
	}

	return c
}

func (c ClusterOp[T]) Populate(model Model[any]) ClusterOp[T] {
	c.operations = append(c.operations, func() {
		c.populate(model)
	})
	return c
}

func (c ClusterOp[T]) populate(model Model[any]) ClusterOp[T] {
	if c.result == nil {
		fmt.Println("result is nil")
		return c
	}
	r := **c.result

	// finding the id field
	var idField string
	for k, v := range c.model.Schema.Definitions {
		if v.IsRefID && v.Ref == model.Name {
			idField = k
			break
		}
	}

	id, ok := r[idField]
	if !ok {
		fmt.Println("Incorrect id field. Field not found")
		return c
	}
	str, ok := id.(string)
	if !ok {
		fmt.Println("Incorrect id field. Not a string")
		return c
	}
	str = strings.TrimPrefix(str, "ObjectID(\"")
	str = strings.TrimSuffix(str, "\")")
	objectID, err := primitive.ObjectIDFromHex(str)
	if err != nil {
		fmt.Println("Incorrect id field. Not a valid ObjectID: " + str)
		fmt.Println(err)
		return c
	}

	res := model.FindByID(objectID).Exec()
	if res == nil {
		fmt.Println("Result is nil")
		return c
	}

	if reflect.TypeOf(res) != reflect.TypeOf(primitive.D{}) {
		fmt.Println("Result is not a primitive.D")
		return c
	}
	resMap := convertDToMap(res.(primitive.D))

	r[model.Name] = resMap
	*c.result = &r

	return c
}

func (c ClusterOp[T]) Exec() any {
	r := c.model.Exec()
	if r == nil {
		fmt.Println("result is nil")
		return nil
	}
	rmap := structToMap(r)
	*c.result = &rmap
	for _, op := range c.operations {
		(op)()
	}
	return **c.result
}
