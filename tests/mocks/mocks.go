package e_mocks

import (
	"os"

	e_test_base "github.com/elcengine/elemental/tests/base"
	e_utils "github.com/elcengine/elemental/utils"
)

var (
	DEFAULT_DATASOURCE   = e_utils.Coalesce(os.Getenv("DEFAULT_DATASOURCE"), "mongodb+srv://akalankaperera128:pFAnQVXE6vrbcXNk@default.ynr156r.mongodb.net/elemental") // This is a test M0 cluster there, so it is safe to use in tests.
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
	Ciri = e_test_base.User{
		Name: "Ciri",
	}
	Geralt = e_test_base.User{
		Name:       "Geralt",
		Age:        100,
		Occupation: "Witcher",
		Weapons:    []string{"Silver sword", "Mahakaman battle hammer", "Battle Axe", "Crossbow", "Steel sword"},
		School:     &WolfSchool,
	}
	Eredin = e_test_base.User{
		Name: "Eredin",
	}
	Caranthir = e_test_base.User{
		Name:       "Caranthir",
		Age:        120,
		Occupation: "Mage",
		Weapons:    []string{"Staff"},
	}
	Imlerith = e_test_base.User{
		Name:       "Imlerith",
		Age:        150,
		Occupation: "General",
		Weapons:    []string{"Mace", "Battle Axe"},
	}
	Yennefer = e_test_base.User{
		Name:       "Yennefer",
		Occupation: "Mage",
		Age:        100,
	}
	Vesemir = e_test_base.User{
		Name:       "Vesemir",
		Occupation: "Witcher",
		Age:        300,
		Weapons:    []string{"Silver sword", "Steel sword", "Crossbow"},
		Retired:    true,
		School:     &WolfSchool,
	}
)

var Users = []e_test_base.User{
	Ciri,
	Geralt,
	Eredin,
	Caranthir,
	Imlerith,
	Yennefer,
	Vesemir,
}
