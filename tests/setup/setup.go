package e_test_setup

import (
	"context"

	e_connection "github.com/elcengine/elemental/connection"
	e_test_base "github.com/elcengine/elemental/tests/base"
	e_mocks "github.com/elcengine/elemental/tests/mocks"
)

func Connection() {
	e_connection.ConnectURI(e_mocks.DEFAULT_DB_URI)
}

func Seed() {
	e_test_base.UserModel.InsertMany(e_mocks.Users).Exec()
}

func SeededConnection() {
	Connection()
	Seed()
}

func Teardown() {
	e_connection.UseDefault().Drop(context.TODO())
	e_connection.Use(e_mocks.SECONDARY_DB).Drop(context.TODO())
	e_connection.Use(e_mocks.TERTIARY_DB).Drop(context.TODO())
	e_connection.Use(e_mocks.TEMPORARY_DB_1).Drop(context.TODO())
	e_connection.Use(e_mocks.TEMPORARY_DB_2).Drop(context.TODO())
	e_connection.Use(e_mocks.TEMPORARY_DB_3).Drop(context.TODO())
	e_connection.Disconnect()
}
