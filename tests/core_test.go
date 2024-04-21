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
	Weapons    []string           `json:"weapons" bson:"weapons"`
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
	"Occupation": {
		Type: reflect.String,
	},
	"Weapons": {
		Type:    reflect.Slice,
		Default: []string{},
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
					Weapons:    []string{"Silver sword", "Mahakaman battle hammer", "Battle Axe", "Crossbow", "Steel sword"},
				},
				{
					Name: "Eredin Bréacc Glas",
				},
				{
					Name:       "Caranthir",
					Age:        120,
					Occupation: "Mage",
					Weapons:    []string{"Staff"},
				},
				{
					Name:       "Imlerith",
					Age:        150,
					Occupation: "General",
					Weapons:    []string{"Mace", "Battle Axe"},
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
		Convey("Find all with a limit of 2", func() {
			users := UserModel.Find().Limit(2).Exec().([]User)
			So(len(users), ShouldEqual, 2)
		})
		Convey("Find all with a limit of 2 and skip 2", func() {
			Convey("In order of skip -> limit", func() {
				users := UserModel.Find().Skip(2).Limit(2).Exec().([]User)
				So(len(users), ShouldEqual, 2)
				So(users[0].Name, ShouldEqual, "Eredin Bréacc Glas")
				So(users[1].Name, ShouldEqual, "Caranthir")
			})
			Convey("In order of limit -> skip", func() {
				users := UserModel.Find().Limit(2).Skip(2).Exec().([]User)
				So(len(users), ShouldEqual, 2)
				So(users[0].Name, ShouldEqual, "Eredin Bréacc Glas")
				So(users[1].Name, ShouldEqual, "Caranthir")
			})
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
		Convey("Find all where age is 100", func() {
			users := UserModel.Where("age").Equals(100).Exec().([]User)
			So(len(users), ShouldEqual, 1)
			So(users[0].Name, ShouldEqual, "Geralt of Rivia")
		})
		Convey("Find all where age is greater than 50", func() {
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
		Convey("Find where age is between 90 and 110", func() {
			Convey("In conjuntion with find", func() {
				users := UserModel.Find(primitive.M{"$and": []primitive.M{
					{"age": primitive.M{"$gte": 90}},
					{"age": primitive.M{"$lte": 110}},
				}}).Exec().([]User)
				So(len(users), ShouldEqual, 1)
				So(users[0].Name, ShouldEqual, "Geralt of Rivia")
			})
			Convey("In conjuntion with where", func() {
				users := UserModel.Where("age").GreaterThanOrEquals(90).Where("age").LessThanOrEquals(110).Exec().([]User)
				So(len(users), ShouldEqual, 1)
				So(users[0].Name, ShouldEqual, "Geralt of Rivia")
			})
			Convey("In conjuntion with between", func() {
				users := UserModel.Where("age").Between(90, 110).Exec().([]User)
				So(len(users), ShouldEqual, 1)
				So(users[0].Name, ShouldEqual, "Geralt of Rivia")
			})
		})
		Convey("Find where age is 120 or 150", func() {
			Convey("In conjuntion with find", func() {
				users := UserModel.Find(primitive.M{"$or": []primitive.M{
					{"age": 120},
					{"age": 150},
				}}).Exec().([]User)
				So(len(users), ShouldEqual, 2)
				So(users[0].Name, ShouldEqual, "Caranthir")
				So(users[1].Name, ShouldEqual, "Imlerith")
			})
			Convey("In conjuntion with in", func() {
				users := UserModel.Where("age").In(120, 150).Exec().([]User)
				So(len(users), ShouldEqual, 2)
				So(users[0].Name, ShouldEqual, "Caranthir")
				So(users[1].Name, ShouldEqual, "Imlerith")
			})
		})
		Convey("Find where age is not 18", func() {
			Convey("In conjuntion with find", func() {
				users := UserModel.Find(primitive.M{"age": primitive.M{"$ne": 18}}).Exec().([]User)
				So(len(users), ShouldEqual, 3)
			})
			Convey("In conjuntion with not equals", func() {
				users := UserModel.Where("age").NotEquals(18).Exec().([]User)
				So(len(users), ShouldEqual, 3)
			})
			Convey("In conjuntion with not in", func() {
				users := UserModel.Where("age").NotIn(18).Exec().([]User)
				So(len(users), ShouldEqual, 3)
			})
		})
		Convey("Find where weapon list contains Battle Axe", func() {
			Convey("In conjuntion with find", func() {
				users := UserModel.Find(primitive.M{"weapons": "Battle Axe"}).Exec().([]User)
				So(len(users), ShouldEqual, 2)
				So(users[0].Name, ShouldEqual, "Geralt of Rivia")
				So(users[1].Name, ShouldEqual, "Imlerith")
			})
			Convey("In conjuntion with element match", func() {
				users := UserModel.Where("weapons").ElementMatches(primitive.M{"$eq": "Battle Axe"}).Exec().([]User)
				So(len(users), ShouldEqual, 2)
				So(users[0].Name, ShouldEqual, "Geralt of Rivia")
				So(users[1].Name, ShouldEqual, "Imlerith")
			})
		})
		Convey("Find where occupation exists", func() {
			users := UserModel.Where("occupation").Exists(true).Exec().([]User)
				So(len(users), ShouldEqual, 4)
		})
		Convey("Find where occupation does not exist", func() {
			users := UserModel.Where("occupation").Exists(false).Exec().([]User)
			So(len(users), ShouldEqual, 2)
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