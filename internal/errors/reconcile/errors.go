package reconcile

import "errors"

var (
	ErrAmountMismatch              = errors.New("amount mismatch")
	ErrOutgoingTransactionNotFound = errors.New("outgoing transaction not found")
)
