package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect(ctx context.Context, opts *options.ClientOptions, dbName string) (db *Database, err error) {
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return
	}

	db = &Database{client.Database(dbName)}
	return
}
