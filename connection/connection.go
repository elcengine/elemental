// Deprecated: Use package "github.com/elcengine/elemental/core" instead
//
// This package will be completely removed once v2.0.0 is released
package e_connection

import (
	elemental "github.com/elcengine/elemental/core"
	"go.mongodb.org/mongo-driver/mongo"
)

// Elemental connection options
//
// Deprecated: Use 'elemental.ConnectionOptions' instead
type ConnectionOptions = elemental.ConnectionOptions

// Connect to a new data source with custom options
//
// Deprecated: Use 'elemental.Connect' instead
func Connect(opts ConnectionOptions) mongo.Client {
	return elemental.Connect(opts)
}

// Simplest form of connect with just a URI and no options
//
// Deprecated: Use 'elemental.Connect' instead
func ConnectURI(uri string) mongo.Client {
	return elemental.Connect(uri)
}

// Get the database connection for a given alias or the default connection if no alias is provided
//
// Deprecated: Use 'elemental.GetConnection' instead
var GetConnection = elemental.GetConnection

// Disconnect a set of connections by alias or disconnect all connections if no alias is provided
//
// Deprecated: Use 'elemental.Disconnect' instead
var Disconnect = elemental.Disconnect

// Use a specific database on a connection
//
// Deprecated: Use 'elemental.UseDatabase' instead
var Use = elemental.UseDatabase

// Use the default database on a connection. Uses the default connection if no alias is provided
//
// Deprecated: Use 'elemental.UseDefaultDatabase' instead
var UseDefault = elemental.UseDefaultDatabase

// Drops all databases across a given client or all clients if no alias is provided
//
// Deprecated: Use 'elemental.DropAllDatabases' instead
var DropAll = elemental.DropAllDatabases
