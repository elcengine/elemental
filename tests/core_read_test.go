package tests

import (
	"errors"
	"testing"

	elemental "github.com/elcengine/elemental/core"
	"github.com/elcengine/elemental/tests/fixtures/mocks"
	ts "github.com/elcengine/elemental/tests/fixtures/setup"

	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCoreRead(t *testing.T) {
	t.Parallel()

	ts.SeededConnection(t.Name())

	UserModel := UserModel.SetDatabase(t.Name())

	Convey("Read users", t, func() {
		Convey("Find all users", func() {
			users := UserModel.Find().ExecTT()
			So(users, ShouldHaveLength, len(mocks.Users))
		})
		Convey("Find all with a limit of 2", func() {
			users := UserModel.Find().Limit(2).ExecTT()
			So(users, ShouldHaveLength, 2)
		})
		Convey("Find all with a limit of 2 and skip 2", func() {
			Convey("In order of skip -> limit", func() {
				users := UserModel.Find().Skip(2).Limit(2).ExecTT()
				So(users, ShouldHaveLength, 2)
				So(users[0].Name, ShouldEqual, mocks.Eredin.Name)
				So(users[1].Name, ShouldEqual, mocks.Caranthir.Name)
			})
			Convey("In order of limit -> skip", func() {
				users := UserModel.Find().Limit(2).Skip(2).ExecTT()
				So(users, ShouldHaveLength, 2)
				So(users[0].Name, ShouldEqual, mocks.Eredin.Name)
				So(users[1].Name, ShouldEqual, mocks.Caranthir.Name)
			})
		})
		Convey("Find all users with a filter query", func() {
			users := UserModel.Find(primitive.M{"name": mocks.Ciri.Name}).ExecTT()
			So(users, ShouldHaveLength, 1)
			So(users[0].Name, ShouldEqual, mocks.Ciri.Name)
		})
		Convey("Find all users with a filter query overridden by another filter query", func() {
			users := UserModel.Find(
				primitive.M{"name": mocks.Geralt.Name, "occupation": mocks.Geralt.Occupation},
				primitive.M{"name": mocks.Vesemir.Name},
			).ExecTT()
			So(users, ShouldHaveLength, 1)
			So(users[0].Name, ShouldEqual, mocks.Vesemir.Name)
		})
		Convey("Find all users with a filter query which has no matching documents", func() {
			users := UserModel.Find(primitive.M{"name": "Yarpen Zigrin"}).ExecTT()
			So(users, ShouldHaveLength, 0)
			Convey("With or fail", func() {
				So(func() {
					UserModel.Find(primitive.M{"name": "Yarpen Zigrin"}).OrFail().Exec()
				}, ShouldPanicWith, errors.New("no results found matching the given query"))
			})
			Convey("With or fail and custom error", func() {
				err := errors.New("no user found")
				So(func() {
					UserModel.Find(primitive.M{"name": "Yarpen Zigrin"}).OrFail(err).Exec()
				}, ShouldPanicWith, err)
			})
		})
		Convey("Find a user with a filter query", func() {
			user := UserModel.FindOne(primitive.M{"age": mocks.Geralt.Age}).ExecPtr()
			So(user, ShouldNotBeNil)
			So(user.Name, ShouldEqual, mocks.Geralt.Name)
			Convey("With or fail", func() {
				So(func() {
					UserModel.FindOne(primitive.M{"name": "Yarpen Zigrin"}).OrFail().Exec()
				}, ShouldPanicWith, errors.New("no results found matching the given query"))
			})
			Convey("With or fail and custom error", func() {
				err := errors.New("no user found")
				So(func() {
					UserModel.FindOne(primitive.M{"name": "Yarpen Zigrin"}).OrFail(err).Exec()
				}, ShouldPanicWith, err)
			})
		})
		Convey("Find a user with a filter query which has no matching documents", func() {
			user := UserModel.FindOne(primitive.M{"name": "Yarpen Zigrin"}).ExecPtr()
			So(user, ShouldBeNil)
		})
		Convey("Find first user", func() {
			user := UserModel.FindOne().ExecPtr()
			So(user, ShouldNotBeNil)
			So(user.Name, ShouldEqual, mocks.Ciri.Name)
		})
		Convey("Find user by ID", func() {
			user := UserModel.FindOne().ExecPtr()
			userById := UserModel.FindByID(user.ID).ExecPtr()
			So(userById, ShouldNotBeNil)
			So(userById.Name, ShouldEqual, mocks.Ciri.Name)

			Convey("Find user by ID (Hex String)", func() {
				userById := UserModel.FindByID(user.ID.Hex()).ExecPtr()
				So(userById, ShouldNotBeNil)
				So(userById.Name, ShouldEqual, mocks.Ciri.Name)
			})
			Convey("Find user by ID (Object ID pointer)", func() {
				userById := UserModel.FindByID(&user.ID).ExecPtr()
				So(userById, ShouldNotBeNil)
				So(userById.Name, ShouldEqual, mocks.Ciri.Name)
			})
		})
		Convey("Count users", func() {
			count := UserModel.CountDocuments().ExecInt()
			So(count, ShouldEqual, len(mocks.Users))
		})
		Convey("Find all users in descending order of age", func() {
			Convey("In conjuntion with a primitive map", func() {
				users := UserModel.Find().Sort(primitive.M{"age": -1}).ExecTT()
				So(users[0].Name, ShouldEqual, mocks.Vesemir.Name)
				So(users[1].Name, ShouldEqual, mocks.Imlerith.Name)
				So(users[2].Name, ShouldEqual, mocks.Caranthir.Name)
				So(users[3].Name, ShouldEqual, mocks.Geralt.Name)
			})
			Convey("In conjuntion with key-value args", func() {
				users := UserModel.Find().Sort("age", -1).ExecTT()
				So(users[0].Name, ShouldEqual, mocks.Vesemir.Name)
				So(users[1].Name, ShouldEqual, mocks.Imlerith.Name)
				So(users[2].Name, ShouldEqual, mocks.Caranthir.Name)
				So(users[3].Name, ShouldEqual, mocks.Geralt.Name)
			})
		})
		Convey("Find all users in descending order of age but ascending order of name", func() {
			users := UserModel.Find().Sort("age", -1, "name", 1).ExecTT()
			So(users[0].Name, ShouldEqual, mocks.Vesemir.Name)
			So(users[1].Name, ShouldEqual, mocks.Imlerith.Name)
			So(users[2].Name, ShouldEqual, mocks.Caranthir.Name)
			So(users[3].Name, ShouldEqual, mocks.Geralt.Name)
			So(users[4].Name, ShouldEqual, mocks.Yennefer.Name)
		})
		Convey("Find all users in descending order of age and name", func() {
			users := UserModel.Find().Sort("age", -1, "name", -1).ExecTT()
			So(users[0].Name, ShouldEqual, mocks.Vesemir.Name)
			So(users[1].Name, ShouldEqual, mocks.Imlerith.Name)
			So(users[2].Name, ShouldEqual, mocks.Caranthir.Name)
			So(users[3].Name, ShouldEqual, mocks.Yennefer.Name)
			So(users[4].Name, ShouldEqual, mocks.Geralt.Name)
		})
		Convey("Must panic when finding with invalid sort arguments", func() {
			So(func() {
				UserModel.Find().Sort("age", 1, "name").Exec()
			}, ShouldPanicWith, elemental.ErrMustPairSortArguments)
		})
		Convey("Find all distinct witcher schools", func() {
			schools := UserModel.Distinct("school").ExecSS()
			So(schools, ShouldHaveLength, 2)
			So(schools, ShouldContain, mocks.WolfSchool)
			So(schools, ShouldContain, "")
		})
	})
}
