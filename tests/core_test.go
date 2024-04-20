package e_tests

import (
	"elemental/connection"
	"elemental/core"
	"elemental/tests/mocks"
	"fmt"
	"reflect"
	"testing"

	"github.com/clubpay/qlubkit-go"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID   primitive.ObjectID `json:"_id" bson:"_id"`
	Name string             `json:"name" bson:"name"`
	Age  int                `json:"age" bson:"age"`
}

var UserModel = elemental.NewModel[User]("User", elemental.NewSchema(map[string]elemental.Field{
	"Name": {
		Type:     reflect.String,
		Required: true,
		Index: options.IndexOptions{
			Unique: qkit.ValPtr(true),
		},
	},
	"Age": {
		Type: reflect.Int,
		Default: 18,
	},
}, elemental.SchemaOptions{
	Collection: "users",
}))

func TestCore(t *testing.T) {
	e_connection.ConnectURI(e_mocks.URI)
	Convey("Test basic crud operations", t, func() {
		e_connection.ConnectURI(e_mocks.URI)
		Convey("Create a user", func() {
			user := User{
				Name: "Akalanka",
			}
			u := UserModel.Create(user)
			fmt.Println("newly created user id", u)
			So(user.ID, ShouldNotBeNil)
		})
	})
}
