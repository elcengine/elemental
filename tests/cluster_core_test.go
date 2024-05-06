package e_tests

import (
	"elemental/connection"
	"elemental/tests/setup"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCoreReadPopulateOnMultipleClusters(t *testing.T) {
	DB_URI1 := "mongodb+srv://myAtlasDBUser:dummypass@myatlasclusteredu.edy0gtm.mongodb.net/elemental"
	DB_URI2 := "mongodb+srv://second:dummypass@cluster0.pqvjtbp.mongodb.net/elemental"

	e_connection.ConnectURI(DB_URI1)
	e_connection.Connect(e_connection.ConnectionOptions{
		URI: DB_URI2,
		Alias: "secondary",
	})
	seed()

	defer e_test_setup.Teardown()
	

	Convey("Find with populated fields on multiple clusters", t, func() {
		Convey("Populate a with multiple calls", func() {
			bestiary := BestiaryModel.UseCluster("secondary").Find().Populate("monster").Populate("kingdom").Exec().([]Bestiary)
			So(bestiary, ShouldHaveLength, 3)
			So(bestiary[0].Monster.Name, ShouldEqual, "Katakan")
			So(bestiary[0].Kingdom.Name, ShouldEqual, "Nilfgaard")
		})
		Convey("Populate with a single call on multiple clusters", func() {
			bestiary := BestiaryModel.UseCluster("secondary").Find().Populate("monster", "kingdom").Exec().([]Bestiary)
			So(bestiary, ShouldHaveLength, 3)
			So(bestiary[0].Monster.Name, ShouldEqual, "Katakan")
			So(bestiary[0].Kingdom.Name, ShouldEqual, "Nilfgaard")
		})
		Convey("Populate with a single call (Comma separated string)", func() {
			bestiary := BestiaryModel.UseCluster("secondary").Find().Populate("monster kingdom").Exec().([]Bestiary)
			So(bestiary, ShouldHaveLength, 3)
			So(bestiary[0].Monster.Name, ShouldEqual, "Katakan")
			So(bestiary[0].Kingdom.Name, ShouldEqual, "Nilfgaard")
		})
	})
}

func seed() {
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
}

func seedWithAlias(Alias string) {
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
}