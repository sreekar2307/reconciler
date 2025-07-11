package seed

import (
	"context"
	"fmt"
	"github.com/sreekar2307/reconciler/cmd"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func Run(ctx context.Context, deps *cmd.Deps) error {
	reconDatabase := deps.Client.Database(deps.ReconDatabase)
	col := reconDatabase.Collection("incoming_transactions")
	dayFromNow := time.Now().UTC().Add(-24 * time.Hour)
	incomingTxns := bson.A{
		bson.M{"txn_id": "txn_001", "amount": 100.0, "reconciled": false, "timestamp": dayFromNow},
		bson.M{"txn_id": "txn_002", "amount": 250.0, "reconciled": false, "timestamp": dayFromNow},
		bson.M{"txn_id": "txn_003", "amount": 500.0, "reconciled": false, "timestamp": dayFromNow},
	}
	outgoingTxns := bson.A{
		bson.M{"txn_id": "txn_001", "amount": 100.0, "reconciled": false, "timestamp": dayFromNow},
		bson.M{"txn_id": "txn_003", "amount": 500.0, "reconciled": false, "timestamp": dayFromNow},
		bson.M{"txn_id": "txn_004", "amount": 300.0, "reconciled": false, "timestamp": dayFromNow},
	}
	err := col.InsertMany(ctx, incomingTxns)
	if err != nil {
		return fmt.Errorf("failed to seed incoming transactions: %w", err)
	}
	col = reconDatabase.Collection("outgoing_transactions")
	err = col.InsertMany(ctx, outgoingTxns)
	if err != nil {
		return fmt.Errorf("failed to seed outgoing transactions: %w", err)
	}
	return nil
}
