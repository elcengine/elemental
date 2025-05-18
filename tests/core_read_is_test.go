package tests

import (
	"testing"

	"github.com/elcengine/elemental/tests/fixtures/mocks"
	ts "github.com/elcengine/elemental/tests/fixtures/setup"

	"github.com/samber/lo"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson"
)

func TestCoreReadIs(t *testing.T) {
	t.Parallel()

	ts.SeededConnection(t.Name())

	UserModel := UserModel.SetDatabase(t.Name())

	Convey("Read users where", t, func() {
		Convey("Name is of type string", func() {
			users := UserModel.Where("name").IsType(bson.TypeString).ExecTT()
			So(len(users), ShouldEqual, len(mocks.Users))
		})
		Convey("School is null", func() {
			users := UserModel.Where("school").IsNull().ExecTT()
			So(len(users), ShouldEqual, len(lo.Filter(mocks.Users, func(u User, _ int) bool {
				return u.School == nil
			})))
		})
	})
}
