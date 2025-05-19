package tests

import (
	"testing"

	"github.com/elcengine/elemental/tests/fixtures/mocks"
	ts "github.com/elcengine/elemental/tests/fixtures/setup"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestCoreReplace(t *testing.T) {
	t.Parallel()

	ts.SeededConnection(t.Name())

	UserModel := UserModel.SetDatabase(t.Name())

	Convey("Replace users", t, func() {
		Convey("Replace user by name", func() {
			result := UserModel.ReplaceOne(&primitive.M{"name": mocks.Geralt.Name}, User{
				Name: "Radovid",
			}).Exec().(*mongo.UpdateResult)
			So(result.MatchedCount, ShouldEqual, 1)
			So(UserModel.Where("name", mocks.Geralt).Exec(), ShouldBeEmpty)
			replacedUser := UserModel.FindOne().Where("name", "Radovid").ExecPtr()
			So(replacedUser, ShouldNotBeNil)
			So(replacedUser.Name, ShouldEqual, "Radovid")
			So(replacedUser.Age, ShouldEqual, 0)
		})

		Convey("Replace user by ID", func() {
			user := UserModel.FindOne().Where("name", mocks.Vesemir.Name).ExecT()
			result := UserModel.ReplaceByID(user.ID, User{
				Name: "Jormungandr",
			}).Exec().(*mongo.UpdateResult)
			So(result.MatchedCount, ShouldEqual, 1)
			replacedUser := UserModel.FindByID(user.ID).ExecPtr()
			So(replacedUser, ShouldNotBeNil)
			So(replacedUser.Name, ShouldEqual, "Jormungandr")
			So(replacedUser.Age, ShouldEqual, 0)
		})
	})
}
