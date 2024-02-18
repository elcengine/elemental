package e_connection

import (
	"context"
	"elemental/constants"
	"elemental/utils"
	"time"

	"github.com/clubpay/qlubkit-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

var client *mongo.Client
var defaultDatabase string

const connectionTimeout = 30 * time.Second

type ConnectionOptions struct {
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
	cs := qkit.Ok(connstring.ParseAndValidate(clientOpts.GetURI()))
	defaultDatabase = cs.Database
	ctx, cancel := context.WithTimeout(context.Background(), *qkit.Coalesce(clientOpts.ConnectTimeout, qkit.ValPtr(connectionTimeout)))
	defer cancel()
	client := qkit.Must(mongo.Connect(ctx, clientOpts))
	e_utils.Must(client.Ping(ctx, readpref.Primary()))
	return *client
}

func Disconnect() error {
	err := client.Disconnect(context.Background())
	return err
}

func Use(databaseName string) *mongo.Database {
	return client.Database(databaseName)
}

func UseDefault() *mongo.Database {
	return client.Database(qkit.Coalesce(defaultDatabase, "test"))
}
