package elemental

import (
	"fmt"
	"reflect"
	// "go.mongodb.org/mongo-driver/bson/primitive"
)

// "reflect"

// "go.mongodb.org/mongo-driver/bson/primitive"

type ClusterOp[T any] struct {
	model      *Model[T]
	result     *any
	operations []Operation
}

// type Operation[T any] func(c *ClusterOp[T])
type Operation func()

func Cluster[T any](model *Model[T], connection *string) ClusterOp[T] {
	fmt.Println("ClusterOp constructed")
	if connection != nil {
		model.SetConnection(*connection)
	}

	c := ClusterOp[T]{
		model:      model,
		result:     nil,
		operations: []Operation{},
	}

	return c
}

// func UseCluster[T any](m Model[T], connection *string, op Operation[T]) ClusterOp[T] {
// 	return Cluster(&m, connection, &op)
// }

func (c ClusterOp[T]) Populate(connectionName string, collection []string) {
	res := c.model.Find().Populate("monster").Populate("kingdom").Exec()
	c.result = &res
}

func (c ClusterOp[T]) PopulateOp(model any) ClusterOp[T] {
	c.operations = append(c.operations, func() {
		fn := reflect.ValueOf(c.populate)
		args := []reflect.Value{
			reflect.ValueOf(model),
		}
		fn.Call(args)
	})
	return c
}

func (c ClusterOp[T]) populate(model any) {
	fmt.Printf("hi")
	vModel := reflect.ValueOf(model)
	if vModel.Kind() == reflect.Ptr {
		vModel = vModel.Elem()
	}
	if vModel.Kind() != reflect.Struct {
		fmt.Errorf("expected a struct or a pointer to struct")
	}

	modelName := vModel.FieldByName("Name").Interface().(string)
	fmt.Printf("Model name: %s\n", modelName)
	if c.result == nil {
		fmt.Errorf("prev op result is nil")
		return
	}
	prevOpResult := *c.result
	vPrevOpResult := reflect.ValueOf(prevOpResult)
	id := vPrevOpResult.FieldByName("ID").Interface().(string)
	fmt.Printf("ID: %s\n", id)
	// objectId, err := primitive.ObjectIDFromHex(id)
	// if err != nil {
	// 	panic(err)
	// }
	// result := model.FindByID(objectId).Exec()
	// prevOpResult.(map[string]any)[modelName] = result
}

func (c ClusterOp[T]) Exec() any {
	r := c.model.Exec()
	fmt.Println("executed Exec")
	c.result = &r
	if c.result == nil {
		fmt.Errorf("result is nil")
	}
	for _, op := range c.operations {
		(op)()
	}
	return *c.result
}

// User.ClusterOp((c ClusterOp) => {
//   c.populate("")
// })
