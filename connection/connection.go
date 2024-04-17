package e_connection

import (
	"context"
	"elemental/constants"
	"elemental/utils"
	"golang.org/x/exp/maps"
	"time"

	"github.com/clubpay/qlubkit-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

const connectionTimeout = 30 * time.Second

var clients = make(map[string]mongo.Client)
var defaultDatabases = make(map[string]string)

type ConnectionOptions struct {
	Alias         string
	URI           string
	ClientOptions *options.ClientOptions
}

func Connect(opts ConnectionOptions) mongo.Client {
	clientOpts := qkit.Coalesce(opts.ClientOptions, options.Client()).SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1))
	if clientOpts.GetURI() == "" {
		clientOpts = clientOpts.ApplyURI(opts.URI)
		if clientOpts.GetURI() == "" {
			panic(e_constants.ErrURIRequired)
		}
	}
	cs := qkit.Must(connstring.ParseAndValidate(clientOpts.GetURI()))
	opts.Alias = qkit.Coalesce(opts.Alias, "default")
	defaultDatabases[opts.Alias] = cs.Database
	ctx, cancel := context.WithTimeout(context.Background(), *qkit.Coalesce(clientOpts.ConnectTimeout, qkit.ValPtr(connectionTimeout)))
	defer cancel()
	client := qkit.Must(mongo.Connect(ctx, clientOpts))
	e_utils.Must(client.Ping(ctx, readpref.Primary()))
	return *client
}

// Simplest form of connect with just a URI and no options
func ConnectURI(uri string) mongo.Client {
	return Connect(ConnectionOptions{URI: uri})
}

// Get the database connection for a given alias or the default connection if no alias is provided
//
// @param alias - The alias of the connection to get
func GetConnection(alias ...string) mongo.Client {
	return clients[qkit.Coalesce(alias[0], "default")]
}

// Disconnect a set of connections by alias or disconnect all connections if no alias is provided
//
// @param aliases - The aliases of the connections to disconnect
func Disconnect(aliases ...string) error {
	if len(aliases) == 0 {
		aliases = maps.Keys(clients)
	}
	for _, alias := range aliases {
		err := qkit.ValPtr(clients[alias]).Disconnect(context.Background())
		if err != nil {
			return err
		}
		delete(clients, alias)
		delete(defaultDatabases, alias)
	}
	return nil
}

// Use a specific database on a connection
//
// @param database - The name of the database to use
//
// @param alias - The alias of the connection to use
func Use(database string, alias ...string) *mongo.Database {
	return qkit.ValPtr(clients[qkit.Coalesce(alias[0], "default")]).Database(database)
}

// Use the default database on a connection. Uses the default connection if no alias is provided
//
// @param alias - The alias of the connection to use
func UseDefault(alias ...string) *mongo.Database {
	return qkit.ValPtr(clients[qkit.Coalesce(alias[0], "default")]).Database(qkit.Coalesce(defaultDatabases[qkit.Coalesce(alias[0], "default")], "test"))
}
