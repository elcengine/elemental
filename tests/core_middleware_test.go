package e_tests

import (
	"elemental/core"
	"elemental/tests/setup"
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCoreMiddleware(t *testing.T) {
	e_test_setup.Connection()
	defer e_test_setup.Teardown()

	invokedHooks := make(map[string]bool)

	var CastleModel = elemental.NewModel[Castle]("Castle-For-Middleware", elemental.NewSchema(map[string]elemental.Field{
		"Name": {
			Type:     reflect.String,
			Required: true,
		},
	}))

	CastleModel.PreSave(func(castle Castle) bool {
		invokedHooks["preSave"] = true
		return true
	})

	CastleModel.PostSave(func(castle Castle) bool {
		invokedHooks["postSave"] = true
		return true
	})

	CastleModel.Create(Castle{Name: "Aretuza"})

	Convey("Pre hooks", t, func() {
		Convey("Save", func() {
			So(invokedHooks["preSave"], ShouldBeTrue)
		})
	})
	Convey("Post hooks", t, func() {
		Convey("Save", func() {
			So(invokedHooks["postSave"], ShouldBeTrue)
		})
	})
}
