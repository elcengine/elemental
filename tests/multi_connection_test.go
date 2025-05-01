package e_tests

import (
	"testing"

	e_connection "github.com/elcengine/elemental/connection"
	e_mocks "github.com/elcengine/elemental/tests/mocks"
	e_test_setup "github.com/elcengine/elemental/tests/setup"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMultiConnection(t *testing.T) {
	t.Parallel()

	e_test_setup.Connection(t.Name())

	e_connection.Connect(e_connection.ConnectionOptions{
		URI:   e_mocks.SECONDARY_DATASOURCE,
		Alias: "second",
	})

	MonsterModel := MonsterModel.SetDatabase(t.Name())

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

	Convey("Insert and read monsters from default datasource", t, func() {
		MonsterModel.InsertMany(monstersData).Exec()
		So(len(MonsterModel.Find().Exec().([]Monster)), ShouldEqual, len(monstersData))
	})

	Convey("Insert and read monsters from secondary data source", t, func() {
		MonsterModel.SetConnection("second").InsertMany(monstersData).Exec()
		So(len(MonsterModel.SetConnection("second").Find().Exec().([]Monster)), ShouldEqual, len(monstersData))
	})
}
