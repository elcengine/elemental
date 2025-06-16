package tests

import (
	"errors"
	"testing"

	elemental "github.com/elcengine/elemental/core"
	ts "github.com/elcengine/elemental/tests/fixtures/setup"
	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCoreReadPopulate(t *testing.T) {
	t.Parallel()

	ts.Connection(t.Name())

	MonsterModel := MonsterModel.SetDatabase(t.Name())
	KingdomModel := KingdomModel.SetDatabase(t.Name())
	BestiaryModel := BestiaryModel.SetDatabase(t.Name())

	monsters := MonsterModel.InsertMany([]Monster{
		{
			Name:     "Katakan",
			Category: "Vampire",
		},
		{
			Name:     "Drowner",
			Category: "Drowner",
		},
		{
			Name:     "Nekker",
			Category: "Nekker",
		},
	}).ExecTT()

	kingdoms := KingdomModel.InsertMany([]Kingdom{
		{
			Name: "Nilfgaard",
		},
		{
			Name: "Redania",
		},
		{
			Name: "Skellige",
		},
	}).ExecTT()

	BestiaryModel.InsertMany([]Bestiary{
		{
			Monster: monsters[0],
			Kingdom: kingdoms[0],
		},
		{
			Monster: monsters[1],
			Kingdom: kingdoms[1],
		},
		{
			Monster: monsters[2],
			Kingdom: kingdoms[2],
		},
	}).Exec()

	SoKatakan := func(bestiary DetailedBestiary) {
		t.Helper()
		So(bestiary.Monster.Name, ShouldEqual, "Katakan")
		So(bestiary.Monster.Category, ShouldEqual, "Vampire")
		So(bestiary.Kingdom.Name, ShouldEqual, "Nilfgaard")
		So(bestiary.Kingdom.ID, ShouldEqual, kingdoms[0].ID)
	}

	SoDrowner := func(bestiary DetailedBestiary) {
		t.Helper()
		So(bestiary.Monster.Name, ShouldEqual, "Drowner")
		So(bestiary.Monster.Category, ShouldEqual, "Drowner")
		So(bestiary.Kingdom.Name, ShouldEqual, "Redania")
		So(bestiary.Kingdom.ID, ShouldEqual, kingdoms[1].ID)
	}

	Convey("Find with populated fields", t, func() {
		Convey("Populate with multiple calls", func() {
			Convey("With bson field name", func() {
				bestiaries := make([]DetailedBestiary, 0)
				BestiaryModel.Find().Populate("monster").Populate("kingdom").ExecInto(&bestiaries)
				So(bestiaries, ShouldHaveLength, 3)
				SoKatakan(bestiaries[0])
				SoDrowner(bestiaries[1])
			})
			Convey("With model field name", func() {
				bestiaries := make([]DetailedBestiary, 0)
				BestiaryModel.Find().Populate("Monster").Populate("Kingdom").ExecInto(&bestiaries)
				So(bestiaries, ShouldHaveLength, 3)
				SoKatakan(bestiaries[0])
				SoDrowner(bestiaries[1])
			})
			Convey("With OrFail", func() {
				bestiaries := make([]DetailedBestiary, 0)
				So(func() {
					BestiaryModel.Find(primitive.M{"monster": uuid.New()}).Populate("monster").Populate("kingdom").OrFail().ExecInto(&bestiaries)
				}, ShouldPanicWith, errors.New("no results found matching the given query"))
			})
		})
		Convey("Populate with a single call", func() {
			Convey("With bson field name", func() {
				var bestiaries []DetailedBestiary
				BestiaryModel.Find().Populate("monster", "kingdom").ExecInto(&bestiaries)
				So(bestiaries, ShouldHaveLength, 3)
				SoKatakan(bestiaries[0])
				SoDrowner(bestiaries[1])
			})
			Convey("With model field name", func() {
				var bestiaries []DetailedBestiary
				BestiaryModel.Find().Populate("Monster", "Kingdom").ExecInto(&bestiaries)
				So(bestiaries, ShouldHaveLength, 3)
				SoKatakan(bestiaries[0])
				SoDrowner(bestiaries[1])
			})
		})
		Convey("Populate with a single call (Space separated string)", func() {
			var bestiaries []DetailedBestiary
			BestiaryModel.Find().Populate("monster kingdom").ExecInto(&bestiaries)
			So(bestiaries, ShouldHaveLength, 3)
			SoKatakan(bestiaries[0])
			SoDrowner(bestiaries[1])
		})
		Convey("Populate with select", func() {
			var bestiaries []DetailedBestiary
			BestiaryModel.Find().Populate(primitive.M{
				"path":   "monster",
				"select": primitive.M{"name": 1},
			}, "kingdom").ExecInto(&bestiaries)
			So(bestiaries, ShouldHaveLength, 3)
			So(bestiaries[0].Monster.Name, ShouldEqual, "Katakan")
			So(bestiaries[0].Monster.Category, ShouldEqual, "")
			So(bestiaries[0].Kingdom.Name, ShouldEqual, "Nilfgaard")
		})
		Convey("Populate with a pipeline", func() {
			var bestiaries []DetailedBestiary
			BestiaryModel.Find().Populate("kingdom", primitive.M{
				"path": "monster",
				"pipeline": []primitive.M{
					{"$addFields": primitive.M{"name": primitive.M{"$concat": []string{"Rare ", "$name"}}}},
				},
			}).ExecInto(&bestiaries)
			So(bestiaries, ShouldHaveLength, 3)
			So(bestiaries[0].Monster.Name, ShouldEqual, "Rare Katakan")
			So(bestiaries[0].Monster.Category, ShouldEqual, "Vampire")
			So(bestiaries[0].Kingdom.Name, ShouldEqual, "Nilfgaard")
		})
		Convey("Populate just one field", func() {
			var bestiaries []GenericBestiary[Monster, primitive.ObjectID]
			BestiaryModel.Find().Populate("monster").ExecInto(&bestiaries)
			So(bestiaries, ShouldHaveLength, 3)
			So(bestiaries[0].Monster.Name, ShouldEqual, "Katakan")
			So(bestiaries[0].Monster.Category, ShouldEqual, "Vampire")
			So(bestiaries[0].Kingdom, ShouldEqual, kingdoms[0].ID)
		})
		Convey("Populate a model with just collection references", func() {
			BestiaryModel := elemental.NewModel[Bestiary](uuid.NewString(), elemental.NewSchema(map[string]elemental.Field{
				"Monster": {
					Type:       elemental.ObjectID,
					Collection: "monsters",
				},
				"Kingdom": {
					Type:       elemental.ObjectID,
					Collection: "kingdoms",
				},
			}, elemental.SchemaOptions{
				Collection: "bestiary",
			})).SetDatabase(t.Name())
			var bestiaries []DetailedBestiary
			BestiaryModel.Find().Populate("monster", "kingdom").ExecInto(&bestiaries)
			So(bestiaries, ShouldHaveLength, 3)
			SoKatakan(bestiaries[0])
			SoDrowner(bestiaries[1])
		})
	})
}
