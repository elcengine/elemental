package e_tests

import (
	"errors"
	"testing"

	e_constants "github.com/elcengine/elemental/constants"
	e_mocks "github.com/elcengine/elemental/tests/mocks"
	e_test_setup "github.com/elcengine/elemental/tests/setup"

	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCoreRead(t *testing.T) {
	t.Parallel()

	e_test_setup.SeededConnection(t.Name())

	UserModel := UserModel.SetDatabase(t.Name())

	Convey("Read users", t, func() {
		Convey("Find all users", func() {
			users := UserModel.Find().ExecTT()
			So(users, ShouldHaveLength, len(e_mocks.Users))
		})
		Convey("Find all with a limit of 2", func() {
			users := UserModel.Find().Limit(2).ExecTT()
			So(users, ShouldHaveLength, 2)
		})
		Convey("Find all with a limit of 2 and skip 2", func() {
			Convey("In order of skip -> limit", func() {
				users := UserModel.Find().Skip(2).Limit(2).ExecTT()
				So(users, ShouldHaveLength, 2)
				So(users[0].Name, ShouldEqual, e_mocks.Eredin.Name)
				So(users[1].Name, ShouldEqual, e_mocks.Caranthir.Name)
			})
			Convey("In order of limit -> skip", func() {
				users := UserModel.Find().Limit(2).Skip(2).ExecTT()
				So(users, ShouldHaveLength, 2)
				So(users[0].Name, ShouldEqual, e_mocks.Eredin.Name)
				So(users[1].Name, ShouldEqual, e_mocks.Caranthir.Name)
			})
		})
		Convey("Find all users with a filter query", func() {
			users := UserModel.Find(primitive.M{"name": e_mocks.Ciri.Name}).ExecTT()
			So(users, ShouldHaveLength, 1)
			So(users[0].Name, ShouldEqual, e_mocks.Ciri.Name)
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
			user := UserModel.FindOne(primitive.M{"age": e_mocks.Geralt.Age}).ExecPtr()
			So(user, ShouldNotBeNil)
			So(user.Name, ShouldEqual, e_mocks.Geralt.Name)
		})
		Convey("Find a user with a filter query which has no matching documents", func() {
			user := UserModel.FindOne(primitive.M{"name": "Yarpen Zigrin"}).ExecPtr()
			So(user, ShouldBeNil)
		})
		Convey("Find first user", func() {
			user := UserModel.FindOne().ExecPtr()
			So(user, ShouldNotBeNil)
			So(user.Name, ShouldEqual, e_mocks.Ciri.Name)
		})
		Convey("Find user by ID", func() {
			user := UserModel.FindOne().ExecPtr()
			userById := UserModel.FindByID(user.ID).ExecPtr()
			So(userById, ShouldNotBeNil)
			So(userById.Name, ShouldEqual, e_mocks.Ciri.Name)
		})
		Convey("Count users", func() {
			count := UserModel.CountDocuments().ExecInt()
			So(count, ShouldEqual, len(e_mocks.Users))
		})
		Convey("Find all users in descending order of age", func() {
			Convey("In conjuntion with a primitive map", func() {
				users := UserModel.Find().Sort(primitive.M{"age": -1}).ExecTT()
				So(users[0].Name, ShouldEqual, e_mocks.Vesemir.Name)
				So(users[1].Name, ShouldEqual, e_mocks.Imlerith.Name)
				So(users[2].Name, ShouldEqual, e_mocks.Caranthir.Name)
				So(users[3].Name, ShouldEqual, e_mocks.Geralt.Name)
			})
			Convey("In conjuntion with key-value args", func() {
				users := UserModel.Find().Sort("age", -1).ExecTT()
				So(users[0].Name, ShouldEqual, e_mocks.Vesemir.Name)
				So(users[1].Name, ShouldEqual, e_mocks.Imlerith.Name)
				So(users[2].Name, ShouldEqual, e_mocks.Caranthir.Name)
				So(users[3].Name, ShouldEqual, e_mocks.Geralt.Name)
			})
		})
		Convey("Find all users in descending order of age but ascending order of name", func() {
			users := UserModel.Find().Sort("age", -1, "name", 1).ExecTT()
			So(users[0].Name, ShouldEqual, e_mocks.Vesemir.Name)
			So(users[1].Name, ShouldEqual, e_mocks.Imlerith.Name)
			So(users[2].Name, ShouldEqual, e_mocks.Caranthir.Name)
			So(users[3].Name, ShouldEqual, e_mocks.Geralt.Name)
			So(users[4].Name, ShouldEqual, e_mocks.Yennefer.Name)
		})
		Convey("Find all users in descending order of age and name", func() {
			users := UserModel.Find().Sort("age", -1, "name", -1).ExecTT()
			So(users[0].Name, ShouldEqual, e_mocks.Vesemir.Name)
			So(users[1].Name, ShouldEqual, e_mocks.Imlerith.Name)
			So(users[2].Name, ShouldEqual, e_mocks.Caranthir.Name)
			So(users[3].Name, ShouldEqual, e_mocks.Yennefer.Name)
			So(users[4].Name, ShouldEqual, e_mocks.Geralt.Name)
		})
		Convey("Must panic when finding with invalid sort arguments", func() {
			So(func() {
				UserModel.Find().Sort("age", 1, "name").Exec()
			}, ShouldPanicWith, e_constants.ErrMustPairSortArguments)
		})
		Convey("Find all distinct witcher schools", func() {
			schools := UserModel.Distinct("school").ExecStringSlice()
			So(schools, ShouldHaveLength, 2)
			So(schools, ShouldContain, e_mocks.WolfSchool)
			So(schools, ShouldContain, "")
		})
	})
}
