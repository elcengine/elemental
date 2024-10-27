package e_tests

import (
	"testing"

	"github.com/elcengine/elemental/connection"
	"github.com/elcengine/elemental/tests/setup"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMultiConnection(t *testing.T) {
	Convey("Read users where", t, func() {
		DB_URI1 := "mongodb+srv://first:dummypass@first.ulhfa.mongodb.net/elemental?retryWrites=true&w=majority&appName=myAtlasClusterEDU"
		DB_URI2 := "mongodb+srv://second:dummypass@cluster0.pqvjtbp.mongodb.net/elemental?retryWrites=true&w=majority&appName=myAtlasClusterEDU"
		e_connection.ConnectURI(DB_URI1)
		e_connection.Connect(e_connection.ConnectionOptions{
			URI:   DB_URI2,
			Alias: "second",
		})

		defer e_test_setup.Teardown()

		monstersData := []Monster{
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
			{
				Name:     "Leshen",
				Category: "Relict",
			},
			{
				Name:     "Fiend",
				Category: "Relict",
			},
			{
				Name:     "Griffin",
				Category: "Hybrid",
			},
			{
				Name:     "Ekimma",
				Category: "Vampire",
			},
			{
				Name:     "Werewolf",
				Category: "Cursed One",
			},
			{
				Name:     "Basilisk",
				Category: "Draconid",
			},
			{
				Name:     "Chort",
				Category: "Relict",
			},
			{
				Name:     "Forktail",
				Category: "Draconid",
			},
			{
				Name:     "Harpie",
				Category: "Hybrid",
			},
			{
				Name:     "Succubus",
				Category: "Relict",
			},
		}

		MonsterModel.InsertMany(monstersData).Exec()
		MonsterModel.SetConnection("second").InsertMany(monstersData).Exec()
	})
}
