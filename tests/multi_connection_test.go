package e_tests

import (
	"testing"

	"github.com/elcengine/elemental/connection"
	"github.com/elcengine/elemental/tests/mocks"
	"github.com/elcengine/elemental/tests/setup"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMultiConnection(t *testing.T) {
	Convey("Read users where", t, func() {
		e_connection.ConnectURI(e_mocks.DEFAULT_DB_URI)
		e_connection.Connect(e_connection.ConnectionOptions{
			URI:   e_mocks.SECONDARY_DB_URI,
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

		So(len(MonsterModel.Find().Exec().([]Monster)), ShouldEqual, len(monstersData))

		So(len(MonsterModel.SetConnection("second").Find().Exec().([]Monster)), ShouldEqual, len(monstersData))
	})
}
