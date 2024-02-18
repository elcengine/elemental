package e_tests

import (
	"elemental/connection"
	"elemental/core"
	"elemental/tests/mocks"
	"fmt"
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	elemental.Model[User]
	ID   primitive.ObjectID `bson:"_id"`
	Name string             `bson:"name"`
}

func (u User) Schema() elemental.Schema {
	return elemental.NewSchema(map[string]elemental.Field{
		"ID": {
			Disabled: true,
		},
		"Name": {
			Type:     reflect.String,
			Required: true,
		},
	}, elemental.SchemaOptions{
		Collection: "users",
	})
}

func TestCore(t *testing.T) {
	Convey("Test basic crud operations", t, func() {
		e_connection.ConnectLite(e_mocks.URI)
		Convey("Create a user", func() {
			user := User{
				Name: "John",
			}
			fmt.Println(24234)
			user.Create()
			So(user.ID, ShouldNotBeNil)
		})
	})
}
