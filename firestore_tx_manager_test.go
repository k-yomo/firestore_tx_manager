package tx_manager

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"reflect"
	"testing"

	"cloud.google.com/go/firestore"
)

func TestNewTxManager(t *testing.T) {
	store, _ := firestore.NewClient(context.Background(), "test")
	type args struct {
		store *firestore.Client
	}
	tests := []struct {
		name string
		args args
		want *txManager
	}{
		{name: "Initialize transaction manager with given store", args: args{store: store}, want: &txManager{store: store}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTxManager(tt.args.store); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTxManager() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_txManager_RunTx(t *testing.T) {
	store, _ := firestore.NewClient(context.Background(), "test")
	tm := NewTxManager(store)
	type args struct {
		ctx context.Context
		f   func(ctx context.Context) error
	}
	tests := []struct {
		name    string
		t       *txManager
		args    args
		wantErr bool
	}{
		{
			name: "Rollback transaction when error occurred",
			t:    tm,
			args: args{ctx: context.Background(),
				f: func(ctx context.Context) error {
					tx, _ := GetTx(ctx)
					testCollection := tm.store.Collection("test")
					ref := testCollection.Doc("1")
					invalidRef := testCollection.Doc("")
					_ = tx.Set(ref, struct{ num int }{num: 1})
					err := tx.Set(invalidRef, struct{ num int }{num: 2})
					return err
				},
			},
			wantErr: true,
		},
		{
			name: "Commit transaction when no error occurred",
			t:    tm,
			args: args{ctx: context.Background(),
				f: func(ctx context.Context) error {
					tx, _ := GetTx(ctx)
					testCollection := tm.store.Collection("test")
					ref1 := testCollection.Doc("1")
					ref2 := testCollection.Doc("2")
					_ = tx.Set(ref1, struct{ num int }{num: 1})
					_ = tx.Set(ref2, struct{ num int }{num: 2})
					return nil
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.t.RunTx(tt.args.ctx, tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("txManager.RunTx() error = %v, wantErr %v", err, tt.wantErr)
			}
			testCollection := tm.store.Collection("test")
			ref, err := testCollection.Doc("1").Get(context.Background())
			if tt.wantErr {
				if status.Code(err) != codes.NotFound {
					t.Errorf("txManager.RunTx() got = %v, wantErr %v", ref.Data(), codes.NotFound)
				}
			} else {
				if status.Code(err) == codes.NotFound {
					t.Errorf("txManager.RunTx() got = %v, want %v", codes.NotFound, struct{ num int }{num: 1})
				}
			}
		})
	}
}

func TestGetTx(t *testing.T) {
	tx := &firestore.Transaction{}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name  string
		args  args
		want  *firestore.Transaction
		want1 bool
	}{
		{
			name:  "Get tranasction when tx is binded to context",
			args:  args{ctx: context.WithValue(context.Background(), txKey, tx)},
			want:  tx,
			want1: true,
		},
		{
			name:  "Get nil when tx is not binded to context",
			args:  args{ctx: context.Background()},
			want:  nil,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetTx(tt.args.ctx)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTx() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetTx() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

