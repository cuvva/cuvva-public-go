package mongodb

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	*mongo.Database
}

func (d Database) Collection(name string, opts ...*options.CollectionOptions) *Collection {
	return &Collection{d.Database.Collection(name, opts...)}
}

func (d Database) SetupSchemas(ctx context.Context, fs fs.FS, collectionNames []string) error {
	for _, colName := range collectionNames {
		file, err := fs.Open(fmt.Sprintf("%s.json", colName))
		if err != nil {
			return err
		}

		data, err := ioutil.ReadAll(file)
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
