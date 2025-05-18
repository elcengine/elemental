// Initialization and teardown functions for all test suites.
package ts

import (
	"strings"

	elemental "github.com/elcengine/elemental/core"
	"github.com/elcengine/elemental/tests/fixtures"
	"github.com/elcengine/elemental/tests/fixtures/mocks"
)

func Connection(databaseName string) {
	elemental.Connect(strings.Replace(mocks.DEFAULT_DATASOURCE, mocks.DEFAULT_DB_NAME, databaseName, 1))
}

func Seed(databaseName string) {
	fixtures.UserModel.SetDatabase(databaseName).InsertMany(mocks.Users).Exec()
}

func SeededConnection(databaseName string) {
	Connection(databaseName)
	Seed(databaseName)
}

func Teardown() {
	elemental.DropAllDatabases()
	elemental.Disconnect()
}
