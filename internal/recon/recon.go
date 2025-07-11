package recon

import (
	"context"
	"errors"
	"fmt"
	"log"
	dbErrors "github.com/sreekar2307/reconciler/internal/errors/db"
	"github.com/sreekar2307/reconciler/internal/repository"
)

type Reconcile interface {
	Reconcile(ctx context.Context) error
}

type reconcile struct {
	repository repository.Repository
}

func NewReconcile(repo repository.Repository) Reconcile {
	return &reconcile{
		repository: repo,
	}
}

func (r reconcile) Reconcile(ctx context.Context) error {
	incomingTransactions, err := r.repository.FindUnReconciledIncomingTransactions(ctx)
	if err != nil {
		return fmt.Errorf("failed to find unreconciled incoming transactions: %w", err)
	}
	for _, inTxn := range incomingTransactions {
		outTxn, err := r.repository.FindOutgoingTransactionByID(ctx, inTxn.TxnID)
		if err != nil {
			if errors.Is(err, dbErrors.ErrRecordNotFound) {
				log.Println("reconcile: outgoing transaction not found for incoming transaction ID:", inTxn.TxnID)
				continue
			}
			return fmt.Errorf("failed to find outgoing transaction by ID %s: %w", inTxn.TxnID, err)
		}
		if inTxn.Amount != outTxn.Amount || inTxn.Currency != outTxn.Currency {
			log.Println("reconcile: amount or currency mismatch for transaction ID:", inTxn.TxnID)
			continue
		}
		if err := r.repository.SetReconciled(ctx, inTxn.TxnID); err != nil {
			return fmt.Errorf("failed to set transaction %s as reconciled: %w", inTxn.TxnID, err)
		}
	}

	return nil
}
