package e_test_setup

import (
	"context"
	"elemental/connection"
	"elemental/tests/base"
	"elemental/tests/mocks"
)

func Connection() {
	e_connection.ConnectURI(e_mocks.URI)
	e_connection.UseDefault().Drop(context.TODO())
}

func Seed() {
	e_test_base.UserModel.InsertMany(e_mocks.Users)
}

func SeededConnection() {
	Connection()
	Seed()
}