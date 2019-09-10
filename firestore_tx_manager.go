package tx_manager

import (
	"cloud.google.com/go/firestore"
	"context"
)

type ctxKey string

var txKey = ctxKey("tx_manager")

// TxManager can manage firestore transaction
type TxManager interface {
	RunTx(context.Context, func(context.Context) error) error
}

// txManager represents firestore transaction manager
type txManager struct {
	store *firestore.Client
}

// NewTxManager initializes txManager
func NewTxManager(store *firestore.Client) *txManager {
	return &txManager{store: store}
}

// RunTx runs tx_manager
func (t *txManager) RunTx(ctx context.Context, f func(ctx context.Context) error) error {
	return t.store.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		ctx = context.WithValue(ctx, txKey, tx)
		return f(ctx)
	})
}

// GetTx extracts tx_manager from context
func GetTx(ctx context.Context) (*firestore.Transaction, bool) {
	tx, ok := ctx.Value(txKey).(*firestore.Transaction)
	return tx, ok
}

