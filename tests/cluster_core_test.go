package e_tests

import (
	// "fmt"
	"testing"

	"elemental/connection"
	"elemental/core"
	"elemental/tests/base"
	"elemental/tests/setup"

	"github.com/samber/lo"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCluster(t *testing.T) {
	Convey("Read users where", t, func() {
		DB_URI := "mongodb+srv://myAtlasDBUser:u8CFO9C9ALu4MUAZ@myatlasclusteredu.edy0gtm.mongodb.net/elemental?retryWrites=true&w=majority&appName=myAtlasClusterEDU"
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
			bestiary := BestiaryModel.UseCluster(lo.ToPtr("second"), func(c *elemental.ClusterOp[e_test_base.Bestiary]) {
				c.Populate("second", []string{"monster", "kingdom"})
			}).Exec().([]Bestiary)
			So(bestiary, ShouldHaveLength, 3)
			So(bestiary[0].Monster.Name, ShouldEqual, "Katakan")
			So(bestiary[0].Kingdom.Name, ShouldEqual, "Nilfgaard")
		})
		Convey("Populate with a single call", func() {
			bestiary := BestiaryModel.UseCluster(lo.ToPtr("second"), func(c *elemental.ClusterOp[e_test_base.Bestiary]) {
				c.Populate("second", []string{"monster", "kingdom"})
			}).Exec().([]Bestiary)
			So(bestiary, ShouldHaveLength, 3)
			So(bestiary[0].Monster.Name, ShouldEqual, "Katakan")
			So(bestiary[0].Kingdom.Name, ShouldEqual, "Nilfgaard")
		})
		Convey("Populate with a single call (Comma separated string)", func() {
			bestiary := BestiaryModel.UseCluster(lo.ToPtr("second"), func(c *elemental.ClusterOp[e_test_base.Bestiary]) {
				c.Populate("second", []string{"monster", "kingdom"})
			}).Exec().([]Bestiary)
			So(bestiary, ShouldHaveLength, 3)
			So(bestiary[0].Monster.Name, ShouldEqual, "Katakan")
			So(bestiary[0].Kingdom.Name, ShouldEqual, "Nilfgaard")
		})
	})
}
