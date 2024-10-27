package e_tests

import (
	// "fmt"
	"fmt"

	// "reflect"
	"testing"

	"github.com/elcengine/elemental/connection"
	"github.com/elcengine/elemental/tests/setup"
	"github.com/samber/lo"
	. "github.com/smartystreets/goconvey/convey"
)

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

		for i, monster := range monsters {
			fmt.Printf("Monster %d ID: %s\n", i, monster.ID)
		}

		BestiaryWithIDModel.SetConnection("second")
		beasts := BestiaryWithIDModel.InsertMany([]BestiaryWithID{
			{
				MonsterID: string(monsters[0].ID.String()),
			},
			{
				MonsterID: string(monsters[1].ID.String()),
			},
			{
				MonsterID: string(monsters[2].ID.String()),
			},
		}).Exec().([]BestiaryWithID)

		for i, beast := range beasts {
			fmt.Printf("Beast %d ID: %s\n", i, beast.ID)
		}

		Convey("Populate a with multiple calls", func() {
			bestiary := BestiaryWithIDModel.FindByID(beasts[0].ID).UseCluster(lo.ToPtr("second")).Populate(MonsterModel.FlexibleClone()).Exec()
			if bestiary == nil {
				fmt.Println("Bestiary is nil")
				return
			}
			bestiaryMap, ok := bestiary.(map[string]any)
			if !ok {
				fmt.Println("Bestiary is not a map")
				return
			}
			monster := bestiaryMap["Monster"]
			if monster == nil {
				fmt.Println("Monster is nil")
				return
			}

			So(bestiaryMap, ShouldNotBeNil)
			So(bestiary, ShouldNotBeNil)
		})
	})
}
