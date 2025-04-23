package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Connect(opts *options.ClientOptions, dbName string) (db *Database, err error) {
	client, err := mongo.Connect(opts)
	if err != nil {
		return
	}

	if opts.ConnectTimeout == nil || *opts.ConnectTimeout <= 0 {
		return nil, fmt.Errorf("invalid connect timeout: %v", opts.ConnectTimeout)
	}

	ctx, cancel := context.WithTimeout(context.Background(), *opts.ConnectTimeout)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		return
	}

	db = &Database{client.Database(dbName)}
	return
}
