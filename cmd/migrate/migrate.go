package migrate

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/sreekar2307/reconciler/cmd"
)

func Run(ctx context.Context, deps *cmd.Deps) error {
	client := deps.Client
	reconDatabase := client.Database(deps.ReconDatabase)
	if err := createCollectionWithValidation(ctx, reconDatabase, "incoming_transactions"); err != nil {
		return fmt.Errorf("failed to create incoming_transactions collection: %w", err)
	}
	if err := createCollectionWithValidation(ctx, reconDatabase, "outgoing_transactions"); err != nil {
		return fmt.Errorf("failed to create outgoing_transactions collection: %w", err)
	}
	if err := createTxnIDIndex(ctx, reconDatabase, "incoming_transactions"); err != nil {
		return fmt.Errorf("failed to create index on incoming_transactions: %w", err)
	}
	if err := createTxnIDIndex(ctx, reconDatabase, "outgoing_transactions"); err != nil {
		return fmt.Errorf("failed to create index on outgoing_transactions: %w", err)
	}
	return nil

}

func createCollectionWithValidation(ctx context.Context, db *mongo.Database, name string) error {
	command := bson.D{
		{"create", name},
		{"validator", bson.M{
			"$jsonSchema": bson.M{
				"bsonType": "object",
				"required": []string{"txn_id", "amount", "reconciled"},
				"properties": bson.M{
					"txn_id":     bson.M{"bsonType": "string"},
					"amount":     bson.M{"bsonType": "double"},
					"reconciled": bson.M{"bsonType": "bool"},
				},
			},
		}},
	}
	return db.RunCommand(ctx, command).Err()
}

func createTxnIDIndex(ctx context.Context, db *mongo.Database, collName string) error {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "txn_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := db.Collection(collName).Indexes().CreateOne(ctx, indexModel)
	return err
}
