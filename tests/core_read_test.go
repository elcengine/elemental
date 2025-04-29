package e_tests

import (
	"github.com/elcengine/elemental/constants"
	"github.com/elcengine/elemental/tests/mocks"
	"github.com/elcengine/elemental/tests/setup"
	"github.com/elcengine/elemental/utils"
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCoreRead(t *testing.T) {

	e_test_setup.SeededConnection()

	defer e_test_setup.Teardown()

	Convey("Read users", t, func() {
		Convey("Find all users", func() {
			t.Parallel()
			users := UserModel.Find().Exec().([]User)
			So(users, ShouldHaveLength, len(e_mocks.Users))
		})
		Convey("Find all with a limit of 2", func() {
			t.Parallel()
			users := UserModel.Find().Limit(2).Exec().([]User)
			So(users, ShouldHaveLength, 2)
		})
		Convey("Find all with a limit of 2 and skip 2", func() {
			t.Parallel()
			Convey("In order of skip -> limit", func() {
				t.Parallel()
				users := UserModel.Find().Skip(2).Limit(2).Exec().([]User)
				So(users, ShouldHaveLength, 2)
				So(users[0].Name, ShouldEqual, e_mocks.Eredin.Name)
				So(users[1].Name, ShouldEqual, e_mocks.Caranthir.Name)
			})
			Convey("In order of limit -> skip", func() {
				t.Parallel()
				users := UserModel.Find().Limit(2).Skip(2).Exec().([]User)
				So(users, ShouldHaveLength, 2)
				So(users[0].Name, ShouldEqual, e_mocks.Eredin.Name)
				So(users[1].Name, ShouldEqual, e_mocks.Caranthir.Name)
			})
		})
		Convey("Find all users with a filter query", func() {
			t.Parallel()
			users := UserModel.Find(primitive.M{"name": e_mocks.Ciri.Name}).Exec().([]User)
			So(users, ShouldHaveLength, 1)
			So(users[0].Name, ShouldEqual, e_mocks.Ciri.Name)
		})
		Convey("Find all users with a filter query which has no matching documents", func() {
			t.Parallel()
			users := UserModel.Find(primitive.M{"name": "Yarpen Zigrin"}).Exec().([]User)
			So(users, ShouldHaveLength, 0)
			Convey("With or fail", func() {
				t.Parallel()
				So(func() {
					UserModel.Find(primitive.M{"name": "Yarpen Zigrin"}).OrFail().Exec()
				}, ShouldPanicWith, errors.New("no results found matching the given query"))
			})
			Convey("With or fail and custom error", func() {
				t.Parallel()
				err := errors.New("no user found")
				So(func() {
					UserModel.Find(primitive.M{"name": "Yarpen Zigrin"}).OrFail(err).Exec()
				}, ShouldPanicWith, err)
			})
		})
		Convey("Find a user with a filter query", func() {
			t.Parallel()
			user := UserModel.FindOne(primitive.M{"age": e_mocks.Geralt.Age}).Exec()
			So(user, ShouldNotBeNil)
			So(e_utils.Cast[User](user).Name, ShouldEqual, e_mocks.Geralt.Name)
		})
		Convey("Find a user with a filter query which has no matching documents", func() {
			t.Parallel()
			user := UserModel.FindOne(primitive.M{"name": "Yarpen Zigrin"}).Exec()
			So(user, ShouldBeNil)
		})
		Convey("Find first user", func() {
			t.Parallel()
			user := UserModel.FindOne().Exec()
			So(user, ShouldNotBeNil)
			So(e_utils.Cast[User](user).Name, ShouldEqual, e_mocks.Ciri.Name)
		})
		Convey("Find user by ID", func() {
			t.Parallel()
			user := e_utils.Cast[User](UserModel.FindOne().Exec())
			userById := UserModel.FindByID(user.ID).Exec()
			So(userById, ShouldNotBeNil)
			So(e_utils.Cast[User](userById).Name, ShouldEqual, e_mocks.Ciri.Name)
		})
		Convey("Count users", func() {
			t.Parallel()
			count := UserModel.CountDocuments().Exec().(int64)
			So(count, ShouldEqual, len(e_mocks.Users))
		})
		Convey("Find all users in descending order of age", func() {
			t.Parallel()
			Convey("In conjuntion with a primitive map", func() {
				t.Parallel()
				users := UserModel.Find().Sort(primitive.M{"age": -1}).Exec().([]User)
				So(users[0].Name, ShouldEqual, e_mocks.Vesemir.Name)
				So(users[1].Name, ShouldEqual, e_mocks.Imlerith.Name)
				So(users[2].Name, ShouldEqual, e_mocks.Caranthir.Name)
				So(users[3].Name, ShouldEqual, e_mocks.Geralt.Name)
			})
			Convey("In conjuntion with key-value args", func() {
				t.Parallel()
				users := UserModel.Find().Sort("age", -1).Exec().([]User)
				So(users[0].Name, ShouldEqual, e_mocks.Vesemir.Name)
				So(users[1].Name, ShouldEqual, e_mocks.Imlerith.Name)
				So(users[2].Name, ShouldEqual, e_mocks.Caranthir.Name)
				So(users[3].Name, ShouldEqual, e_mocks.Geralt.Name)
			})
		})
		Convey("Find all users in descending order of age but ascending order of name", func() {
			t.Parallel()
			users := UserModel.Find().Sort("age", -1, "name", 1).Exec().([]User)
			So(users[0].Name, ShouldEqual, e_mocks.Vesemir.Name)
			So(users[1].Name, ShouldEqual, e_mocks.Imlerith.Name)
			So(users[2].Name, ShouldEqual, e_mocks.Caranthir.Name)
			So(users[3].Name, ShouldEqual, e_mocks.Geralt.Name)
			So(users[4].Name, ShouldEqual, e_mocks.Yennefer.Name)
		})
		Convey("Find all users in descending order of age and name", func() {
			t.Parallel()
			users := UserModel.Find().Sort("age", -1, "name", -1).Exec().([]User)
			So(users[0].Name, ShouldEqual, e_mocks.Vesemir.Name)
			So(users[1].Name, ShouldEqual, e_mocks.Imlerith.Name)
			So(users[2].Name, ShouldEqual, e_mocks.Caranthir.Name)
			So(users[3].Name, ShouldEqual, e_mocks.Yennefer.Name)
			So(users[4].Name, ShouldEqual, e_mocks.Geralt.Name)
		})
		Convey("Must panic when finding with invalid sort arguments", func() {
			t.Parallel()
			So(func() {
				UserModel.Find().Sort("age", 1, "name").Exec()
			}, ShouldPanicWith, e_constants.ErrMustPairSortArguments)
		})
		Convey("Find all distinct witcher schools", func() {
			t.Parallel()
			schools := UserModel.Distinct("school").Exec().([]string)
			So(schools, ShouldHaveLength, 2)
			So(schools, ShouldContain, e_mocks.WolfSchool)
			So(schools, ShouldContain, "")
		})
	})
}
