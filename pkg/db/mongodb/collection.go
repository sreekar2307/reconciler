package mongodb

import (
	"context"
	"errors"
	"fmt"
	dbErrors "github.com/sreekar2307/reconciler/internal/errors/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Collection interface {
	Find(context.Context, bson.M) ([]bson.Raw, error)
	FindOne(context.Context, bson.M) (bson.Raw, error)
	UpdateOne(context.Context, bson.M, bson.M) error
	CreateIndex(context.Context, Index) error
	InsertMany(context.Context, bson.A) error
}

type collection struct {
	mongoCollection *mongo.Collection
}

func (c *collection) Find(ctx context.Context, filter bson.M) ([]bson.Raw, error) {
	cursor, err := c.mongoCollection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}
	defer cursor.Close(ctx)

	results := make([]bson.Raw, 0)
	for cursor.Next(ctx) {
		results = append(results, cursor.Current)
	}
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode documents: %w", err)
	}
	return results, nil
}

func (c *collection) FindOne(ctx context.Context, filter bson.M) (bson.Raw, error) {
	result := c.mongoCollection.FindOne(ctx, filter)
	if err := result.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, dbErrors.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to find one document: %w", result.Err())
	}
	return result.Raw()
}

func (c *collection) UpdateOne(ctx context.Context, filter bson.M, update bson.M) error {
	_, err := c.mongoCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update one document: %w", err)
	}
	return nil
}

type Index struct {
	Keys bson.D
	Opts *IndexOptions
}

type IndexOptions struct {
	Unique *bool
}

func (c *collection) CreateIndex(ctx context.Context, index Index) error {
	_, err := c.mongoCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: index.Keys,
		Options: &options.IndexOptions{
			Unique: index.Opts.Unique,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to update one document: %w", err)
	}
	return nil
}

func (c *collection) InsertMany(ctx context.Context, documents bson.A) error {
	_, err := c.mongoCollection.InsertMany(ctx, documents)
	if err != nil {
		return fmt.Errorf("failed to insert many documents: %w", err)
	}
	return nil
}
