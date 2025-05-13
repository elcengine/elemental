package e_tests

import (
	"testing"

	e_test_base "github.com/elcengine/elemental/tests/base"
	e_mocks "github.com/elcengine/elemental/tests/mocks"
	e_test_setup "github.com/elcengine/elemental/tests/setup"
	"github.com/samber/lo"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestCoreUpdate(t *testing.T) {
	t.Parallel()

	e_test_setup.SeededConnection(t.Name())

	UserModel := UserModel.SetDatabase(t.Name())

	Convey("Update users", t, func() {
		Convey("Find and update first user", func() {
			user := UserModel.FindOneAndUpdate(nil, primitive.M{"name": "Zireael"}).ExecT()
			So(user.Name, ShouldEqual, e_mocks.Ciri.Name)
			updatedUser := UserModel.FindOne(primitive.M{"name": "Zireael"}).ExecT()
			So(updatedUser.Age, ShouldEqual, e_test_base.DefaultAge)
		})
		Convey("Find and update first user and return updated document", func() {
			opts := options.FindOneAndUpdateOptions{}
			opts.SetReturnDocument(options.After)
			user := UserModel.FindOneAndUpdate(nil, User{
				Name: "Zireael",
			}, &opts).ExecT()
			So(user.Name, ShouldEqual, "Zireael")
			Convey("In conjunction with New", func() {
				user := UserModel.FindOneAndUpdate(nil, User{
					Name: "Swallow",
				}).New().ExecT()
				So(user.Name, ShouldEqual, "Swallow")
			})
		})
		Convey("Find and update user by ID", func() {
			user := UserModel.FindOne().Where("name", e_mocks.Geralt.Name).ExecT()
			updatedUser := UserModel.FindByIDAndUpdate(user.ID, User{
				Name: "White Wolf",
			}).ExecT()
			So(updatedUser.Name, ShouldEqual, e_mocks.Geralt.Name)
			So(updatedUser.Age, ShouldEqual, e_mocks.Geralt.Age)
			updatedUser = UserModel.FindByID(user.ID).ExecT()
			So(updatedUser.Name, ShouldEqual, "White Wolf")
		})
		Convey("Update user with upsert", func() {
			UserModel.UpdateOne(&primitive.M{"name": "Triss"}, User{
				Name: "Triss",
			}).Upsert().Exec()
			So(UserModel.Where("name", "Triss").Exec(), ShouldHaveLength, 1)
		})
		Convey("Update user by ID", func() {
			user := UserModel.FindOne().Where("name", "Triss").ExecT()
			UserModel.UpdateByID(user.ID, User{
				Name: "Triss Merigold",
			}).Exec()
			updatedUser := UserModel.FindByID(user.ID).ExecT()
			So(updatedUser.Name, ShouldEqual, "Triss Merigold")
		})
		Convey("Update a user document", func() {
			user := UserModel.FindOne().Where("name", e_mocks.Eredin.Name).ExecT()
			user.Age = 200
			UserModel.Save(user).Exec()
			updatedUser := UserModel.FindByID(user.ID).ExecT()
			So(updatedUser.Age, ShouldEqual, 200)
		})
		Convey("Find and replace first user", func() {
			user := UserModel.FindOneAndReplace(nil, User{
				Name: "Zireael",
			}).ExecT()
			So(user.Name, ShouldEqual, "Swallow")
			updatedUser := UserModel.FindOne(primitive.M{"name": "Zireael"}).ExecT()
			So(updatedUser.Age, ShouldEqual, 0)
		})
		Convey("Find and replace user by ID", func() {
			user := UserModel.FindOne().Where("name", e_mocks.Imlerith.Name).ExecT()
			updatedUser := UserModel.FindByIDAndReplace(user.ID, User{
				Name: "Imlerith",
				Age:  151,
			}).ExecT()
			So(updatedUser.Name, ShouldEqual, e_mocks.Imlerith.Name)
			updatedUser = UserModel.FindByID(user.ID).ExecT()
			So(updatedUser.Age, ShouldEqual, 151)
			So(updatedUser.Occupation, ShouldBeZeroValue)
		})
		Convey("Update all remaining users to have only daggers", func() {
			UserModel.UpdateMany(nil, User{
				Weapons: []string{"Dagger"},
			}).Exec()
			users := UserModel.Where("weapons", "Dagger").ExecTT()
			So(users, ShouldHaveLength, len(e_mocks.Users)+1)
		})
		Convey("Increment age of a user", func() {
			UserModel.Where("name", e_mocks.Vesemir.Name).Inc("age", 1).Exec()
			updatedUser := UserModel.FindOne().Where("name", e_mocks.Vesemir.Name).ExecT()
			So(updatedUser.Age, ShouldEqual, e_mocks.Vesemir.Age+1)
		})
		Convey("Decrement age of a user", func() {
			UserModel.Where("name", e_mocks.Vesemir.Name).Dec("age", 1).Exec()
			updatedUser := UserModel.FindOne().Where("name", e_mocks.Vesemir.Name).ExecT()
			So(updatedUser.Age, ShouldEqual, e_mocks.Vesemir.Age)
		})
		Convey("Multiply age of a user", func() {
			UserModel.Where("name", e_mocks.Vesemir.Name).Mul("age", 2).Exec()
			updatedUser := UserModel.FindOne().Where("name", e_mocks.Vesemir.Name).ExecT()
			So(updatedUser.Age, ShouldEqual, e_mocks.Vesemir.Age*2)
		})
		Convey("Divide age of a user", func() {
			UserModel.Where("name", e_mocks.Vesemir.Name).Div("age", 2).Exec()
			updatedUser := UserModel.FindOne().Where("name", e_mocks.Vesemir.Name).ExecT()
			So(updatedUser.Age, ShouldEqual, e_mocks.Vesemir.Age)
		})
		Convey("Rename occupation field to profession", func() {
			UserModel.Rename("occupation", "profession").Exec()
			user := UserModel.FindOne().Where("name", e_mocks.Vesemir.Name).ExecT()
			So(user.Occupation, ShouldBeZeroValue)
		})
		Convey("Unset school of a user", func() {
			UserModel.Where("name", e_mocks.Vesemir.Name).Unset("school").Exec()
			updatedUser := UserModel.FindOne().Where("name", e_mocks.Vesemir.Name).ExecT()
			So(updatedUser.School, ShouldBeZeroValue)
		})
		Convey("Add a weapon to a user", func() {
			UserModel.Where("name", e_mocks.Vesemir.Name).Push("weapons", "Xiphos").Exec()
			updatedUser := UserModel.FindOne().Where("name", e_mocks.Vesemir.Name).ExecT()
			So(updatedUser.Weapons, ShouldContain, "Xiphos")
		})
		Convey("Add multiple weapons to a user", func() {
			UserModel.Where("name", e_mocks.Vesemir.Name).Push("weapons", "Ulfberht", "Mace").Exec()
			updatedUser := UserModel.FindOne().Where("name", e_mocks.Vesemir.Name).ExecT()
			So(updatedUser.Weapons, ShouldContain, "Xiphos")
			So(updatedUser.Weapons, ShouldContain, "Ulfberht")
			So(updatedUser.Weapons, ShouldContain, "Mace")
		})
		Convey("Remove a weapon from a user", func() {
			UserModel.Where("name", e_mocks.Vesemir.Name).Pull("weapons", "Xiphos").Exec()
			updatedUser := UserModel.FindOne().Where("name", e_mocks.Vesemir.Name).ExecT()
			So(updatedUser.Weapons, ShouldNotContain, "Xiphos")
		})
		Convey("Remove multiple weapons from a user", func() {
			UserModel.Where("name", e_mocks.Vesemir.Name).PullAll("weapons", "Ulfberht", "Mace").Exec()
			updatedUser := UserModel.FindOne().Where("name", e_mocks.Vesemir.Name).ExecT()
			So(updatedUser.Weapons, ShouldNotContain, "Ulfberht")
			So(updatedUser.Weapons, ShouldNotContain, "Mace")
		})
		Convey("Remove last weapon from a user", func() {
			user := UserModel.FindOne().Where("name", e_mocks.Vesemir.Name).ExecT()
			UserModel.Where("name", e_mocks.Vesemir.Name).Pop("weapons").Exec()
			updatedUser := UserModel.FindOne().Where("name", e_mocks.Vesemir.Name).ExecT()
			So(len(updatedUser.Weapons), ShouldEqual, len(user.Weapons)-1)
		})
		Convey("Try adding the same weapon multiple times to a user", func() {
			UserModel.Where("name", e_mocks.Vesemir.Name).AddToSet("weapons", "Longsword", "Longsword", "Longsword").Exec()
			updatedUser := UserModel.FindOne().Where("name", e_mocks.Vesemir.Name).ExecT()
			So(len(lo.Filter(updatedUser.Weapons, func(w string, _ int) bool {
				return w == "Longsword"
			})), ShouldEqual, 1)
		})
		Convey("Remove first weapon from a user", func() {
			user := UserModel.FindOne().Where("name", e_mocks.Vesemir.Name).ExecT()
			UserModel.Where("name", e_mocks.Vesemir.Name).Shift("weapons").Exec()
			updatedUser := UserModel.FindOne().Where("name", e_mocks.Vesemir.Name).ExecT()
			So(len(updatedUser.Weapons), ShouldEqual, len(user.Weapons)-1)
		})
	})
}
