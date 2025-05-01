package e_tests

import (
	"github.com/elcengine/elemental/tests/mocks"
	"github.com/elcengine/elemental/tests/setup"
	"github.com/elcengine/elemental/utils"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCoreDelete(t *testing.T) {
	t.Parallel()

	var LocalUserModel = UserModel.Clone().SetCollection("users_for_delete")

	e_test_setup.Connection()

	LocalUserModel.InsertMany(e_mocks.Users).Exec()

	defer e_test_setup.Teardown()

	Convey("Delete users", t, func() {
		Convey("Find and delete first user", func() {
			user := e_utils.Cast[User](LocalUserModel.FindOneAndDelete().Exec())
			So(user.Name, ShouldEqual, e_mocks.Ciri.Name)
			So(LocalUserModel.FindByID(user.ID).Exec(), ShouldBeNil)
		})
		Convey("Find and delete user by ID", func() {
			user := e_utils.Cast[User](LocalUserModel.FindOne().Exec())
			So(user.Name, ShouldEqual, e_mocks.Geralt.Name)
			deletedUser := e_utils.Cast[User](LocalUserModel.FindByIdAndDelete(user.ID).Exec())
			So(deletedUser.Name, ShouldEqual, e_mocks.Geralt.Name)
			So(LocalUserModel.FindByID(user.ID).Exec(), ShouldBeNil)
		})
		Convey("Delete a user document", func() {
			user := e_utils.Cast[User](LocalUserModel.FindOne().Exec())
			So(user.Name, ShouldEqual, e_mocks.Eredin.Name)
			LocalUserModel.Delete(user).Exec()
			So(LocalUserModel.FindByID(user.ID).Exec(), ShouldBeNil)
		})
		Convey("Delete a user by ID", func() {
			user := e_utils.Cast[User](LocalUserModel.FindOne().Exec())
			So(user.Name, ShouldEqual, e_mocks.Caranthir.Name)
			LocalUserModel.DeleteByID(user.ID).Exec()
			So(LocalUserModel.FindByID(user.ID).Exec(), ShouldBeNil)
		})
		Convey("Delete all remaining users", func() {
			LocalUserModel.DeleteMany().Exec()
			So(LocalUserModel.Find().Exec(), ShouldBeEmpty)
		})
	})
}
