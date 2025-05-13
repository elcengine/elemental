package e_tests

import (
	"github.com/elcengine/elemental/tests/mocks"
	"github.com/elcengine/elemental/tests/setup"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCoreDelete(t *testing.T) {
	t.Parallel()

	e_test_setup.Connection(t.Name())

	UserModel := UserModel.SetDatabase(t.Name())

	UserModel.InsertMany(e_mocks.Users).Exec()

	Convey("Delete users", t, func() {
		Convey("Find and delete first user", func() {
			user := UserModel.FindOneAndDelete().ExecT()
			So(user.Name, ShouldEqual, e_mocks.Ciri.Name)
			So(UserModel.FindByID(user.ID).Exec(), ShouldBeNil)
		})
		Convey("Find and delete user by ID", func() {
			user := UserModel.FindOne().ExecT()
			So(user.Name, ShouldEqual, e_mocks.Geralt.Name)
			deletedUser := UserModel.FindByIdAndDelete(user.ID).ExecT()
			So(deletedUser.Name, ShouldEqual, e_mocks.Geralt.Name)
			So(UserModel.FindByID(user.ID).Exec(), ShouldBeNil)
		})
		Convey("Delete a user document", func() {
			user := UserModel.FindOne().ExecT()
			So(user.Name, ShouldEqual, e_mocks.Eredin.Name)
			UserModel.Delete(user).Exec()
			So(UserModel.FindByID(user.ID).Exec(), ShouldBeNil)
		})
		Convey("Delete a user by ID", func() {
			user := UserModel.FindOne().ExecT()
			So(user.Name, ShouldEqual, e_mocks.Caranthir.Name)
			UserModel.DeleteByID(user.ID).Exec()
			So(UserModel.FindByID(user.ID).Exec(), ShouldBeNil)
		})
		Convey("Delete all remaining users", func() {
			UserModel.DeleteMany().Exec()
			So(UserModel.Find().Exec(), ShouldBeEmpty)
		})
	})
}
