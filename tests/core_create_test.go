package e_tests

import (
	"elemental/tests/base"
	"elemental/tests/mocks"
	"elemental/tests/setup"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCoreCreate(t *testing.T) {

	e_test_setup.Connection()

	defer e_test_setup.Teardown()

	Convey("Create users", t, func() {
		Convey("Create a single user", func() {
			user := UserModel.Create(e_mocks.Ciri)
			So(user.ID, ShouldNotBeNil)
			So(user.Name, ShouldEqual, e_mocks.Ciri.Name)
			So(user.Age, ShouldEqual, e_test_base.DefaultAge)
			So(user.CreatedAt.Unix(), ShouldBeBetweenOrEqual, time.Now().Add(-10*time.Second).Unix(), time.Now().Unix())
			So(user.UpdatedAt.Unix(), ShouldBeBetweenOrEqual, time.Now().Add(-10*time.Second).Unix(), time.Now().Unix())
		})
		Convey("Create many users", func() {
			users := UserModel.InsertMany(e_mocks.Users[1:])
			So(len(users), ShouldEqual, len(e_mocks.Users[1:]))
			So(users[0].ID, ShouldNotBeNil)
			So(users[1].ID, ShouldNotBeNil)
			So(users[0].Name, ShouldEqual, e_mocks.Geralt.Name)
			So(users[1].Name, ShouldEqual, e_mocks.Eredin.Name)
			So(users[0].Age, ShouldEqual, e_mocks.Geralt.Age)
			So(users[1].Age, ShouldEqual, e_test_base.DefaultAge)
		})
	})
}
