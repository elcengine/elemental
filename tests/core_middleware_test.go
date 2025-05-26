package tests

import (
	"testing"

	elemental "github.com/elcengine/elemental/core"
	ts "github.com/elcengine/elemental/tests/fixtures/setup"

	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestCoreMiddleware(t *testing.T) {
	t.Parallel()

	ts.Connection(t.Name())

	invokedHooks := make(map[string]bool)

	CastleModel := elemental.NewModel[Castle]("Castle-For-Middleware", elemental.NewSchema(map[string]elemental.Field{
		"Name": {
			Type:     elemental.String,
			Required: true,
		},
	})).SetDatabase(t.Name())

	CastleModel.PreSave(func(castle *bson.M) bool {
		invokedHooks["preSave"] = true
		return true
	})

	CastleModel.PostSave(func(castle *bson.M) bool {
		invokedHooks["postSave"] = true
		return true
	})

	CastleModel.PostSave(func(castle *bson.M) bool {
		invokedHooks["postSaveSecond"] = true
		return false
	})

	CastleModel.PostSave(func(castle *bson.M) bool {
		invokedHooks["postSaveThird"] = true
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

	CastleModel.PreDeleteOne(func(filters *primitive.M) bool {
		invokedHooks["preDeleteOne"] = true
		return true
	})

	CastleModel.PostDeleteOne(func(result *mongo.DeleteResult, err error) bool {
		invokedHooks["postDeleteOne"] = true
		return true
	})

	CastleModel.PreDeleteMany(func(filters *primitive.M) bool {
		invokedHooks["preDeleteMany"] = true
		return true
	})

	CastleModel.PostDeleteMany(func(result *mongo.DeleteResult, err error) bool {
		invokedHooks["postDeleteMany"] = true
		return true
	})

	CastleModel.PostFind(func(castle *[]Castle) bool {
		invokedHooks["postFind"] = true
		return true
	})

	CastleModel.PreFindOneAndUpdate(func(filter *primitive.M, doc any) bool {
		invokedHooks["preFindOneAndUpdate"] = true
		return true
	})

	CastleModel.PostFindOneAndUpdate(func(castle *Castle) bool {
		invokedHooks["postFindOneAndUpdate"] = true
		return true
	})

	CastleModel.PreFindOneAndDelete(func(filters *primitive.M) bool {
		invokedHooks["preFindOneAndDelete"] = true
		return true
	})

	CastleModel.PostFindOneAndDelete(func(castle *Castle) bool {
		invokedHooks["postFindOneAndDelete"] = true
		return true
	})

	CastleModel.PreFindOneAndReplace(func(castle *primitive.M, doc any) bool {
		invokedHooks["preFindOneAndReplace"] = true
		return true
	})

	CastleModel.PostFindOneAndReplace(func(castle *Castle) bool {
		invokedHooks["postFindOneAndReplace"] = true
		return true
	})

	CastleModel.Create(Castle{Name: "Aretuza"}).Exec()

	CastleModel.Create(Castle{Name: "Maverick"}).Exec()

	CastleModel.Create(Castle{Name: "Robert"}).Exec()

	CastleModel.FindOneAndReplace(&primitive.M{"name": "Robert"}, Castle{Name: "Jack"}).Exec()

	CastleModel.UpdateOne(&primitive.M{"name": "Aretuza"}, Castle{Name: "Kaer Morhen"}).Exec()

	CastleModel.Find().Exec()

	CastleModel.FindOneAndUpdate(&primitive.M{"name": "Maverick"}, primitive.M{"name": "Maverickk"}).Exec()

	CastleModel.DeleteOne(primitive.M{"name": "Kaer Morhen"}).Exec()

	CastleModel.FindOneAndDelete(primitive.M{"name": "Jack"}).Exec()

	CastleModel.DeleteMany(primitive.M{"name": primitive.M{"$in": []string{"Aretuza", "Maverick"}}}).Exec()

	Convey("Pre hooks", t, func() {
		Convey("Save", func() {
			So(invokedHooks["preSave"], ShouldBeTrue)
		})
		Convey("UpdateOne", func() {
			So(invokedHooks["preUpdateOne"], ShouldBeTrue)
		})
		Convey("DeleteOne", func() {
			So(invokedHooks["preDeleteOne"], ShouldBeTrue)
		})
		Convey("DeleteMany", func() {
			So(invokedHooks["preDeleteMany"], ShouldBeTrue)
		})
		Convey("FindOneAndUpdate", func() {
			So(invokedHooks["preFindOneAndUpdate"], ShouldBeTrue)
		})
		Convey("FindOneAndDelete", func() {
			So(invokedHooks["preFindOneAndDelete"], ShouldBeTrue)
		})
		Convey("FindOneAndReplace", func() {
			So(invokedHooks["preFindOneAndReplace"], ShouldBeTrue)
		})
	})

	Convey("Post hooks", t, func() {
		Convey("Save", func() {
			So(invokedHooks["postSave"], ShouldBeTrue)
			So(invokedHooks["postSaveSecond"], ShouldBeTrue)
			So(invokedHooks["postSaveThird"], ShouldBeFalse)
		})
		Convey("UpdateOne", func() {
			So(invokedHooks["postUpdateOne"], ShouldBeTrue)
		})
		Convey("DeleteOne", func() {
			So(invokedHooks["postDeleteOne"], ShouldBeTrue)
		})
		Convey("DeleteMany", func() {
			So(invokedHooks["postDeleteMany"], ShouldBeTrue)
		})
		Convey("Find", func() {
			So(invokedHooks["postFind"], ShouldBeTrue)
		})
		Convey("FindOneAndUpdate", func() {
			So(invokedHooks["postFindOneAndUpdate"], ShouldBeTrue)
		})
		Convey("FindOneAndDelete", func() {
			So(invokedHooks["postFindOneAndDelete"], ShouldBeTrue)
		})
		Convey("FindOneAndReplace", func() {
			So(invokedHooks["postFindOneAndReplace"], ShouldBeTrue)
		})
	})
}
