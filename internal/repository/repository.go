package repository

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"github.com/sreekar2307/reconciler/internal/errors/db"
	"github.com/sreekar2307/reconciler/internal/model"
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
	client              *mongo.Client
	reconDatabase       *mongo.Database
	incomingCollections *mongo.Collection
	outgoingCollections *mongo.Collection
}

const (
	incomingTransactionsCollection = "incoming_transactions"
	outgoingTransactionsCollection = "outgoing_transactions"
)

func NewRepository(client *mongo.Client, reconDatabase string) Repository {
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
	cursor, err := r.incomingCollections.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find unreconciled incoming transactions: %w", err)
	}
	defer func() {
		_ = cursor.Close(ctx)
	}()
	var transactions []*model.Transaction
	if err := cursor.All(ctx, &transactions); err != nil {
		return nil, fmt.Errorf("failed to decode unreconciled incoming transactions: %w", err)
	}
	return transactions, nil
}

func (r repository) FindUnReconciledOutgoingTransactions(ctx context.Context) ([]*model.Transaction, error) {
	filter := bson.M{
		"reconciled": false,
		"timestamp":  bson.M{"$lt": time.Now().Add(-10 * time.Minute)},
	}
	cursor, err := r.outgoingCollections.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find unreconciled outgoing transactions: %w", err)
	}
	defer func() {
		_ = cursor.Close(ctx)
	}()
	var transactions []*model.Transaction
	if err := cursor.All(ctx, &transactions); err != nil {
		return nil, fmt.Errorf("failed to decode unreconciled outgoing transactions: %w", err)
	}
	return transactions, nil
}

func (r repository) FindOutgoingTransactionByID(ctx context.Context, txnID string) (*model.Transaction, error) {
	outTxn := new(model.Transaction)
	if err := r.outgoingCollections.FindOne(ctx, bson.M{"txn_id": txnID}).Decode(&outTxn); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, db.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to find outgoing transaction by ID %s: %w", txnID, err)
	}
	return outTxn, nil
}

func (r repository) FindIncomingTransactionByID(ctx context.Context, txnID string) (*model.Transaction, error) {
	outTxn := new(model.Transaction)
	if err := r.incomingCollections.FindOne(ctx, bson.M{"txn_id": txnID}).Decode(&outTxn); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, db.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to find incoming transaction by ID %s: %w", txnID, err)
	}
	return outTxn, nil
}

func (r repository) SetReconciled(ctx context.Context, txnID string) error {
	session, err := r.client.StartSession()
	if err != nil {
		return fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)
	_, err = session.WithTransaction(ctx, func(ctx mongo.SessionContext) (any, error) {
		if _, err := r.incomingCollections.UpdateOne(ctx, bson.M{"txn_id": txnID}, bson.M{"$set": bson.M{"reconciled": true}}); err != nil {
			return nil, fmt.Errorf("failed to update incoming transaction %s as reconciled: %w", txnID, err)
		}
		if _, err := r.outgoingCollections.UpdateOne(ctx, bson.M{"txn_id": txnID}, bson.M{"$set": bson.M{"reconciled": true}}); err != nil {
			return nil, fmt.Errorf("failed to update outgoing transaction %s as reconciled: %w", txnID, err)
		}
		return nil, nil
	}, &options.TransactionOptions{
		ReadConcern:  readconcern.Majority(),
		WriteConcern: writeconcern.Majority(),
	})
	if err != nil {
		return fmt.Errorf("failed to set reconciled for transaction %s: %w", txnID, err)
	}
	return nil
}
