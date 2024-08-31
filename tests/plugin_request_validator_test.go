package e_tests

import (
	"github.com/elcengine/elemental/tests/base"
	"github.com/elcengine/elemental/tests/mocks"
	"github.com/elcengine/elemental/tests/setup"
	"github.com/elcengine/elemental/utils"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestRequestValidator(t *testing.T) {

	e_test_setup.SeededConnection()

	defer e_test_setup.Teardown()

	Convey("Basic validations", t, func() {
		Convey("Exists", func() {
			user := e_utils.Cast[User](UserModel.FindOneAndUpdate(nil, primitive.M{"name": "Zireael"}).Exec())
			So(user.Name, ShouldEqual, e_mocks.Ciri.Name)
			updatedUser := e_utils.Cast[User](UserModel.FindOne(primitive.M{"name": "Zireael"}).Exec())
			So(updatedUser.Age, ShouldEqual, e_test_base.DefaultAge)
		})
	})
}
