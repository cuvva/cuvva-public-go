package mongodb

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Database struct {
	*mongo.Database
}

func (d Database) Collection(name string, opts options.Lister[options.CollectionOptions]) *Collection {
	return &Collection{d.Database.Collection(name, opts)}
}

func (d Database) SetupSchemas(ctx context.Context, fs fs.FS, collectionNames []string) error {
	for _, colName := range collectionNames {
		file, err := fs.Open(fmt.Sprintf("%s.json", colName))
		if err != nil {
			return err
		}

		data, err := io.ReadAll(file)
		if err != nil {
			return err
		}

		var schema interface{}
		if err := json.Unmarshal(data, &schema); err != nil {
			return err
		}

		if err := d.RunCommand(ctx, bson.D{
			{"collMod", colName},
			{"validationLevel", "strict"},
			{"validationAction", "error"},
			{"validator", bson.M{
				"$jsonSchema": schema,
			}},
		}).Err(); err != nil {
			return err
		}
	}

	return nil
}

func (db *Database) DoTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return db.DoTxWithOptions(ctx, options.Session(), fn)
}

func (db *Database) DoTxWithOptions(ctx context.Context, opts *options.SessionOptionsBuilder, fn func(ctx context.Context) error) error {
	return db.Client().UseSessionWithOptions(ctx, opts, func(ctx context.Context) error {
		return mongo.WithSession(ctx, mongo.SessionFromContext(ctx), func(ctx context.Context) error {
			return fn(ctx)
		})
	})
}
