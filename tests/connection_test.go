package e_tests

import (
	"fmt"
	"strings"
	"testing"

	e_connection "github.com/elcengine/elemental/connection"
	e_constants "github.com/elcengine/elemental/constants"
	e_mocks "github.com/elcengine/elemental/tests/mocks"
	e_test_setup "github.com/elcengine/elemental/tests/setup"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestConnection(t *testing.T) {
	defer e_test_setup.Teardown()

	Convey("Connect to a local database", t, func() {
		Convey("Simplest form of connect with just a URI", func() {
			client := e_connection.Connect(e_connection.ConnectionOptions{
				URI: strings.Replace(e_mocks.DEFAULT_DATASOURCE, "elemental", t.Name(), 1),
			})
			So(client, ShouldNotBeNil)
			Convey("Should use the default database", func() {
				So(e_connection.UseDefault().Name(), ShouldEqual, t.Name())
			})
			Convey("Should use the specified database", func() {
				DATABASE := fmt.Sprintf("%s_%s", t.Name(), "secondary")
				So(e_connection.Use(DATABASE).Name(), ShouldEqual, DATABASE)
			})
		})
		Convey("Connect with a URI specified within client options", func() {
			opts := options.Client().ApplyURI(strings.Replace(e_mocks.DEFAULT_DATASOURCE, "elemental", t.Name(), 1))
			client := e_connection.Connect(e_connection.ConnectionOptions{
				ClientOptions: opts,
			})
			So(client, ShouldNotBeNil)
		})
		Convey("Connect with no URI", func() {
			So(func() {
				e_connection.Connect(e_connection.ConnectionOptions{})
			}, ShouldPanicWith, e_constants.ErrURIRequired)
		})
	})
}
