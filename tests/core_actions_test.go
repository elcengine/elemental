package e_tests

import (
	"context"
	"slices"
	"testing"

	"github.com/elcengine/elemental/tests/setup"
	"go.mongodb.org/mongo-driver/bson/primitive"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCoreActions(t *testing.T) {
	t.Parallel()

	e_test_setup.Connection(t.Name())

	UserModel := UserModel.SetDatabase(t.Name())

	Convey("Ping the database", t, func() {
		So(UserModel.Ping(), ShouldBeNil)
	})

	Convey("Drop indexes", t, func() {
		Convey("Drop all indexes used by a model", func() {
			UserModel.SyncIndexes()
			So(UserModel.NumberOfIndexes(), ShouldEqual, 2)
			UserModel.DropIndexes()
			So(UserModel.NumberOfIndexes(), ShouldEqual, 1)
		})

		Convey("Drop a specific index used by a model", func() {
			UserModel.SyncIndexes()
			So(UserModel.NumberOfIndexes(), ShouldEqual, 2)
			UserModel.DropIndex("name_1")
			So(UserModel.NumberOfIndexes(), ShouldEqual, 1)
		})
	})

	Convey("Drop collection used by a model", t, func() {
		collections, _ := UserModel.Database().ListCollectionNames(context.TODO(), primitive.M{})
		So(slices.Contains(collections, UserModel.Collection().Name()), ShouldBeTrue)
		UserModel.Drop()
		collections, _ = UserModel.Database().ListCollectionNames(context.TODO(), primitive.M{})
		So(slices.Contains(collections, UserModel.Collection().Name()), ShouldBeFalse)
	})
}
