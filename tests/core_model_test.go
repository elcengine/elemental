package tests

import (
	"reflect"
	"testing"

	elemental "github.com/elcengine/elemental/core"
	ts "github.com/elcengine/elemental/tests/fixtures/setup"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCoreModel(t *testing.T) {
	t.Parallel()

	ts.Connection(t.Name())

	UserModel := UserModel.SetDatabase(t.Name())

	Convey("Should clone model", t, func() {
		ClonedModel := UserModel.Clone()
		So(ClonedModel.Name, ShouldEqual, UserModel.Name)
		So(ClonedModel.Cloned, ShouldBeTrue)
	})

	Convey("Should use cached model", t, func() {
		CachedUserModel := elemental.NewModel[User]("User", elemental.NewSchema(map[string]elemental.Field{
			"Name": {
				Type: reflect.String,
			},
		}))
		So(CachedUserModel.Schema.Definitions, ShouldHaveLength, len(UserModel.Schema.Definitions))
	})
}
