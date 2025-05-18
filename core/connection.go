package elemental

import (
	"context"
	"maps"
	"slices"
	"sync"
	"time"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

// The default connection timeout for elemental connections
const ConnectionTimeout = 30 * time.Second

var clients = make(map[string]*mongo.Client)
var defaultDatabases = make(map[string]string)
var mu sync.RWMutex

// Elemental connection options
type ConnectionOptions struct {
	Alias         string                 // The alias of the connection, if not provided, it will be set to "default"
	URI           string                 // The connection string to connect to the database
	ClientOptions *options.ClientOptions // The options to use when creating the client
	PoolMonitor   *event.PoolMonitor     // The underlying event pool monitor to use when creating the client
}

// Connect to a new data source.
//
// @param arg - The connection string or ConnectionOptions struct
//
// Example:
//
//		client1 := Connect("mongodb://localhost:27017")
//	    // or
//		client2 := Connect(ConnectionOptions{
//			Alias: "secondary",
//			URI: "mongodb://localhost:27018",
//			ClientOptions: options.Client().SetMaxPoolSize(10),
//		})
func Connect(arg any) mongo.Client {
	mu.Lock()
	defer mu.Unlock()

	opts := ConnectionOptions{}

	if _, ok := arg.(string); ok {
		opts.URI = arg.(string)
	} else if _, ok := arg.(ConnectionOptions); ok {
		opts = arg.(ConnectionOptions)
	} else {
		panic(ErrInvalidConnectionArgument)
	}

	opts.Alias = lo.CoalesceOrEmpty(opts.Alias, "default")
	clientOpts := lo.CoalesceOrEmpty(opts.ClientOptions, options.Client()).
		SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1)).
		SetPoolMonitor(lo.CoalesceOrEmpty(opts.PoolMonitor, defaultPoolMonitor(opts.Alias)))
	if clientOpts.GetURI() == "" {
		if opts.URI == "" {
			panic(ErrURIRequired)
		}
		clientOpts = clientOpts.ApplyURI(opts.URI)
	}
	cs, err := connstring.ParseAndValidate(clientOpts.GetURI())
	if err != nil {
		panic(err)
	}
	defaultDatabases[opts.Alias] = cs.Database
	ctx, cancel := context.WithTimeout(context.Background(), *lo.CoalesceOrEmpty(clientOpts.ConnectTimeout, lo.ToPtr(ConnectionTimeout)))
	defer cancel()
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		panic(err)
	}
	lo.Must0(client.Ping(ctx, readpref.Primary()))
	clients[opts.Alias] = client
	triggerEventIfRegistered(opts.Alias, EventDeploymentDiscovered)
	return *client
}

// Get the database connection for a given alias or the default connection if no alias is provided
//
// @param alias - The alias of the connection to get
func GetConnection(alias ...string) *mongo.Client {
	client, ok := clients[lo.CoalesceOrEmpty(lo.FirstOrEmpty(alias), "default")]
	if !ok {
		return &mongo.Client{}
	}
	return client
}

// Same as 'GetConnection' method
var GetClient = GetConnection

// Disconnect a set of connections by alias or disconnect all connections if no alias is provided
//
// @param aliases - The aliases of the connections to disconnect
func Disconnect(aliases ...string) error {
	mu.Lock()
	defer mu.Unlock()

	if len(aliases) == 0 {
		aliases = slices.AppendSeq(aliases, maps.Keys(clients))
	}
	for _, alias := range aliases {
		err := clients[alias].Disconnect(context.Background())
		if err != nil {
			return err
		}
		delete(clients, alias)
		delete(defaultDatabases, alias)
	}
	return nil
}

// UseDatabase a specific database on a connection
//
// @param database - The name of the database to use
//
// @param alias - The alias of the connection to use
func UseDatabase(database string, alias ...string) *mongo.Database {
	return GetConnection(alias...).
		Database(lo.CoalesceOrEmpty(database, defaultDatabases[lo.CoalesceOrEmpty(lo.FirstOrEmpty(alias), "default")]))
}

// Use the default database on a connection. Uses the default connection if no alias is provided
//
// @param alias - The alias of the connection to use
func UseDefaultDatabase(alias ...string) *mongo.Database {
	return GetConnection(alias...).
		Database(lo.CoalesceOrEmpty(defaultDatabases[lo.CoalesceOrEmpty(lo.FirstOrEmpty(alias), "default")], "test"))
}

// Drops all databases across a given client or all clients if no alias is provided
func DropAllDatabases(alias ...string) {
	for key, client := range clients {
		if len(alias) > 0 && !lo.Contains(alias, key) {
			continue
		}
		databases, err := client.ListDatabaseNames(context.Background(), bson.D{{}}, options.ListDatabases().SetNameOnly(true))
		if err != nil {
			panic(err)
		}
		for _, db := range databases {
			client.Database(db).Drop(context.Background())
		}
	}
}

// Pings the primary server of a given connection or the default connection if no alias is provided
// This will timeout after 5 seconds
func Ping(alias ...string) error {
	client := GetConnection(alias...)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	return nil
}
