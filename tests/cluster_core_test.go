package e_tests

import (
	// "fmt"
	"fmt"
	"testing"

	e_connection "elemental/connection"
	// elemental "elemental/core"
	// e_test_base "elemental/tests/base"
	e_test_setup "elemental/tests/setup"

	"github.com/samber/lo"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCluster(t *testing.T) {
	Convey("Read users where", t, func() {
		// DB_URI := "mongodb+srv://myAtlasDBUser:dummypass@myatlasclusteredu.edy0gtm.mongodb.net/elemental?retryWrites=true&w=majority&appName=myAtlasClusterEDU"
		DB_URI := "mongodb+srv://second:dummypass@cluster0.pqvjtbp.mongodb.net/elemental?retryWrites=true&w=majority&appName=myAtlasClusterEDU"
		e_connection.ConnectURI(DB_URI)
		e_connection.Connect(e_connection.ConnectionOptions{
			URI:   DB_URI,
			Alias: "second",
		})
		defer e_test_setup.Teardown()
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

		// elemental.UseCluster(BestiaryModel, lo.ToPtr("second"), func(c *elemental.ClusterOp[e_test_base.Bestiary]) {
		// 	c.Populate("second", []string{"monster", "kingdom"})
		// }).Exec()
		// (lo.ToPtr("second"),
		// 	func(c *elemental.ClusterOp[e_test_base.Bestiary]) {
		// 		c.Populate("second", )
		// 	},
		// ).Exec()

		Convey("Populate a with multiple calls", func() {
			bestiary := BestiaryModel.UseCluster(lo.ToPtr("second")).Exec().([]any)
			So(bestiary, ShouldHaveLength, 3)
			So(bestiary[0].(map[string]any)["monster"].(map[string]any)["Name"], ShouldEqual, "Katakan")
			// So(bestiary[0].Kingdom.Name, ShouldEqual, "Nilfgaard")
		})

		Convey("Populate a with multiple calls", func() {
			bestiary := BestiaryModel.UseCluster(lo.ToPtr("second")).Exec().([]Bestiary)
			So(bestiary, ShouldHaveLength, 3)
			So(bestiary[0].Monster.Name, ShouldEqual, "Katakan")
			So(bestiary[0].Kingdom.Name, ShouldEqual, "Nilfgaard")
		})
		Convey("Populate with a single call", func() {
			bestiary := BestiaryModel.UseCluster(lo.ToPtr("second")).Exec().([]Bestiary)
			So(bestiary, ShouldHaveLength, 3)
			So(bestiary[0].Monster.Name, ShouldEqual, "Katakan")
			So(bestiary[0].Kingdom.Name, ShouldEqual, "Nilfgaard")
		})
		Convey("Populate with a single call (Comma separated string)", func() {
			bestiary := BestiaryModel.UseCluster(lo.ToPtr("second")).Exec().([]Bestiary)
			So(bestiary, ShouldHaveLength, 3)
			So(bestiary[0].Monster.Name, ShouldEqual, "Katakan")
			So(bestiary[0].Kingdom.Name, ShouldEqual, "Nilfgaard")
		})
	})
}

func TestClusterWithID(t *testing.T) {
	Convey("Read users where", t, func() {
		DB_URI := "mongodb+srv://second:dummypass@cluster0.pqvjtbp.mongodb.net/elemental?retryWrites=true&w=majority&appName=myAtlasClusterEDU"
		e_connection.ConnectURI(DB_URI)
		e_connection.Connect(e_connection.ConnectionOptions{
			URI:   DB_URI,
			Alias: "second",
		})
		defer e_test_setup.Teardown()

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

		BestiaryWithIDModel.SetConnection("second")
		BestiaryWithIDModel.InsertMany([]BestiaryWithID{
			{
				MonsterID: string(monsters[0].ID.String()),
			},
			{
				MonsterID: string(monsters[1].ID.String()),
			},
			{
				MonsterID: string(monsters[2].ID.String()),
			},
		}).Exec()

		// elemental.UseCluster(BestiaryModel, lo.ToPtr("second"), func(c *elemental.ClusterOp[e_test_base.Bestiary]) {
		// 	c.Populate("second", []string{"monster", "kingdom"})
		// }).Exec()
		// (lo.ToPtr("second"),
		// 	func(c *elemental.ClusterOp[e_test_base.Bestiary]) {
		// 		c.Populate("second", )
		// 	},
		// ).Exec()

		Convey("Populate a with multiple calls", func() {
			bestiary := BestiaryWithIDModel.UseCluster(lo.ToPtr("second")).PopulateOp(MonsterModel).Exec().([]BestiaryWithID)
			So(bestiary, ShouldHaveLength, 3)
			fmt.Printf("\nMonsterID: %s\n",bestiary[0].MonsterID)
			fmt.Printf("\nMonsterID: %s\n",bestiary[1].MonsterID)
			fmt.Printf("\nMonsterID: %s\n",bestiary[2].MonsterID)
			// So(bestiary[0].MonsterID, ShouldEqual, "Katakan")
		})
	})
}
