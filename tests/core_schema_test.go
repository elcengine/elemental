package e_tests

import (
	"elemental/core"
	"elemental/tests/mocks"
	"elemental/tests/setup"
	"fmt"
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCoreSchemaOptions(t *testing.T) {
	e_test_setup.Connection()
	defer e_test_setup.Teardown()

	Convey("Schema variations", t, func() {
		Convey(fmt.Sprintf("Should use the default database of %s", e_mocks.DEFAULT_DB), func() {
			So(UserModel.Collection().Database().Name(), ShouldEqual, e_mocks.DEFAULT_DB)
		})
		Convey("Should automatically create a collection with given name", func() {
			So(UserModel.Collection().Name(), ShouldEqual, "users")
		})
		Convey("Collection should be a plural of the model name if not specified", func() {
			var CastleModel = elemental.NewModel[Castle]("Castle", elemental.NewSchema(map[string]elemental.Field{
				"Name": {
					Type:     reflect.String,
					Required: true,
				},
			}, elemental.SchemaOptions{}))
			CastleModel.Create(Castle{Name: "Kaer Morhen"})
			So(CastleModel.Collection().Name(), ShouldEqual, "castles")
		})
		Convey("Should use the specified database", func() {
			var MonsterModel = elemental.NewModel[Monster]("Monster", elemental.NewSchema(map[string]elemental.Field{}, elemental.SchemaOptions{
				Database: e_mocks.SECONDARY_DB,
			}))
			MonsterModel.Create(Monster{Name: "Nekker"})
			So(MonsterModel.Collection().Database().Name(), ShouldEqual, e_mocks.SECONDARY_DB)
		})
	})
}
