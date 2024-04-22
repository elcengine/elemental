package e_tests

import (
	"elemental/tests/base"
	"elemental/tests/mocks"
	"elemental/tests/setup"
	"elemental/utils"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestCoreUpdate(t *testing.T) {

	e_test_setup.SeededConnection()

	defer e_test_setup.Teardown()

	Convey("Update users", t, func() {
		Convey("Find and update first user", func() {
			user := e_utils.Cast[User](UserModel.FindOneAndUpdate(nil, primitive.M{"name": "Zireael"}).Exec())
			So(user.Name, ShouldEqual, e_mocks.Ciri.Name)
			updatedUser := e_utils.Cast[User](UserModel.FindOne(primitive.M{"name": "Zireael"}).Exec())
			So(updatedUser.Age, ShouldEqual, e_test_base.DefaultAge)
		})
		Convey("Find and update first user and return updated document", func() {
			opts := options.FindOneAndUpdateOptions{}
			opts.SetReturnDocument(options.After)
			user := e_utils.Cast[User](UserModel.FindOneAndUpdate(nil, User{
				Name: "Zireael",
			}, &opts).Exec())
			So(user.Name, ShouldEqual, "Zireael")
			Convey("In conjunction with New", func() {
				user := e_utils.Cast[User](UserModel.FindOneAndUpdate(nil, User{
					Name: "Swallow",
				}).New().Exec())
				So(user.Name, ShouldEqual, "Swallow")
			})
		})
		Convey("Find and update user by ID", func() {
			user := e_utils.Cast[User](UserModel.FindOne().Where("name", e_mocks.Geralt.Name).Exec())
			updatedUser := e_utils.Cast[User](UserModel.FindByIDAndUpdate(user.ID, User{
				Name: "White Wolf",
			}).Exec())
			So(updatedUser.Name, ShouldEqual, e_mocks.Geralt.Name)
			So(updatedUser.Age, ShouldEqual, e_mocks.Geralt.Age)
			updatedUser = e_utils.Cast[User](UserModel.FindByID(user.ID).Exec())
			So(updatedUser.Name, ShouldEqual, "White Wolf")
		})
		Convey("Update user with upsert", func() {
			UserModel.UpdateOne(&primitive.M{"name": "Triss"}, User{
				Name: "Triss",
			}).Upsert().Exec()
			So(UserModel.Where("name", "Triss").Exec(), ShouldHaveLength, 1)
		})
		Convey("Update a user document", func() {
			user := e_utils.Cast[User](UserModel.FindOne().Where("name", e_mocks.Eredin.Name).Exec())
			user.Age = 200
			UserModel.Save(user)
			updatedUser := e_utils.Cast[User](UserModel.FindByID(user.ID).Exec())
			So(updatedUser.Age, ShouldEqual, 200)
		})
		Convey("Find and replace first user", func() {
			user := e_utils.Cast[User](UserModel.FindOneAndReplace(nil, User{
				Name: "Zireael",
			}).Exec())
			So(user.Name, ShouldEqual, "Swallow")
			updatedUser := e_utils.Cast[User](UserModel.FindOne(primitive.M{"name": "Zireael"}).Exec())
			So(updatedUser.Age, ShouldEqual, 0)
		})
		Convey("Find and replace user by ID", func() {
			user := e_utils.Cast[User](UserModel.FindOne().Where("name", e_mocks.Imlerith.Name).Exec())
			updatedUser := e_utils.Cast[User](UserModel.FindByIDAndReplace(user.ID, User{
				Name: "Imlerith",
				Age:  151,
			}).Exec())
			So(updatedUser.Name, ShouldEqual, e_mocks.Imlerith.Name)
			updatedUser = e_utils.Cast[User](UserModel.FindByID(user.ID).Exec())
			So(updatedUser.Age, ShouldEqual, 151)
			So(updatedUser.Occupation, ShouldBeZeroValue)
		})
		Convey("Update all remaining users to have only daggers", func() {
			UserModel.UpdateMany(nil, User{
				Weapons: []string{"Dagger"},
			}).Exec()
			users := UserModel.Where("weapons", "Dagger").Exec()
			So(users, ShouldHaveLength, len(e_mocks.Users)+1)
		})
	})
}
