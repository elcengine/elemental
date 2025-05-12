package e_test_setup

import (
	"strings"

	elemental "github.com/elcengine/elemental/core"
	e_test_base "github.com/elcengine/elemental/tests/base"
	e_mocks "github.com/elcengine/elemental/tests/mocks"
)

func Connection(databaseName string) {
	elemental.Connect(strings.Replace(e_mocks.DEFAULT_DATASOURCE, e_mocks.DEFAULT_DB_NAME, databaseName, 1))
}

func Seed(databaseName string) {
	e_test_base.UserModel.SetDatabase(databaseName).InsertMany(e_mocks.Users).Exec()
}

func SeededConnection(databaseName string) {
	Connection(databaseName)
	Seed(databaseName)
}

func Teardown() {
	elemental.DropAllDatabases()
	elemental.Disconnect()
}
