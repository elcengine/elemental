package tests

import (
	"os"
	"testing"

	elemental "github.com/elcengine/elemental/core"
	"github.com/elcengine/elemental/tests/fixtures/mocks"
	ts "github.com/elcengine/elemental/tests/fixtures/setup"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMultiConnection(t *testing.T) {
	t.Parallel()

	if os.Getenv("CI") == "" {
		t.Skip("Skipping test in non-CI environment")
	}

	ts.Connection(t.Name())

	elemental.Connect(elemental.ConnectionOptions{
		URI:   mocks.SECONDARY_DATASOURCE,
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
		So(len(MonsterModel.Find().ExecTT()), ShouldEqual, len(monstersData))
	})

	Convey("Insert and read monsters from secondary data source", t, func() {
		MonsterModel.SetConnection("second").InsertMany(monstersData).Exec()
		So(len(MonsterModel.SetConnection("second").Find().ExecTT()), ShouldEqual, len(monstersData))
	})
}
