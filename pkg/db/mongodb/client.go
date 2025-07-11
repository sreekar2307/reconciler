package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client interface {
	Database(string) Database
	StartSession(context.Context) (Session, error)
}

type client struct {
	mongoClient *mongo.Client
}

func (c *client) Database(name string) Database {
	return &database{
		mongoDatabase: c.mongoClient.Database(name),
	}
}

func (c *client) StartSession(_ context.Context) (Session, error) {
	s, err := c.mongoClient.StartSession()
	if err != nil {
		return nil, fmt.Errorf("starting session: %w", err)
	}
	return &session{s}, nil
}

func NewClient(ctx context.Context, uri string) (Client, error) {
	c, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("connecting to MongoDB URI %s: %w", uri, err)
	}
	return &client{c}, nil
}
