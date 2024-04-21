package e_tests

import (
	"elemental/connection"
	"elemental/core"
	"elemental/tests/mocks"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/clubpay/qlubkit-go"
	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id"`
	Name       string             `json:"name" bson:"name"`
	Age        int                `json:"age" bson:"age"`
	Occupation string             `json:"occupation" bson:"occupation"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
}

var UserModel = elemental.NewModel[User]("User", elemental.NewSchema(map[string]elemental.Field{
	"Name": {
		Type:     reflect.String,
		Required: true,
		Index: options.IndexOptions{
			Unique: qkit.ValPtr(true),
		},
	},
	"Occupation": {
		Type: reflect.String,
	},
	"Age": {
		Type:    reflect.Int,
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
			name := fmt.Sprintf("Ciri-%d", uuid.New().ID())
			user := UserModel.Create(User{
				Name: name,
			})
			So(user.ID, ShouldNotBeNil)
			So(user.Name, ShouldEqual, name)
			So(user.Age, ShouldEqual, 18)
			So(user.CreatedAt.Unix(), ShouldBeBetweenOrEqual, time.Now().Add(-10*time.Second).Unix(), time.Now().Unix())
			So(user.UpdatedAt.Unix(), ShouldBeBetweenOrEqual, time.Now().Add(-10*time.Second).Unix(), time.Now().Unix())
		})
	})
}
