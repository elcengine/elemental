package e_tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/elcengine/elemental/connection"
	"github.com/elcengine/elemental/tests/base"
	"github.com/elcengine/elemental/tests/mocks"
	"github.com/elcengine/elemental/tests/setup"

	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCoreCreate(t *testing.T) {
	t.Parallel()

	e_test_setup.Connection(t.Name())

	UserModel := UserModel.SetDatabase(t.Name())
	MonsterModel := MonsterModel.SetDatabase(t.Name())

	Convey("Create users", t, func() {
		Convey("Create a single user", func() {
			user := UserModel.Create(e_mocks.Ciri).Exec().(User)
			So(user.ID, ShouldNotBeNil)
			So(user.Name, ShouldEqual, e_mocks.Ciri.Name)
			So(user.Age, ShouldEqual, e_test_base.DefaultAge)
			So(user.CreatedAt.Unix(), ShouldBeBetweenOrEqual, time.Now().Add(-10*time.Second).Unix(), time.Now().Unix())
			So(user.UpdatedAt.Unix(), ShouldBeBetweenOrEqual, time.Now().Add(-10*time.Second).Unix(), time.Now().Unix())
		})
		Convey("Create many users", func() {
			users := UserModel.InsertMany(e_mocks.Users[1:]).Exec().([]User)
			So(len(users), ShouldEqual, len(e_mocks.Users[1:]))
			So(users[0].ID, ShouldNotBeNil)
			So(users[1].ID, ShouldNotBeNil)
			So(users[0].Name, ShouldEqual, e_mocks.Geralt.Name)
			So(users[1].Name, ShouldEqual, e_mocks.Eredin.Name)
			So(users[0].Age, ShouldEqual, e_mocks.Geralt.Age)
			So(users[1].Age, ShouldEqual, e_test_base.DefaultAge)
		})
		Convey("Create a single user in a different database", func() {
			TEMPORARY_DB := fmt.Sprintf("%s_%s", t.Name(), "temporary_1")
			user := UserModel.Create(e_mocks.Ciri).SetDatabase(TEMPORARY_DB).Exec().(User)
			So(user.ID, ShouldNotBeNil)
			var newUser User
			e_connection.Use(TEMPORARY_DB).Collection(UserModel.Collection().Name()).FindOne(context.TODO(), primitive.M{"_id": user.ID}).Decode(&newUser)
			So(newUser.Name, ShouldEqual, e_mocks.Ciri.Name)
		})
		Convey("Create a single user in a different collection in a different database", func() {
			TEMPORARY_DB := fmt.Sprintf("%s_%s", t.Name(), "temporary_2")
			user := UserModel.Create(e_mocks.Geralt).SetDatabase(TEMPORARY_DB).SetCollection("witchers").Exec().(User)
			So(user.ID, ShouldNotBeNil)
			var newUser User
			e_connection.Use(TEMPORARY_DB).Collection("witchers").FindOne(context.TODO(), primitive.M{"_id": user.ID}).Decode(&newUser)
			So(newUser.Name, ShouldEqual, e_mocks.Geralt.Name)
		})
	})
	Convey("Create a monster which has a sub schema with defaults", t, func() {
		monster := MonsterModel.Create(Monster{
			Name:     "Katakan",
			Category: "Vampire",
		}).Exec().(Monster)
		So(monster.ID, ShouldNotBeNil)
		So(monster.Name, ShouldEqual, "Katakan")
		So(monster.CreatedAt.Unix(), ShouldBeBetweenOrEqual, time.Now().Add(-10*time.Second).Unix(), time.Now().Unix())
		So(monster.UpdatedAt.Unix(), ShouldBeBetweenOrEqual, time.Now().Add(-10*time.Second).Unix(), time.Now().Unix())
		So(monster.Weaknesses.Signs, ShouldContain, "Igni")
		So(monster.Weaknesses.InvulnerableTo, ShouldContain, "Steel")
	})
}
