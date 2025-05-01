package e_tests

import (
	"testing"

	e_test_setup "github.com/elcengine/elemental/tests/setup"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCoreReadPopulate(t *testing.T) {
	t.Parallel()

	e_test_setup.Connection()

	defer e_test_setup.Teardown()

	var LocalMonsterModel = MonsterModel.Clone().SetCollection("monsters_for_populate")
	var LocalKingdomModel = KingdomModel.Clone().SetCollection("kingdoms_for_populate")
	var LocalBestiaryModel = BestiaryModel.Clone().SetCollection("bestiaries_for_populate")

	monsters := LocalMonsterModel.InsertMany([]Monster{
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

	kingdoms := LocalKingdomModel.InsertMany([]Kingdom{
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

	LocalBestiaryModel.InsertMany([]Bestiary{
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
			bestiary := LocalBestiaryModel.Find().Populate("monster").Populate("kingdom").Exec().([]Bestiary)
			So(bestiary, ShouldHaveLength, 3)
			So(bestiary[0].Monster.Name, ShouldEqual, "Katakan")
			So(bestiary[0].Kingdom.Name, ShouldEqual, "Nilfgaard")
		})
		Convey("Populate with a single call", func() {
			bestiary := LocalBestiaryModel.Find().Populate("monster", "kingdom").Exec().([]Bestiary)
			So(bestiary, ShouldHaveLength, 3)
			So(bestiary[0].Monster.Name, ShouldEqual, "Katakan")
			So(bestiary[0].Kingdom.Name, ShouldEqual, "Nilfgaard")
		})
		Convey("Populate with a single call (Comma separated string)", func() {
			bestiary := LocalBestiaryModel.Find().Populate("monster kingdom").Exec().([]Bestiary)
			So(bestiary, ShouldHaveLength, 3)
			So(bestiary[0].Monster.Name, ShouldEqual, "Katakan")
			So(bestiary[0].Kingdom.Name, ShouldEqual, "Nilfgaard")
		})
	})
}
