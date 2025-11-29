package repository

import (
	"context"
)

// TransactionRepositoryはトランザクション操作を管理します。
type TransactionRepository interface {
	Run(ctx context.Context, f func(txCtx context.Context) error) error
}
