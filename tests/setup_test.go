package e_tests

import (
	"elemental/connection"
	"elemental/constants"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	uri = "mongodb://localhost:27017"
)

func TestSpec(t *testing.T) {
	Convey("Connect to a local database", t, func() {
		Convey("Simplest form of connect with just a URI", func() {
			client := e_connection.Connect(e_connection.ConnectionOptions{
				URI: uri,
			})
			So(client, ShouldNotBeNil)
		})
		Convey("Connect with a URI specified within client options", func() {
			opts := options.Client().ApplyURI(uri)
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
