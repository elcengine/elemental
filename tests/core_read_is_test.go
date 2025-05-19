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

			users = UserModel.Where("name").IsString().ExecTT()
			So(len(users), ShouldEqual, len(mocks.Users))
		})
		Convey("School is null", func() {
			users := UserModel.Where("school").IsNull().ExecTT()
			So(len(users), ShouldEqual, len(lo.Filter(mocks.Users, func(u User, _ int) bool {
				return u.School == nil
			})))
		})
		Convey("Retired is boolean", func() {
			users := UserModel.Where("retired").IsBoolean().ExecTT()
			So(len(users), ShouldEqual, len(mocks.Users))
		})
		Convey("Age is int32", func() {
			users := UserModel.Where("age").IsInt32().ExecTT()
			So(len(users), ShouldEqual, len(mocks.Users))
		})
		Convey("Age is int64", func() {
			users := UserModel.Where("age").IsInt64().ExecTT()
			So(len(users), ShouldEqual, 0)
		})
		Convey("Age is double", func() {
			users := UserModel.Where("age").IsDouble().ExecTT()
			So(len(users), ShouldEqual, 0)
		})
		Convey("Weapons is arrayw", func() {
			users := UserModel.Where("weapons").IsArray().ExecTT()
			So(len(users), ShouldEqual, len(mocks.Users))
		})
		Convey("Name is binary", func() {
			users := UserModel.Where("name").IsBinary().ExecTT()
			So(len(users), ShouldEqual, 0)
		})
		Convey("Name is undefined", func() {
			users := UserModel.Where("name").IsUndefined().ExecTT()
			So(len(users), ShouldEqual, 0)
		})
		Convey("ID is object id", func() {
			users := UserModel.Where("_id").IsObjectID().ExecTT()
			So(len(users), ShouldEqual, len(mocks.Users))
		})
		Convey("Created At is date", func() {
			users := UserModel.Where("created_at").IsDateTime().ExecTT()
			So(len(users), ShouldEqual, len(mocks.Users))
		})
		Convey("Name is regex", func() {
			users := UserModel.Where("name").IsRegex().ExecTT()
			So(len(users), ShouldEqual, 0)
		})
		Convey("Name is db pointer", func() {
			users := UserModel.Where("name").IsDBPointer().ExecTT()
			So(len(users), ShouldEqual, 0)
		})
		Convey("Name is javascript", func() {
			users := UserModel.Where("name").IsJavaScript().ExecTT()
			So(len(users), ShouldEqual, 0)
		})
		Convey("Name is symbol", func() {
			users := UserModel.Where("name").IsSymbol().ExecTT()
			So(len(users), ShouldEqual, 0)
		})
		Convey("Name is code with scope", func() {
			users := UserModel.Where("name").IsCodeWithScope().ExecTT()
			So(len(users), ShouldEqual, 0)
		})
		Convey("Updated at is timestamp", func() {
			users := UserModel.Where("updated_at").IsTimestamp().ExecTT()
			So(len(users), ShouldEqual, 0)
		})
		Convey("Name is decimal128", func() {
			users := UserModel.Where("name").IsDecimal128().ExecTT()
			So(len(users), ShouldEqual, 0)
		})
	})
}
