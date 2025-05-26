package tests

import (
	"testing"

	ts "github.com/elcengine/elemental/tests/fixtures/setup"
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
	}).Exec().([]Monster)

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
	}).Exec().([]Kingdom)

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

	Convey("Find with populated fields", t, func() {
		Convey("Populate a with multiple calls", func() {
			bestiary := BestiaryModel.Find().Populate("monster").Populate("kingdom").Exec().([]DetailedBestiary)
			So(bestiary, ShouldHaveLength, 3)
			So(bestiary[0].Monster.Name, ShouldEqual, "Katakan")
			So(bestiary[0].Monster.Category, ShouldEqual, "Vampire")
			So(bestiary[0].Kingdom.Name, ShouldEqual, "Nilfgaard")
		})
		Convey("Populate with a single call", func() {
			bestiary := BestiaryModel.Find().Populate("monster", "kingdom").Exec().([]DetailedBestiary)
			So(bestiary, ShouldHaveLength, 3)
			So(bestiary[0].Monster.Name, ShouldEqual, "Katakan")
			So(bestiary[0].Monster.Category, ShouldEqual, "Vampire")
			So(bestiary[0].Kingdom.Name, ShouldEqual, "Nilfgaard")
		})
		Convey("Populate with a single call (Comma separated string)", func() {
			bestiary := BestiaryModel.Find().Populate("monster kingdom").Exec().([]DetailedBestiary)
			So(bestiary, ShouldHaveLength, 3)
			So(bestiary[0].Monster.Name, ShouldEqual, "Katakan")
			So(bestiary[0].Monster.Category, ShouldEqual, "Vampire")
			So(bestiary[0].Kingdom.Name, ShouldEqual, "Nilfgaard")
		})
		Convey("Populate with select", func() {
			bestiary := BestiaryModel.Find().Populate(primitive.M{
				"path":   "monster",
				"select": primitive.M{"name": 1},
			}, "kingdom").Exec().([]DetailedBestiary)
			So(bestiary, ShouldHaveLength, 3)
			So(bestiary[0].Monster.Name, ShouldEqual, "Katakan")
			So(bestiary[0].Monster.Category, ShouldEqual, "")
			So(bestiary[0].Kingdom.Name, ShouldEqual, "Nilfgaard")
		})
	})
}
