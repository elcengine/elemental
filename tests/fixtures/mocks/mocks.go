// Sample data to be used within test suites.
package mocks

import (
	"os"

	"github.com/elcengine/elemental/tests/fixtures"
	"github.com/samber/lo"
)

var (
	DEFAULT_DATASOURCE = lo.CoalesceOrEmpty(os.Getenv("DEFAULT_DATASOURCE"),
		"mongodb+srv://akalankaperera128:pFAnQVXE6vrbcXNk@default.ynr156r.mongodb.net/elemental") // This is a test M0 cluster there, so it is safe to use in tests.
	SECONDARY_DATASOURCE = os.Getenv("SECONDARY_DATASOURCE")
)

var (
	DEFAULT_DB_NAME = "elemental"
)

var (
	WolfSchool      = "Wolf"
	BearSchool      = "Bear"
	GriffinSchool   = "Griffin"
	ManticoreSchool = "Manticore"
)

var (
	Ciri = fixtures.User{
		Name: "Ciri",
	}
	Geralt = fixtures.User{
		Name:       "Geralt",
		Age:        100,
		Occupation: "Witcher",
		Weapons:    []string{"Silver sword", "Mahakaman battle hammer", "Battle Axe", "Crossbow", "Steel sword"},
		School:     &WolfSchool,
	}
	Eredin = fixtures.User{
		Name: "Eredin",
	}
	Caranthir = fixtures.User{
		Name:       "Caranthir",
		Age:        120,
		Occupation: "Mage",
		Weapons:    []string{"Staff"},
	}
	Imlerith = fixtures.User{
		Name:       "Imlerith",
		Age:        150,
		Occupation: "General",
		Weapons:    []string{"Mace", "Battle Axe"},
	}
	Yennefer = fixtures.User{
		Name:       "Yennefer",
		Occupation: "Mage",
		Age:        100,
	}
	Vesemir = fixtures.User{
		Name:       "Vesemir",
		Occupation: "Witcher",
		Age:        300,
		Weapons:    []string{"Silver sword", "Steel sword", "Crossbow"},
		Retired:    true,
		School:     &WolfSchool,
	}
)

var Users = []fixtures.User{
	Ciri,
	Geralt,
	Eredin,
	Caranthir,
	Imlerith,
	Yennefer,
	Vesemir,
}
