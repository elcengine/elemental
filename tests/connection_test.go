package e_tests

import (
	"fmt"
	"strings"
	"testing"

	e_constants "github.com/elcengine/elemental/constants"
	elemental "github.com/elcengine/elemental/core"
	e_mocks "github.com/elcengine/elemental/tests/mocks"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestConnection(t *testing.T) {
	Convey("Connect to a local database", t, func() {
		Convey("Simplest form of connect with just a URI", func() {
			connectionCreatedFired := false

			elemental.OnConnectionEvent(event.ConnectionCreated, func() {
				connectionCreatedFired = true
			})

			defer elemental.RemoveConnectionEvent(event.ConnectionCreated)

			client := elemental.Connect(strings.Replace(e_mocks.DEFAULT_DATASOURCE, e_mocks.DEFAULT_DB_NAME, t.Name(), 1))
			So(client, ShouldNotBeNil)

			Convey("Should use the default database", func() {
				So(elemental.UseDefaultDatabase().Name(), ShouldEqual, t.Name())
			})
			Convey("Should use the specified database", func() {
				DATABASE := fmt.Sprintf("%s_%s", t.Name(), "secondary")
				So(elemental.UseDatabase(DATABASE).Name(), ShouldEqual, DATABASE)
			})
			Convey("Should fire the connection created event", func() {
				So(connectionCreatedFired, ShouldBeTrue)
			})
			Convey("Should ping the primary server", func() {
				So(elemental.Ping(), ShouldBeNil)
			})
		})
		Convey("Connect with a URI specified within client options", func() {
			opts := options.Client().ApplyURI(strings.Replace(e_mocks.DEFAULT_DATASOURCE, e_mocks.DEFAULT_DB_NAME, t.Name(), 1))
			client := elemental.Connect(elemental.ConnectionOptions{
				ClientOptions: opts,
			})
			So(client, ShouldNotBeNil)
		})
		Convey("Connect with no URI", func() {
			So(func() {
				elemental.Connect(elemental.ConnectionOptions{})
			}, ShouldPanicWith, e_constants.ErrURIRequired)
		})
		Convey("Connect with invalid argument", func() {
			So(func() {
				elemental.Connect(123)
			}, ShouldPanicWith, e_constants.ErrInvalidConnectionArgument)
		})
	})
}
