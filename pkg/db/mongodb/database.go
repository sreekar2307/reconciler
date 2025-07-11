package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Database interface {
	Collection(string) Collection
	RunCommand(context.Context, bson.D) error
}

type database struct {
	mongoDatabase *mongo.Database
}

func (d *database) Collection(name string) Collection {
	return &collection{
		mongoCollection: d.mongoDatabase.Collection(name),
	}
}

func (d *database) RunCommand(ctx context.Context, command bson.D) error {
	if err := d.mongoDatabase.RunCommand(ctx, command).Err(); err != nil {
		return err
	}
	return nil
}
