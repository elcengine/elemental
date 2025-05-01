package e_tests

import (
	"testing"

	e_mocks "github.com/elcengine/elemental/tests/mocks"
	e_test_setup "github.com/elcengine/elemental/tests/setup"

	"github.com/samber/lo"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson"
)

func TestCoreReadIs(t *testing.T) {
	t.Parallel()

	e_test_setup.SeededConnection(t.Name())

	Convey("Read users where", t, func() {
		Convey("Name is of type string", func() {
			users := UserModel.Where("name").IsType(bson.TypeString).Exec().([]User)
			So(len(users), ShouldEqual, len(e_mocks.Users))
		})
		Convey("School is null", func() {
			users := UserModel.Where("school").IsNull().Exec().([]User)
			So(len(users), ShouldEqual, len(lo.Filter(e_mocks.Users, func(u User, _ int) bool {
				return u.School == nil
			})))
		})
	})
}
