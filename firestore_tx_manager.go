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

// CreateWithTx inserts data into document with transaction if exist. If not, it will create as usual.
func CreateWithTx(ctx context.Context, dr *firestore.DocumentRef, data interface{}) error {
	tx, ok := GetTx(ctx)
	if ok {
		return tx.Create(dr, data)
	} else {
		_, err := dr.Create(ctx, data)
		return err
	}
}

// SetWithTx set data to document with transaction if exist. If not, it will set as usual.
func SetWithTx(ctx context.Context, dr *firestore.DocumentRef, data interface{}, opts ...firestore.SetOption) error {
	tx, ok := GetTx(ctx)
	if ok {
		return tx.Set(dr, data, opts...)
	} else {
		_, err := dr.Set(ctx, data, opts...)
		return err
	}
}


// UpdateWithTx updates document with transaction if exist. If not, it will update as usual.
func UpdateWithTx(ctx context.Context, dr *firestore.DocumentRef, data []firestore.Update, preconds ...firestore.Precondition) error {
	tx, ok := GetTx(ctx)
	if ok {
		return tx.Update(dr, data, preconds...)
	} else {
		_, err := dr.Update(ctx, data, preconds...)
		return err
	}
}


// DeleteWithTx deletes document with transaction if exist. If not, it will delete as usual.
func DeleteWithTx(ctx context.Context, dr *firestore.DocumentRef, preconds ...firestore.Precondition) error {
	tx, ok := GetTx(ctx)
	if ok {
		return tx.Delete(dr, preconds...)
	} else {
		_, err := dr.Delete(ctx, preconds...)
		return err
	}
}
