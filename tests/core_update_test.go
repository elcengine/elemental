package tests

import (
	"testing"

	"github.com/elcengine/elemental/tests/fixtures"
	"github.com/elcengine/elemental/tests/fixtures/mocks"
	ts "github.com/elcengine/elemental/tests/fixtures/setup"
	"github.com/samber/lo"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestCoreUpdate(t *testing.T) {
	t.Parallel()

	ts.SeededConnection(t.Name())

	UserModel := UserModel.SetDatabase(t.Name())

	Convey("Update users", t, func() {
		Convey("Find and update first user", func() {
			user := UserModel.FindOneAndUpdate(nil, primitive.M{"name": "Zireael"}).ExecT()
			So(user.Name, ShouldEqual, mocks.Ciri.Name)
			updatedUser := UserModel.FindOne(primitive.M{"name": "Zireael"}).ExecT()
			So(updatedUser.Age, ShouldEqual, fixtures.DefaultUserAge)
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
			user := UserModel.FindOne().Where("name", mocks.Geralt.Name).ExecT()
			updatedUser := UserModel.FindByIDAndUpdate(user.ID, User{
				Name: "White Wolf",
			}).ExecT()
			So(updatedUser.Name, ShouldEqual, mocks.Geralt.Name)
			So(updatedUser.Age, ShouldEqual, mocks.Geralt.Age)
			updatedUser = UserModel.FindByID(user.ID).ExecT()
			So(updatedUser.Name, ShouldEqual, "White Wolf")
		})
		Convey("Update user with upsert", func() {
			UserModel.UpdateOne(&primitive.M{"name": "Triss"}, User{
				Name: "Triss",
			}).Upsert().Exec()
			So(UserModel.Where("name", "Triss").Exec(), ShouldHaveLength, 1)
		})
		Convey("Update user with upsert within options", func() {
			UserModel.UpdateOne(&primitive.M{"name": "Letho"}, User{
				Name: "Letho",
			}, &options.UpdateOptions{
				Upsert: lo.ToPtr(true),
			}).Exec()
			So(UserModel.Where("name", "Letho").Exec(), ShouldHaveLength, 1)
		})
		Convey("Update a user with a pointer document", func() {
			user := User{
				Name: "Foltest",
				Age:  50,
			}
			UserModel.Create(user).Exec()
			UserModel.UpdateOne(&primitive.M{"name": user.Name}, &User{
				Age: 51,
			}).Exec()
			So(UserModel.FindOne().Where("name", user.Name).ExecT().Age, ShouldEqual, 51)
		})
		Convey("Update user by ID", func() {
			user := UserModel.FindOne().Where("name", "Triss").ExecT()
			UserModel.UpdateByID(user.ID, User{
				Name: "Triss Merigold",
			}).Exec()
			updatedUser := UserModel.FindByID(user.ID).ExecT()
			So(updatedUser.Name, ShouldEqual, "Triss Merigold")

			Convey("Update user by ID (Hex String)", func() {
				UserModel.UpdateByID(user.ID.Hex(), User{
					Name: "Triss Merigold the Fearless",
				}).Exec()
				updatedUser := UserModel.FindByID(user.ID).ExecT()
				So(updatedUser.Name, ShouldEqual, "Triss Merigold the Fearless")
			})
		})
		Convey("Update a user document", func() {
			user := UserModel.FindOne().Where("name", mocks.Eredin.Name).ExecT()
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
			user := UserModel.FindOne().Where("name", mocks.Imlerith.Name).ExecT()
			updatedUser := UserModel.FindByIDAndReplace(user.ID, User{
				Name: "Imlerith",
				Age:  151,
			}).ExecT()
			So(updatedUser.Name, ShouldEqual, mocks.Imlerith.Name)
			updatedUser = UserModel.FindByID(user.ID).ExecT()
			So(updatedUser.Age, ShouldEqual, 151)
			So(updatedUser.Occupation, ShouldBeZeroValue)
		})
		Convey("Update all remaining users to have only daggers", func() {
			UserModel.UpdateMany(nil, User{
				Weapons: []string{"Dagger"},
			}).Exec()
			users := UserModel.Where("weapons", "Dagger").ExecTT()
			So(users, ShouldHaveLength, len(mocks.Users)+3)
		})
		Convey("Increment age of a user", func() {
			UserModel.Where("name", mocks.Vesemir.Name).Inc("age", 1).Exec()
			updatedUser := UserModel.FindOne().Where("name", mocks.Vesemir.Name).ExecT()
			So(updatedUser.Age, ShouldEqual, mocks.Vesemir.Age+1)
		})
		Convey("Decrement age of a user", func() {
			UserModel.Where("name", mocks.Vesemir.Name).Dec("age", 1).Exec()
			updatedUser := UserModel.FindOne().Where("name", mocks.Vesemir.Name).ExecT()
			So(updatedUser.Age, ShouldEqual, mocks.Vesemir.Age)
		})
		Convey("Multiply age of a user", func() {
			UserModel.Where("name", mocks.Vesemir.Name).Mul("age", 2).Exec()
			updatedUser := UserModel.FindOne().Where("name", mocks.Vesemir.Name).ExecT()
			So(updatedUser.Age, ShouldEqual, mocks.Vesemir.Age*2)
		})
		Convey("Divide age of a user", func() {
			UserModel.Where("name", mocks.Vesemir.Name).Div("age", 2).Exec()
			updatedUser := UserModel.FindOne().Where("name", mocks.Vesemir.Name).ExecT()
			So(updatedUser.Age, ShouldEqual, mocks.Vesemir.Age)
		})
		Convey("Rename occupation field to profession", func() {
			UserModel.Rename("occupation", "profession").Exec()
			user := UserModel.FindOne().Where("name", mocks.Vesemir.Name).ExecT()
			So(user.Occupation, ShouldBeZeroValue)
		})
		Convey("Unset school of a user", func() {
			UserModel.Where("name", mocks.Vesemir.Name).Unset("school").Exec()
			updatedUser := UserModel.FindOne().Where("name", mocks.Vesemir.Name).ExecT()
			So(updatedUser.School, ShouldBeZeroValue)
		})
		Convey("Set school of a user using set", func() {
			UserModel.Where("name", mocks.Vesemir.Name).Set(primitive.M{
				"school": "Kaer Morhen",
			}).Exec()
			updatedUser := UserModel.FindOne().Where("name", mocks.Vesemir.Name).ExecT()
			So(updatedUser.School, ShouldEqual, lo.ToPtr("Kaer Morhen"))
		})
		Convey("Add a weapon to a user", func() {
			UserModel.Where("name", mocks.Vesemir.Name).Push("weapons", "Xiphos").Exec()
			updatedUser := UserModel.FindOne().Where("name", mocks.Vesemir.Name).ExecT()
			So(updatedUser.Weapons, ShouldContain, "Xiphos")
		})
		Convey("Add multiple weapons to a user", func() {
			UserModel.Where("name", mocks.Vesemir.Name).Push("weapons", "Ulfberht", "Mace").Exec()
			updatedUser := UserModel.FindOne().Where("name", mocks.Vesemir.Name).ExecT()
			So(updatedUser.Weapons, ShouldContain, "Xiphos")
			So(updatedUser.Weapons, ShouldContain, "Ulfberht")
			So(updatedUser.Weapons, ShouldContain, "Mace")
		})
		Convey("Remove a weapon from a user", func() {
			UserModel.Where("name", mocks.Vesemir.Name).Pull("weapons", "Xiphos").Exec()
			updatedUser := UserModel.FindOne().Where("name", mocks.Vesemir.Name).ExecT()
			So(updatedUser.Weapons, ShouldNotContain, "Xiphos")
		})
		Convey("Remove multiple weapons from a user", func() {
			UserModel.Where("name", mocks.Vesemir.Name).PullAll("weapons", "Ulfberht", "Mace").Exec()
			updatedUser := UserModel.FindOne().Where("name", mocks.Vesemir.Name).ExecT()
			So(updatedUser.Weapons, ShouldNotContain, "Ulfberht")
			So(updatedUser.Weapons, ShouldNotContain, "Mace")
		})
		Convey("Remove last weapon from a user", func() {
			user := UserModel.FindOne().Where("name", mocks.Vesemir.Name).ExecT()
			UserModel.Where("name", mocks.Vesemir.Name).Pop("weapons").Exec()
			updatedUser := UserModel.FindOne().Where("name", mocks.Vesemir.Name).ExecT()
			So(len(updatedUser.Weapons), ShouldEqual, len(user.Weapons)-1)
		})
		Convey("Try adding the same weapon multiple times to a user", func() {
			UserModel.Where("name", mocks.Vesemir.Name).AddToSet("weapons", "Longsword", "Longsword", "Longsword").Exec()
			updatedUser := UserModel.FindOne().Where("name", mocks.Vesemir.Name).ExecT()
			So(len(lo.Filter(updatedUser.Weapons, func(w string, _ int) bool {
				return w == "Longsword"
			})), ShouldEqual, 1)
		})
		Convey("Remove first weapon from a user", func() {
			user := UserModel.FindOne().Where("name", mocks.Vesemir.Name).ExecT()
			UserModel.Where("name", mocks.Vesemir.Name).Shift("weapons").Exec()
			updatedUser := UserModel.FindOne().Where("name", mocks.Vesemir.Name).ExecT()
			So(len(updatedUser.Weapons), ShouldEqual, len(user.Weapons)-1)
		})
		Convey("Set age of user only if it is greater than current age", func() {
			UserModel.Where("name", mocks.Vesemir.Name).Max("age", 50).Exec()
			updatedUser := UserModel.FindOne().Where("name", mocks.Vesemir.Name).ExecT()
			So(updatedUser.Age, ShouldEqual, mocks.Vesemir.Age)

			UserModel.Where("name", mocks.Vesemir.Name).Max("age", 301).Exec()
			updatedUser = UserModel.FindOne().Where("name", mocks.Vesemir.Name).ExecT()
			So(updatedUser.Age, ShouldEqual, 301)
		})
		Convey("Set age of user only if it is less than current age", func() {
			UserModel.Where("name", mocks.Yennefer.Name).Min("age", 200).Exec()
			updatedUser := UserModel.FindOne().Where("name", mocks.Yennefer.Name).ExecT()
			So(updatedUser.Age, ShouldEqual, mocks.Yennefer.Age)

			UserModel.Where("name", mocks.Yennefer.Name).Min("age", 80).Exec()
			updatedUser = UserModel.FindOne().Where("name", mocks.Yennefer.Name).ExecT()
			So(updatedUser.Age, ShouldEqual, 80)
		})
	})
}
