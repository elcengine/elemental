package tests

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/elcengine/elemental/tests/fixtures/mocks"
	ts "github.com/elcengine/elemental/tests/fixtures/setup"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCoreDelete(t *testing.T) {
	t.Parallel()

	ts.Connection(t.Name())

	UserModel := UserModel.SetDatabase(t.Name())

	UserModel.InsertMany(mocks.Users).Exec()

	Convey("Delete users", t, func() {
		Convey("Find and delete first user", func() {
			user := UserModel.FindOneAndDelete().ExecT()
			So(user.Name, ShouldEqual, mocks.Ciri.Name)
			So(UserModel.FindByID(user.ID).Exec(), ShouldBeNil)
		})
		Convey("Find and delete user by ID", func() {
			user := UserModel.FindOne().ExecT()
			So(user.Name, ShouldEqual, mocks.Geralt.Name)
			deletedUser := UserModel.FindByIDAndDelete(user.ID).ExecT()
			So(deletedUser.Name, ShouldEqual, mocks.Geralt.Name)
			So(UserModel.FindByID(user.ID).Exec(), ShouldBeNil)

			Convey("Find and delete user by ID (Hex String)", func() {
				user := UserModel.FindOne(primitive.M{"name": mocks.Imlerith.Name}).ExecT()
				So(user.Name, ShouldEqual, mocks.Imlerith.Name)
				deletedUser := UserModel.FindByIDAndDelete(user.ID.Hex()).ExecT()
				So(deletedUser.Name, ShouldEqual, mocks.Imlerith.Name)
				So(UserModel.FindByID(user.ID).Exec(), ShouldBeNil)
			})
		})
		Convey("Delete a user document", func() {
			user := UserModel.FindOne().ExecT()
			So(user.Name, ShouldEqual, mocks.Eredin.Name)
			UserModel.Delete(user).Exec()
			So(UserModel.FindByID(user.ID).Exec(), ShouldBeNil)
		})
		Convey("Delete a user by ID", func() {
			user := UserModel.FindOne().ExecT()
			So(user.Name, ShouldEqual, mocks.Caranthir.Name)
			UserModel.DeleteByID(user.ID).Exec()
			So(UserModel.FindByID(user.ID).Exec(), ShouldBeNil)
		})
		Convey("Delete all remaining users", func() {
			UserModel.DeleteMany().Exec()
			So(UserModel.Find().Exec(), ShouldBeEmpty)
		})
	})
}
