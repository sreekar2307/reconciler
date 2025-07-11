package repository

import (
	"context"
	"fmt"
	"github.com/sreekar2307/reconciler/internal/model"
	"github.com/sreekar2307/reconciler/pkg/db/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type Repository interface {
	FindUnReconciledIncomingTransactions(ctx context.Context) ([]*model.Transaction, error)
	FindUnReconciledOutgoingTransactions(ctx context.Context) ([]*model.Transaction, error)
	FindOutgoingTransactionByID(ctx context.Context, txnID string) (*model.Transaction, error)
	FindIncomingTransactionByID(ctx context.Context, txnID string) (*model.Transaction, error)
	SetReconciled(ctx context.Context, txnID string) error
}

type repository struct {
	client              mongodb.Client
	reconDatabase       mongodb.Database
	incomingCollections mongodb.Collection
	outgoingCollections mongodb.Collection
}

const (
	incomingTransactionsCollection = "incoming_transactions"
	outgoingTransactionsCollection = "outgoing_transactions"
)

func NewRepository(client mongodb.Client, reconDatabase string) Repository {
	reconDb := client.Database(reconDatabase)
	return &repository{
		client:              client,
		reconDatabase:       reconDb,
		incomingCollections: reconDb.Collection(incomingTransactionsCollection),
		outgoingCollections: reconDb.Collection(outgoingTransactionsCollection),
	}
}

func (r repository) FindUnReconciledIncomingTransactions(ctx context.Context) ([]*model.Transaction, error) {
	filter := bson.M{
		"reconciled": false,
		"timestamp":  bson.M{"$lt": time.Now().Add(-10 * time.Minute)},
	}
	results, err := r.incomingCollections.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find unreconciled incoming transactions: %w", err)
	}
	txns := make([]*model.Transaction, len(results))
	for i, result := range results {
		txn := new(model.Transaction)
		txns[i] = txn
		if err := bson.Unmarshal(result, txns[i]); err != nil {
			return nil, fmt.Errorf("failed to decode incoming transaction: %w", err)
		}
	}
	return txns, nil
}

func (r repository) FindUnReconciledOutgoingTransactions(ctx context.Context) ([]*model.Transaction, error) {
	filter := bson.M{
		"reconciled": false,
		"timestamp":  bson.M{"$lt": time.Now().Add(-10 * time.Minute)},
	}
	results, err := r.outgoingCollections.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find unreconciled incoming transactions: %w", err)
	}
	txns := make([]*model.Transaction, len(results))
	for i, result := range results {
		txn := new(model.Transaction)
		txns[i] = txn
		if err := bson.Unmarshal(result, txns[i]); err != nil {
			return nil, fmt.Errorf("failed to decode incoming transaction: %w", err)
		}
	}
	return txns, nil
}

func (r repository) FindOutgoingTransactionByID(ctx context.Context, txnID string) (*model.Transaction, error) {
	outTxn := new(model.Transaction)
	result, err := r.outgoingCollections.FindOne(ctx, bson.M{"txn_id": txnID})
	if err != nil {
		return nil, fmt.Errorf("failed to find outgoing transaction by ID %s: %w", txnID, err)
	}
	if err := bson.Unmarshal(result, outTxn); err != nil {
		return nil, fmt.Errorf("failed to decode outgoing transaction: %w", err)
	}
	return outTxn, nil
}

func (r repository) FindIncomingTransactionByID(ctx context.Context, txnID string) (*model.Transaction, error) {
	inTxn := new(model.Transaction)
	result, err := r.incomingCollections.FindOne(ctx, bson.M{"txn_id": txnID})
	if err != nil {
		return nil, fmt.Errorf("failed to find incoming transaction by ID %s: %w", txnID, err)
	}
	if err := bson.Unmarshal(result, inTxn); err != nil {
		return nil, fmt.Errorf("failed to decode incoming transaction: %w", err)
	}
	return inTxn, nil
}

func (r repository) SetReconciled(ctx context.Context, txnID string) error {
	session, err := r.client.StartSession(ctx)
	if err != nil {
		return fmt.Errorf("failed to start session: %w", err)
	}
	defer session.End(ctx)
	_, err = session.WithTransaction(ctx, func(ctx context.Context) (any, error) {
		if err := r.incomingCollections.UpdateOne(ctx, bson.M{"txn_id": txnID}, bson.M{"$set": bson.M{"reconciled": true}}); err != nil {
			return nil, fmt.Errorf("failed to update incoming transaction %s as reconciled: %w", txnID, err)
		}
		if err := r.outgoingCollections.UpdateOne(ctx, bson.M{"txn_id": txnID}, bson.M{"$set": bson.M{"reconciled": true}}); err != nil {
			return nil, fmt.Errorf("failed to update outgoing transaction %s as reconciled: %w", txnID, err)
		}
		return nil, nil
	}, &mongodb.TransactionOptions{
		ReadConcern:  "majority",
		WriteConcern: "majority",
	})
	if err != nil {
		return fmt.Errorf("failed to set reconciled for transaction %s: %w", txnID, err)
	}
	return nil
}
