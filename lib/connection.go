package elemental

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func Connect(opts options.BSONOptions) {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().SetBSONOptions(&opts))
	if err != nil {
		panic(err)
	}
}

func Disconnect() error {
	err := client.Disconnect(context.Background())
	return err
}

func Use(databaseName string) *mongo.Database {
	return client.Database(databaseName)
}

func UseDefault() *mongo.Database {
	return client.Database("mailman")
}
