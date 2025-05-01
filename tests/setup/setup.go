package e_test_setup

import (
	"fmt"

	e_connection "github.com/elcengine/elemental/connection"
	e_test_base "github.com/elcengine/elemental/tests/base"
	e_mocks "github.com/elcengine/elemental/tests/mocks"
)

func Connection(databaseName string) {
	e_connection.ConnectURI(fmt.Sprintf("%s/%s", e_mocks.DEFAULT_DATASOURCE, databaseName))
}

func Seed() {
	e_test_base.UserModel.InsertMany(e_mocks.Users).Exec()
}

func SeededConnection(databaseName string) {
	Connection(databaseName)
	Seed()
}

func Teardown() {
	e_connection.DropAll()
	e_connection.Disconnect()
}
