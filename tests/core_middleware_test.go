package e_tests

import (
	"elemental/core"
	"elemental/tests/setup"
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

	CastleModel.PreUpdateOne(func(doc any) bool {
		invokedHooks["preUpdateOne"] = true
		return true
	})

	CastleModel.PostUpdateOne(func(result *mongo.UpdateResult, err error) bool {
		invokedHooks["postUpdateOne"] = true
		return true
	})

	CastleModel.Create(Castle{Name: "Aretuza"}).Exec()

	CastleModel.UpdateOne(&primitive.M{"name": "Aretuza"}, Castle{Name: "Kaer Morhen"}).Exec()

	Convey("Pre hooks", t, func() {
		Convey("Save", func() {
			So(invokedHooks["preSave"], ShouldBeTrue)
		})
		Convey("UpdateOne", func() {
			So(invokedHooks["preUpdateOne"], ShouldBeTrue)
		})
	})
	Convey("Post hooks", t, func() {
		Convey("Save", func() {
			So(invokedHooks["postSave"], ShouldBeTrue)
		})
		Convey("UpdateOne", func() {
			So(invokedHooks["postUpdateOne"], ShouldBeTrue)
		})
	})
}
