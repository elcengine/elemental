package tests

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/elcengine/elemental/tests/fixtures/mocks"
	ts "github.com/elcengine/elemental/tests/fixtures/setup"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCoreSoftDelete(t *testing.T) {
	t.Parallel()

	ts.Connection(t.Name())

	UserModel := UserModel.SetDatabase(t.Name())

	UserModel.InsertMany(mocks.Users).Exec()

	UserModel.EnableSoftDelete()

	defer UserModel.DisableSoftDelete()

	SoSoftDelete := func(id any) {
		t.Helper()
		rawUser := map[string]any{}
		UserModel.Collection().FindOne(context.Background(), primitive.M{"_id": id}).Decode(&rawUser)
		So(rawUser["deleted_at"], ShouldNotBeNil)
	}

	Convey("Soft delete users", t, func() {
		Convey("Find and delete first user", func() {
			user := UserModel.FindOneAndDelete().ExecT()
			So(user.Name, ShouldEqual, mocks.Ciri.Name)
			So(UserModel.FindByID(user.ID).Exec(), ShouldBeNil)
			SoSoftDelete(user.ID)
		})
		Convey("Find and delete user by ID", func() {
			user := UserModel.FindOne().ExecT()
			So(user.Name, ShouldEqual, mocks.Geralt.Name)
			deletedUser := UserModel.FindByIDAndDelete(user.ID).ExecT()
			So(deletedUser.Name, ShouldEqual, mocks.Geralt.Name)
			So(UserModel.FindByID(user.ID).Exec(), ShouldBeNil)
			SoSoftDelete(user.ID)

			Convey("Find and delete user by ID (Hex String)", func() {
				user := UserModel.FindOne(primitive.M{"name": mocks.Imlerith.Name}).ExecT()
				So(user.Name, ShouldEqual, mocks.Imlerith.Name)
				deletedUser := UserModel.FindByIDAndDelete(user.ID.Hex()).ExecT()
				So(deletedUser.Name, ShouldEqual, mocks.Imlerith.Name)
				So(UserModel.FindByID(user.ID).Exec(), ShouldBeNil)
				SoSoftDelete(user.ID)
			})
		})
		Convey("Delete a user document", func() {
			user := UserModel.FindOne().ExecT()
			So(user.Name, ShouldEqual, mocks.Eredin.Name)
			UserModel.Delete(user).Exec()
			So(UserModel.FindByID(user.ID).Exec(), ShouldBeNil)
			SoSoftDelete(user.ID)
		})
		Convey("Delete a user by ID", func() {
			user := UserModel.FindOne().ExecT()
			So(user.Name, ShouldEqual, mocks.Caranthir.Name)
			UserModel.DeleteByID(user.ID).Exec()
			So(UserModel.FindByID(user.ID).Exec(), ShouldBeNil)
			SoSoftDelete(user.ID)
		})
		Convey("Delete all remaining users", func() {
			UserModel.DeleteMany().Exec()
			So(UserModel.Find().Exec(), ShouldBeEmpty)

			cursor, err := UserModel.Collection().Find(context.Background(), primitive.M{})

			So(err, ShouldBeNil)
			defer cursor.Close(context.Background())

			rawUsers := []map[string]any{}
			cursor.All(context.Background(), &rawUsers)

			So(len(rawUsers), ShouldEqual, len(mocks.Users))

			for _, rawUser := range rawUsers {
				So(rawUser["deleted_at"], ShouldNotBeNil)
			}

			So(UserModel.CountDocuments().Exec(), ShouldEqual, 0)
		})
	})
}
