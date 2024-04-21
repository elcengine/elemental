package e_tests

import (
	"context"
	"elemental/connection"
	"elemental/core"
	"elemental/tests/mocks"
	"reflect"
	"testing"
	"time"

	"github.com/clubpay/qlubkit-go"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id"`
	Name       string             `json:"name" bson:"name"`
	Age        int                `json:"age" bson:"age"`
	Occupation string             `json:"occupation" bson:"occupation,omitempty"`
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
	"Age": {
		Type:    reflect.Int,
		Default: 18,
	},
}, elemental.SchemaOptions{
	Collection: "users",
}))

func TestCore(t *testing.T) {

	e_connection.ConnectURI(e_mocks.URI)

	e_connection.UseDefault().Drop(context.TODO())

	Convey("Test basic crud operations", t, func() {
		Convey("Create a user", func() {
			name := "Ciri"
			user := UserModel.Create(User{
				Name: name,
			})
			So(user.ID, ShouldNotBeNil)
			So(user.Name, ShouldEqual, name)
			So(user.Age, ShouldEqual, 18)
			So(user.CreatedAt.Unix(), ShouldBeBetweenOrEqual, time.Now().Add(-10*time.Second).Unix(), time.Now().Unix())
			So(user.UpdatedAt.Unix(), ShouldBeBetweenOrEqual, time.Now().Add(-10*time.Second).Unix(), time.Now().Unix())
		})
		Convey("Create many users", func() {
			users := UserModel.InsertMany([]User{
				{
					Name:       "Geralt of Rivia",
					Age:        100,
					Occupation: "Witcher",
				},
				{
					Name: "Eredin Bréacc Glas",
				},
				{
					Name:       "Caranthir",
					Age:        120,
					Occupation: "Mage",
				},
				{
					Name:       "Imlerith",
					Age:        150,
					Occupation: "General",
				},
			})
			So(len(users), ShouldEqual, 4)
			So(users[0].ID, ShouldNotBeNil)
			So(users[1].ID, ShouldNotBeNil)
			So(users[0].Name, ShouldEqual, "Geralt of Rivia")
			So(users[1].Name, ShouldEqual, "Eredin Bréacc Glas")
			So(users[0].Age, ShouldEqual, 100)
			So(users[1].Age, ShouldEqual, 18)
		})
		Convey("Find all users", func() {
			users := UserModel.Find().Exec().([]User)
			So(len(users), ShouldEqual, 5)
		})
		Convey("Filter users", func() {
			users := UserModel.Find(primitive.M{"age": 18}).Exec().([]User)
			So(len(users), ShouldEqual, 2)
			So(users[0].Name, ShouldEqual, "Ciri")
		})
		Convey("Find a user", func() {
			user := qkit.Cast[User](UserModel.FindOne(primitive.M{"age": 100}).Exec())
			So(user, ShouldNotBeNil)
			So(user.Name, ShouldEqual, "Geralt of Rivia")
		})
		Convey("Find first user", func() {
			user := qkit.Cast[User](UserModel.FindOne().Exec())
			So(user, ShouldNotBeNil)
			So(user.Name, ShouldEqual, "Ciri")
		})
		Convey("Find user by ID", func() {
			name := "Yennefer of Vengerberg"
			user := UserModel.Create(User{
				Name:       name,
				Occupation: "Mage",
			})
			found := qkit.Cast[User](UserModel.FindByID(user.ID).Exec())
			So(found, ShouldNotBeNil)
			So(found.Name, ShouldEqual, name)
		})
		Convey("Find users whose age is 100", func() {
			users := UserModel.Where("age").Equals(100).Exec().([]User)
			So(len(users), ShouldEqual, 1)
			So(users[0].Name, ShouldEqual, "Geralt of Rivia")
		})
		Convey("Find where age is greater than 50", func() {
			users := UserModel.Where("age").GreaterThan(50).Exec().([]User)
			So(len(users), ShouldEqual, 3)
			So(users[0].Name, ShouldEqual, "Geralt of Rivia")
		})
		Convey("Find a mage where age is greater than 50", func() {
			Convey("In conjuntion with find", func() {
				users := UserModel.Find(primitive.M{"occupation": "Mage"}).Where("age").GreaterThan(50).Exec().([]User)
				So(len(users), ShouldEqual, 1)
				So(users[0].Name, ShouldEqual, "Caranthir")
			})
			Convey("In conjuntion with find one", func() {
				user := qkit.Cast[User](UserModel.FindOne(primitive.M{"occupation": "Mage"}).Where("age").GreaterThan(50).Exec())
				So(user, ShouldNotBeNil)
				So(user.Name, ShouldEqual, "Caranthir")
			})
			Convey("In conjuntion with equals", func() {
				users := UserModel.Where("age").GreaterThan(50).Where("occupation").Equals("Mage").Exec().([]User)
				So(len(users), ShouldEqual, 1)
				So(users[0].Name, ShouldEqual, "Caranthir")
			})
		})
		Convey("Count users", func() {
			count := UserModel.CountDocuments().Exec().(int64)
			So(count, ShouldEqual, 6)
			Convey("In conjuntion with greater than", func() {
				count := UserModel.Where("age").GreaterThan(50).CountDocuments().Exec().(int64)
				So(count, ShouldEqual, 3)
			})
		})
	})
}
