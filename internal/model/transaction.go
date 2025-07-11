package model

import "time"

type Transaction struct {
	TxnID      string    `bson:"txn_id"`
	Amount     float64   `bson:"amount"`
	Currency   string    `bson:"currency"`
	Timestamp  time.Time `bson:"timestamp"`
	Source     string    `bson:"source"`
	Reconciled bool      `bson:"reconciled"`
}
