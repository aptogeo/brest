package brest

import (
	"context"

	"github.com/uptrace/bun"
)

// The contextKey type is unexported to prevent collisions with context keys defined in
// other packages
type contextKey string

// ValueFromContext retrives Value from context
func ValueFromContext(ctx context.Context, keyName string) interface{} {
	return ctx.Value(contextKey(keyName))
}

// ContextWithValue sets Value to context request
func ContextWithValue(ctx context.Context, keyName string, value interface{}) context.Context {
	return context.WithValue(ctx, contextKey(keyName), value)
}

// DbFromContext retrives Db from context
func DbFromContext(ctx context.Context) *bun.DB {
	v := ValueFromContext(ctx, "db")
	if v == nil {
		return nil
	}
	return v.(*bun.DB)
}

// ContextWithDb sets Db to context request
func ContextWithDb(ctx context.Context, db *bun.DB) context.Context {
	return context.WithValue(ctx, contextKey("db"), db)
}

// TxFromContext retrives Tx from context
func TxFromContext(ctx context.Context) *bun.Tx {
	v := ValueFromContext(ctx, "tx")
	if v == nil {
		return nil
	}
	return v.(*bun.Tx)
}

// ContextWithTx sets Tx to context request
func ContextWithTx(ctx context.Context, tx *bun.Tx) context.Context {
	return context.WithValue(ctx, contextKey("tx"), tx)
}
