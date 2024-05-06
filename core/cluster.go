package elemental

import (
	"reflect"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClusterOp[T any] struct {
	model      *Model[T]
	result     *any
	operations []*Operation[T]
}

type Operation[T any] func(c *ClusterOp[T])

func Cluster[T any](model *Model[T], connection *string, op *Operation[T]) ClusterOp[T] {
	if connection != nil {
		model.SetConnection(*connection)
	}

	c := ClusterOp[T]{
		model:      model,
		result:     nil,
		operations: []*Operation[T]{op},
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

func (c ClusterOp[T]) populat(connectionName string, model any) {
	// check if the model is valid
	// if _, ok := Models[modelName]; !ok {
	// 	panic(e_constants.ErrInvalidModel)
	// }

	r := (*c.result).([]any)[0]
	if reflect.ValueOf(r).Kind() == reflect.Slice {
		r = r.(primitive.D)[0]
	}
	// get the id of the populating model
	// v := reflect.ValueOf(r)
	// for i := 0; i < v.NumField(); i++ {
	// 	n := v.Type().Field(i).Name
	// 	fmt.Println(n)
	// }

	slice := r.(primitive.D)
	castedModel := model.(Model[any])
	for i := range slice {
		if castedModel.Name == slice[i].Key {
		}
		// fmt.Printf("%s : %s\n", slice[i].Key, slice[i].Value)

	}

	// model := e_utils.Cast[Model[T]](Models[modelName])
}

func (c ClusterOp[T]) Exec() any {
	// c.result = &r
	// for _, op := range c.operations {
	// 	(*op)(&c)
	// }
	// r := (*c.model).Exec()
	return c.model.Find().Populate("monster").Populate("kingdom").Exec()
}

// User.ClusterOp((c ClusterOp) => {
//   c.populate("")
// })
