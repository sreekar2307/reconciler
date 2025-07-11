package migrate

import (
	"context"
	"fmt"
	"github.com/sreekar2307/reconciler/cmd"
	"github.com/sreekar2307/reconciler/pkg/db/mongodb"
	"go.mongodb.org/mongo-driver/bson"
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

func createCollectionWithValidation(ctx context.Context, db mongodb.Database, name string) error {
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
	return db.RunCommand(ctx, command)
}

func createTxnIDIndex(ctx context.Context, db mongodb.Database, collName string) error {
	t := true
	indexModel := mongodb.Index{
		Keys: bson.D{{Key: "txn_id", Value: 1}},
		Opts: &mongodb.IndexOptions{Unique: &t},
	}
	return db.Collection(collName).CreateIndex(ctx, indexModel)
}
