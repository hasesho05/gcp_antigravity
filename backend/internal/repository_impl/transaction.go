package repository_impl

import (
	"context"

	"cloud.google.com/go/firestore"
	"gcp_antigravity/backend/internal/repository"
)

// transactionKeyはcontext.Valueのキーとして使用されます。
type transactionKey struct{}

// GetTransactionはコンテキストからFirestoreトランザクションを抽出します。
func GetTransaction(ctx context.Context) (*firestore.Transaction, bool) {
	tx, ok := ctx.Value(transactionKey{}).(*firestore.Transaction)
	return tx, ok
}

type transactionRepository struct {
	client *firestore.Client
}

// NewTransactionRepositoryは新しいTransactionRepositoryを作成します。
func NewTransactionRepository(client *firestore.Client) repository.TransactionRepository {
	return &transactionRepository{client: client}
}

// RunはFirestoreトランザクション内で関数を実行します。
func (r *transactionRepository) Run(ctx context.Context, f func(txCtx context.Context) error) error {
	return r.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		ctxWithTx := context.WithValue(ctx, transactionKey{}, tx)
		return f(ctxWithTx)
	})
}
