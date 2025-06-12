package tests

import (
	"testing"
	"time"

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
		(*castle)["created_at"] = time.Now().AddDate(-1, 0, 0)
		return true
	})

	CastleModel.PostSave(func(castle *bson.M) bool {
		invokedHooks["postSave"] = true
		if name, ok := (*castle)["name"].(string); ok {
			(*castle)["name"] = "Created: " + name
		}
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
		for i := range *castle {
			(*castle)[i].Name = "Modified: " + (*castle)[i].Name
		}
		return true
	})

	CastleModel.PreFindOneAndUpdate(func(filter *primitive.M, doc any) bool {
		invokedHooks["preFindOneAndUpdate"] = true
		return true
	})

	CastleModel.PostFindOneAndUpdate(func(castle *Castle) bool {
		invokedHooks["postFindOneAndUpdate"] = true
		if castle != nil {
			castle.Name = "Updated: " + castle.Name
		}
		return true
	})

	CastleModel.PreFindOneAndDelete(func(filters *primitive.M) bool {
		invokedHooks["preFindOneAndDelete"] = true
		return true
	})

	CastleModel.PostFindOneAndDelete(func(castle *Castle) bool {
		invokedHooks["postFindOneAndDelete"] = true
		if castle != nil {
			castle.Name = "Deleted: " + castle.Name
		}
		return true
	})

	CastleModel.PreFindOneAndReplace(func(castle *primitive.M, doc any) bool {
		invokedHooks["preFindOneAndReplace"] = true
		return true
	})

	CastleModel.PostFindOneAndReplace(func(castle *Castle) bool {
		invokedHooks["postFindOneAndReplace"] = true
		if castle != nil {
			castle.Name = "Replaced: " + castle.Name
		}
		return true
	})

	CastleModel.Create(Castle{Name: "Aretuza"}).Exec()

	CastleModel.Create(Castle{Name: "Rozrog"}).Exec()

	createdCastle := CastleModel.Create(Castle{Name: "Drakenborg"}).ExecT()

	replacedCastle := CastleModel.FindOneAndReplace(&primitive.M{"name": "Drakenborg"}, Castle{Name: "Tesham Mutna"}).ExecPtr()

	CastleModel.UpdateOne(&primitive.M{"name": "Aretuza"}, Castle{Name: "Kaer Morhen"}).Exec()

	updatedCastle := CastleModel.FindOneAndUpdate(&primitive.M{"name": "Rozrog"}, primitive.M{"name": "Rozrog Ruins"}).ExecPtr()

	CastleModel.DeleteOne(primitive.M{"name": "Kaer Morhen"}).Exec()

	deletedCastle := CastleModel.FindOneAndDelete(primitive.M{"name": "Tesham Mutna"}).ExecPtr()

	CastleModel.DeleteMany(primitive.M{"name": primitive.M{"$in": []string{"Aretuza", "Rozrog"}}}).Exec()

	castles := CastleModel.Find().ExecTT()

	Convey("Pre hooks", t, func() {
		Convey("Save", func() {
			So(invokedHooks["preSave"], ShouldBeTrue)
			firstInsertedCastle := CastleModel.FindOne().ExecT()
			So(firstInsertedCastle.CreatedAt.Year(), ShouldEqual, time.Now().AddDate(-1, 0, 0).Year())
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
			So(createdCastle.Name, ShouldEqual, "Created: Drakenborg")
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
			So(castles, ShouldNotBeEmpty)
			for _, castle := range castles {
				So(castle.Name, ShouldStartWith, "Modified: ")
			}
		})
		Convey("FindOneAndUpdate", func() {
			So(invokedHooks["postFindOneAndUpdate"], ShouldBeTrue)
			So(updatedCastle.Name, ShouldEqual, "Updated: Rozrog")
		})
		Convey("FindOneAndDelete", func() {
			So(invokedHooks["postFindOneAndDelete"], ShouldBeTrue)
			So(deletedCastle.Name, ShouldEqual, "Deleted: Tesham Mutna")
		})
		Convey("FindOneAndReplace", func() {
			So(invokedHooks["postFindOneAndReplace"], ShouldBeTrue)
			So(replacedCastle.Name, ShouldEqual, "Replaced: Drakenborg")
		})
	})
}
