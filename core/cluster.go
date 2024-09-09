package elemental

import (
	"fmt"
	"reflect"
	"strings"
	// "unsafe"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/bson/primitive"
)

// "reflect"

// "go.mongodb.org/mongo-driver/bson/primitive"

type ClusterOp[T any] struct {
	model      *Model[T]
	result     **any
	operations []Operation
	val        *int
}

// type Operation[T any] func(c *ClusterOp[T])
type Operation func()

func Cluster[T any](model *Model[T], connection *string) ClusterOp[T] {
	fmt.Println("ClusterOp constructed")
	if connection != nil {
		model.SetConnection(*connection)
	}
	var a any
	b := &a

	c := ClusterOp[T]{
		model:      model,
		result:     &b,
		operations: []Operation{},
	}
	c.val = lo.ToPtr(22)

	return c
}

// func UseCluster[T any](m Model[T], connection *string, op Operation[T]) ClusterOp[T] {
// 	return Cluster(&m, connection, &op)
// }

func (c ClusterOp[T]) PopulateOp(model any) ClusterOp[T] {
	c.operations = append(c.operations, func() {
		c.populate(model)
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
		fmt.Printf("expected a struct or a pointer to struct")
	}

	modelName := vModel.FieldByName("Name").Interface().(string)
	fmt.Printf("Model name: %s\n", modelName)
	if c.result == nil {
		fmt.Printf("prev op result is nil")
		return
	}
	prevOpResult := **c.result
	fmt.Printf("PrevResult in populate\n")
	vPrevOpResult := reflect.ValueOf(prevOpResult)
	fmt.Printf("vPrevOpResult type: %v, kind: %v\n", vPrevOpResult.Type(), vPrevOpResult.Kind())
	// if vPrevOpResult.Kind() == reflect.Slice {
	// 	vPrevOpResult = vPrevOpResult.Index(0)
	// }
	fmt.Printf("modelName + ID: %s\n", modelName+"ID")
	objectIdStr := vPrevOpResult.FieldByName(modelName + "ID").Interface()
	objectIdStr, _ = strings.CutPrefix(objectIdStr.(string), "ObjectID(\"")
	objectIdStr, _ = strings.CutSuffix(objectIdStr.(string), "\")")
	fmt.Printf("objectIdStr: %s\n", objectIdStr)
	objectId, err := primitive.ObjectIDFromHex(objectIdStr.(string))
	primitive.NewObjectID()
	if err != nil {
		fmt.Printf("Error converting string to ObjectID: %s\n", err)
	}
	fmt.Printf("PrevResult Monster ID: %s\n", objectId)

	method := vModel.MethodByName("FindByID")
	var ret reflect.Value
	if method.IsValid() {
		args := []reflect.Value{
			reflect.ValueOf(objectId),
		}
		ret = method.Call(args)[0]
	} else {
		fmt.Println("Method not found or invalid")
	}
	method = ret.MethodByName("Exec")
	if method.IsValid() {
		fmt.Println("Fetching populating model")
		retArr := method.Call(nil)
		// HERE RET COULD EITHER BE A SINGLE OBJECT OR MANY OBJECTS
		ret = retArr[len(retArr)-1]
		fmt.Println("Populated model fetched")
	} else {
		fmt.Println("Method not found or invalid")
	}

	fmt.Printf("fieldToSet type: %v, kind: %v\n", ret.Elem().Type(), ret.Elem().Kind())

	v := reflect.ValueOf(**c.result)
	fmt.Printf("v type: %v, kind: %v\n", v.Type(), v.Kind())
	ps := reflect.ValueOf(&v)
	addrStruct := ps.Elem()
	fmt.Printf("s type: %v, kind: %v\n", addrStruct.Type(), addrStruct.Kind())
	fmt.Printf("s.CanAddr(): %v\n", addrStruct.CanAddr())

	vret := reflect.ValueOf(ret)
	printStructFields(vret)

	// vret := reflect.ValueOf(ret)
	// fmt.Printf("VRET FIELDS\n")

	if addrStruct.CanAddr() {
		fmt.Println("addrStruct is addressable!")

		if addrStruct.Kind() == reflect.Struct {
			fmt.Println("addrStruct is a struct!")
			fName := addrStruct.FieldByName(modelName)
			fp := reflect.ValueOf(&fName)
			fmt.Printf("fp type: %v, kind: %v\n", fp.Type(), fp.Kind())
			f := fp.Elem()
			fmt.Printf("f type: %v, kind: %v\n", f.Type(), f.Kind())
			if f.IsValid() {
				fmt.Printf("Field %s is valid!\n", modelName)
				if f.CanAddr() {
					fmt.Printf("Field %s is addressable!\n", modelName)
					// f.Set(vret)
					f.Set(reflect.ValueOf(ret))
					fmt.Printf("Field %s set!\n", modelName)
				} else {
					fmt.Printf("Field %s is not addressable!\n", modelName)
				}
			} else {
				fmt.Printf("Field %s is not valid!\n", modelName)
			}
		} else {
			fmt.Println("AddrPrevOpResult is not a struct!")
		}
	} else {
		fmt.Println("AddrPrevOpResult is not addressable!")
	}

}

func printStructFields(s interface{}) {
	// Convert the interface to reflect.Value and reflect.Type
	val := reflect.ValueOf(s)
	typ := reflect.TypeOf(s)

	// Ensure the reflect.Value is of Kind struct
	if val.Kind() != reflect.Struct {
		fmt.Println("Not a struct")
		val = val.Elem()
	}

	// Loop through the struct fields
	for i := 0; i < val.NumField(); i++ {
		// Get the field name and value
		fieldName := typ.Field(i).Name

		// Print the field name and value
		fmt.Printf("Field: %s\n", fieldName)
		fmt.Printf("Value: %v\n", val.Field(i))
	}
}

func (c ClusterOp[T]) Exec() any {
	r := c.model.Exec()
	fmt.Println("executed Exec")
	if !reflect.ValueOf(r).IsValid() {
		fmt.Printf("result is invalid")
		return nil
	}
	*c.result = &r
	if c.result == nil {
		fmt.Printf("result is nil")
	}
	*c.val = 33
	for _, op := range c.operations {
		(op)()
	}
	return **c.result
}

// User.ClusterOp((c ClusterOp) => {
//   c.populate("")
// })
