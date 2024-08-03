package postgres_repository

import (
	"context"
	"github.com/jackc/pgx/v4"
)

func TransactionHandler(ctx context.Context, tx pgx.Tx, err *error) {
	if p := recover(); p != nil {
		tx.Rollback(ctx)
		panic(p)
	} else if *err != nil {
		tx.Rollback(ctx)
	} else {
		*err = tx.Commit(ctx)
	}
}
