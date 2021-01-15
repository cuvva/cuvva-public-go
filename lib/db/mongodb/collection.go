package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Collection struct {
	*mongo.Collection
}

func (c Collection) SetupIndexes(models []mongo.IndexModel) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = c.Indexes().CreateMany(ctx, models)
	return
}

func (c Collection) FindAll(ctx context.Context, filter, results interface{}, opts ...*options.FindOptions) (err error) {
	cur, err := c.Find(ctx, filter, opts...)
	if err == nil {
		err = cur.All(ctx, results)
	}
	return
}

func (c Collection) DistinctString(ctx context.Context, fieldName string, filter interface{}, opts ...*options.DistinctOptions) ([]string, error) {
	pipeline := []bson.M{
		{"$match": filter},
		{"$group": bson.M{"_id": nil, "set": bson.M{"$addToSet": "$" + fieldName}}},
	}

	dOpts := options.MergeDistinctOptions(opts...)
	aOpts := &options.AggregateOptions{
		Collation: dOpts.Collation,
		MaxTime:   dOpts.MaxTime,
	}

	cur, err := c.Aggregate(ctx, pipeline, aOpts)
	if err != nil {
		return nil, err
	}
	if !cur.Next(ctx) {
		return []string{}, nil
	}

	var out struct {
		Set []string `bson:"set"`
	}
	if err = cur.Decode(&out); err != nil {
		return nil, err
	}
	return out.Set, nil
}
