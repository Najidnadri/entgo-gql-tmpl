package infrastructure

import (
	"chipin/ent"
	"context"
	"fmt"
)

func NewTxCtx(ctx context.Context, client *ent.Client) (context.Context, *ent.Tx, error) {
	tx, err := client.Tx(ctx)
	if err != nil {
		return ctx, nil, fmt.Errorf("error starting new transaction client: %v", err)
	}
	ctx2 := ent.NewContext(ctx, tx.Client())
	return ctx2, tx, nil
}

func TxRollback(ctx context.Context, tx *ent.Tx) error {
	err := tx.Rollback()
	if err != nil {
		fmt.Println("error rolling back transaction: ", err)
		return fmt.Errorf("error rolling back transaction: %v", err)
	}
	return nil
}

func TxCommit(ctx context.Context, tx *ent.Tx) error {
	err := tx.Commit()
	if err != nil {
		fmt.Println("error committing transaction: ", err)
		return fmt.Errorf("error committing transaction: %v", err)
	}
	return nil
}

func TransactionWrapper(ctx context.Context, client *ent.Client, fn func(ctx context.Context) error) error {
	ctx2, tx, err := NewTxCtx(ctx, client)
	if err != nil {
		return fmt.Errorf("error creating new transaction context: %v", err)
	}
	err = fn(ctx2)
	if err != nil {
		TxRollback(ctx, tx)
		return err
	}
	return TxCommit(ctx, tx)
}
