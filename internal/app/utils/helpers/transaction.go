// Package helpers provides helper functions.
package helpers

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// CommitTx is a helper function to commit or rollback phx Transactions.
func CommitTx(ctx context.Context, tx pgx.Tx, err error) {
	if err != nil {
		tx.Rollback(ctx)
	} else {
		tx.Commit(ctx)
	}
}
