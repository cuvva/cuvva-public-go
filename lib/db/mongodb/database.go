package mongodb

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	*mongo.Database
}

func (d Database) Collection(name string, opts ...*options.CollectionOptions) *Collection {
	return &Collection{d.Database.Collection(name, opts...)}
}
