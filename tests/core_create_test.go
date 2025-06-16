package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	elemental "github.com/elcengine/elemental/core"
	"github.com/elcengine/elemental/tests/fixtures"
	"github.com/elcengine/elemental/tests/fixtures/mocks"
	ts "github.com/elcengine/elemental/tests/fixtures/setup"

	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCoreCreate(t *testing.T) {
	t.Parallel()

	ts.Connection(t.Name())

	UserModel := UserModel.SetDatabase(t.Name())
	MonsterModel := MonsterModel.SetDatabase(t.Name())

	Convey("Create users", t, func() {
		Convey("Create a single user", func() {
			user := UserModel.Create(mocks.Ciri).ExecT()
			So(user.ID.IsZero(), ShouldBeFalse)
			So(user.Name, ShouldEqual, mocks.Ciri.Name)
			So(user.Age, ShouldEqual, fixtures.DefaultUserAge)
			So(user.CreatedAt.Unix(), ShouldBeBetweenOrEqual, time.Now().Add(-10*time.Second).Unix(), time.Now().Unix())
			So(user.UpdatedAt.Unix(), ShouldBeBetweenOrEqual, time.Now().Add(-10*time.Second).Unix(), time.Now().Unix())
		})
		Convey("Create many users", func() {
			users := UserModel.CreateMany(mocks.Users[1:]).ExecTT()
			So(len(users), ShouldEqual, len(mocks.Users[1:]))
			So(users[0].ID.IsZero(), ShouldBeFalse)
			So(users[1].ID.IsZero(), ShouldBeFalse)
			So(users[0].Name, ShouldEqual, mocks.Geralt.Name)
			So(users[1].Name, ShouldEqual, mocks.Eredin.Name)
			So(users[0].Age, ShouldEqual, mocks.Geralt.Age)
			So(users[1].Age, ShouldEqual, fixtures.DefaultUserAge)
		})
		Convey("Create a single user in a different database", func() {
			TEMPORARY_DB := fmt.Sprintf("%s_%s", t.Name(), "temporary_1")
			user := UserModel.Create(mocks.Ciri).SetDatabase(TEMPORARY_DB).ExecT()
			So(user.ID.IsZero(), ShouldBeFalse)
			var newUser User
			elemental.UseDatabase(TEMPORARY_DB).Collection(UserModel.Collection().Name()).FindOne(context.TODO(), primitive.M{"_id": user.ID}).Decode(&newUser)
			So(newUser.Name, ShouldEqual, mocks.Ciri.Name)
		})
		Convey("Create a single user in a different collection in a different database", func() {
			TEMPORARY_DB := fmt.Sprintf("%s_%s", t.Name(), "temporary_2")
			user := UserModel.Create(mocks.Geralt).SetDatabase(TEMPORARY_DB).SetCollection("witchers").ExecT()
			So(user.ID.IsZero(), ShouldBeFalse)
			var newUser User
			elemental.UseDatabase(TEMPORARY_DB).Collection("witchers").FindOne(context.TODO(), primitive.M{"_id": user.ID}).Decode(&newUser)
			So(newUser.Name, ShouldEqual, mocks.Geralt.Name)
		})
	})
	Convey("Create a monster which has a sub schema with defaults", t, func() {
		monster := MonsterModel.Create(Monster{
			Name:     "Katakan",
			Category: "Vampire",
		}).ExecT()
		So(monster.ID.IsZero(), ShouldBeFalse)
		So(monster.Name, ShouldEqual, "Katakan")
		So(monster.CreatedAt.Unix(), ShouldBeBetweenOrEqual, time.Now().Add(-10*time.Second).Unix(), time.Now().Unix())
		So(monster.UpdatedAt.Unix(), ShouldBeBetweenOrEqual, time.Now().Add(-10*time.Second).Unix(), time.Now().Unix())
		So(monster.Weaknesses.Signs, ShouldContain, "Igni")
		So(monster.Weaknesses.InvulnerableTo, ShouldContain, "Steel")
	})
}
