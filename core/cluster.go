package elemental

import (
	"elemental/constants"
	"fmt"
)

type ClusterOp[T any] struct {
	model      *Model[T]
	result     *any
	operations []*Operation[T]
}

type Operation[T any] func(c *ClusterOp[T]) any

func Cluster[T any](model *Model[T], connection *string, op *Operation[T]) ClusterOp[T] {
	fmt.Println(("DATABASE CLUSTER LAYER CALLED"))
	if connection != nil {
		model.SetConnection(*connection)
	}

	c := ClusterOp[T]{
		model:      model,
		result:     nil,
		operations: []*Operation[T]{op},
	}
	c.operations = append(c.operations, op)

	return c
}

func (c ClusterOp[T]) Populate(connectionName string, modelName string) {
	if _, ok := Models[modelName]; !ok {
		panic(e_constants.ErrURIRequired)
	}
	return
}

func (c ClusterOp[T]) Exec() any {
	r := (*c.model).Exec()
	c.result = &r
	for _, op := range c.operations {
		(*op)(&c)
	}
	return c.result
}

// User.ClusterOp((c ClusterOp) => {
//   c.populate("")
// })
