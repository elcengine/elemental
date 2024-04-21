package e_tests

import (
	"elemental/tests/mocks"
	"elemental/tests/setup"
	"testing"

	"github.com/clubpay/qlubkit-go"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCoreRead(t *testing.T) {

	e_test_setup.SeededConnection()

	defer e_test_setup.Teardown()

	Convey("Read users", t, func() {
		Convey("Find all users", func() {
			users := UserModel.Find().Exec().([]User)
			So(len(users), ShouldEqual, 6)
		})
		Convey("Find all with a limit of 2", func() {
			users := UserModel.Find().Limit(2).Exec().([]User)
			So(len(users), ShouldEqual, 2)
		})
		Convey("Find all with a limit of 2 and skip 2", func() {
			Convey("In order of skip -> limit", func() {
				users := UserModel.Find().Skip(2).Limit(2).Exec().([]User)
				So(len(users), ShouldEqual, 2)
				So(users[0].Name, ShouldEqual, e_mocks.Eredin.Name)
				So(users[1].Name, ShouldEqual, e_mocks.Caranthir.Name)
			})
			Convey("In order of limit -> skip", func() {
				users := UserModel.Find().Limit(2).Skip(2).Exec().([]User)
				So(len(users), ShouldEqual, 2)
				So(users[0].Name, ShouldEqual, e_mocks.Eredin.Name)
				So(users[1].Name, ShouldEqual, e_mocks.Caranthir.Name)
			})
		})
		Convey("Find all users with a filter query", func() {
			users := UserModel.Find(primitive.M{"name": e_mocks.Ciri.Name}).Exec().([]User)
			So(len(users), ShouldEqual, 1)
			So(users[0].Name, ShouldEqual, e_mocks.Ciri.Name)
		})
		Convey("Find a user with a filter query", func() {
			user := qkit.Cast[User](UserModel.FindOne(primitive.M{"age": e_mocks.Geralt.Age}).Exec())
			So(user, ShouldNotBeNil)
			So(user.Name, ShouldEqual, e_mocks.Geralt.Name)
		})
		Convey("Find first user", func() {
			user := qkit.Cast[User](UserModel.FindOne().Exec())
			So(user, ShouldNotBeNil)
			So(user.Name, ShouldEqual, e_mocks.Ciri.Name)
		})
		Convey("Find user by ID", func() {
			user := qkit.Cast[User](UserModel.FindOne().Exec())
			user = qkit.Cast[User](UserModel.FindByID(user.ID).Exec())
			So(user, ShouldNotBeNil)
			So(user.Name, ShouldEqual, e_mocks.Ciri.Name)
		})
		Convey("Count users", func() {
			count := UserModel.CountDocuments().Exec().(int64)
			So(count, ShouldEqual, 6)
		})
	})
}
