package brest_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aptogeo/brest"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
)

func TestTransactionalCurrentKO(t *testing.T) {
	db, _ := initTests(t)
	defer db.Close()
	ctx := brest.ContextWithDb(context.Background(), db)
	err := brest.Execute(ctx, func(ctx context.Context, tx *bun.Tx) error {
		_, err := tx.NewInsert().Model(&Todo{Text: "ko"}).Exec(ctx)
		assert.Nil(t, err)
		return errors.New("ko")
	})
	assert.NotNil(t, err)
	count, err := db.NewSelect().Model(&Todo{}).Count(ctx)
	assert.Nil(t, err)
	assert.Equal(t, 0, count)
}

func TestTransactionalCurrentOK(t *testing.T) {
	db, _ := initTests(t)
	ctx := brest.ContextWithDb(context.Background(), db)
	err := brest.Execute(ctx, func(ctx context.Context, tx *bun.Tx) error {
		_, err := tx.NewInsert().Model(&Todo{Text: "ok"}).Exec(ctx)
		return err
	})
	assert.Nil(t, err)
	count, err := db.NewSelect().Model(&Todo{}).Count(ctx)
	assert.Nil(t, err)
	assert.Equal(t, 1, count)
}

func TestTransactionalCurrentOKCurrentKO(t *testing.T) {
	db, _ := initTests(t)
	ctx := brest.ContextWithDb(context.Background(), db)
	err := brest.Execute(ctx, func(ctx context.Context, tx *bun.Tx) error {
		_, err := tx.NewInsert().Model(&Todo{Text: "ok"}).Exec(ctx)
		assert.Nil(t, err)
		return brest.Execute(ctx, func(ctx context.Context, tx *bun.Tx) error {
			_, err := tx.NewInsert().Model(&Todo{Text: "ko"}).Exec(ctx)
			assert.Nil(t, err)
			return errors.New("ko")
		})
	})
	assert.NotNil(t, err)
	count, err := db.NewSelect().Model(&Todo{}).Count(ctx)
	assert.Nil(t, err)
	assert.Equal(t, 0, count)
}

func TestTransactionalCurrentOKCurrentOK(t *testing.T) {
	db, _ := initTests(t)
	ctx := brest.ContextWithDb(context.Background(), db)
	err := brest.Execute(ctx, func(ctx context.Context, tx *bun.Tx) error {
		_, err := tx.NewInsert().Model(&Todo{Text: "ok"}).Exec(ctx)
		assert.Nil(t, err)
		return brest.Execute(ctx, func(ctx context.Context, tx *bun.Tx) error {
			_, err := tx.NewInsert().Model(&Todo{Text: "ok"}).Exec(ctx)
			return err
		})
	})
	assert.Nil(t, err)
	count, err := db.NewSelect().Model(&Todo{}).Count(ctx)
	assert.Nil(t, err)
	assert.Equal(t, 2, count)
}

func TestTransactionalMandatory(t *testing.T) {
	db, _ := initTests(t)
	var err error
	ctx := brest.ContextWithDb(context.Background(), db)
	err = brest.ExecuteWithPropagation(ctx, brest.Mandatory, func(ctx context.Context, tx *bun.Tx) error {
		_, err := tx.NewInsert().Model(&Todo{Text: "ok"}).Exec(ctx)
		return err
	})
	assert.NotNil(t, err)
	count, err := db.NewSelect().Model(&Todo{}).Count(ctx)
	assert.Nil(t, err)
	assert.Equal(t, 0, count)

	err = brest.Execute(ctx, func(ctx context.Context, tx *bun.Tx) error {
		return brest.ExecuteWithPropagation(ctx, brest.Current, func(ctx context.Context, tx *bun.Tx) error {
			_, err := tx.NewInsert().Model(&Todo{Text: "ok"}).Exec(ctx)
			return err
		})
	})
	assert.Nil(t, err)
	count, err = db.NewSelect().Model(&Todo{}).Count(ctx)
	assert.Nil(t, err)
	assert.Equal(t, 1, count)
}

func TestTransactionalSavepointKO(t *testing.T) {
	db, _ := initTests(t)
	var err error
	ctx := brest.ContextWithDb(context.Background(), db)
	err = brest.ExecuteWithPropagation(ctx, brest.Savepoint, func(ctx context.Context, tx *bun.Tx) error {
		_, err := tx.NewInsert().Model(&Todo{Text: "ko"}).Exec(ctx)
		assert.Nil(t, err)
		return errors.New("ko")
	})
	assert.Nil(t, err)
	count, err := db.NewSelect().Model(&Todo{}).Count(ctx)
	assert.Nil(t, err)
	assert.Equal(t, 0, count)
}

func TestTransactionalSavepointOK(t *testing.T) {
	db, _ := initTests(t)
	var err error
	ctx := brest.ContextWithDb(context.Background(), db)
	err = brest.ExecuteWithPropagation(ctx, brest.Savepoint, func(ctx context.Context, tx *bun.Tx) error {
		_, err := tx.NewInsert().Model(&Todo{Text: "ko"}).Exec(ctx)
		assert.Nil(t, err)
		return err
	})
	assert.Nil(t, err)
	count, err := db.NewSelect().Model(&Todo{}).Count(ctx)
	assert.Nil(t, err)
	assert.Equal(t, 1, count)
}

func TestTransactionalSavepointOKSavepointOK(t *testing.T) {
	db, _ := initTests(t)
	var err error
	ctx := brest.ContextWithDb(context.Background(), db)
	err = brest.ExecuteWithPropagation(ctx, brest.Savepoint, func(ctx context.Context, tx *bun.Tx) error {
		_, err := tx.NewInsert().Model(&Todo{Text: "ok"}).Exec(ctx)
		assert.Nil(t, err)
		return brest.ExecuteWithPropagation(ctx, brest.Savepoint, func(ctx context.Context, tx *bun.Tx) error {
			_, err := tx.NewInsert().Model(&Todo{Text: "ok"}).Exec(ctx)
			assert.Nil(t, err)
			return err
		})
	})
	assert.Nil(t, err)
	count, err := db.NewSelect().Model(&Todo{}).Count(ctx)
	assert.Nil(t, err)
	assert.Equal(t, 2, count)
}

func TestTransactionalSavepointOKSavepointKO(t *testing.T) {
	db, _ := initTests(t)
	var err error
	ctx := brest.ContextWithDb(context.Background(), db)
	err = brest.ExecuteWithPropagation(ctx, brest.Savepoint, func(ctx context.Context, tx *bun.Tx) error {
		_, err := tx.NewInsert().Model(&Todo{Text: "ok"}).Exec(ctx)
		assert.Nil(t, err)
		return brest.ExecuteWithPropagation(ctx, brest.Savepoint, func(ctx context.Context, tx *bun.Tx) error {
			_, err := tx.NewInsert().Model(&Todo{Text: "ko"}).Exec(ctx)
			assert.Nil(t, err)
			return errors.New("ko")
		})
	})
	assert.Nil(t, err)
	count, err := db.NewSelect().Model(&Todo{}).Count(ctx)
	assert.Nil(t, err)
	assert.Equal(t, 1, count)
}

func TestTransactionalCurrentOKSavepointKO(t *testing.T) {
	db, _ := initTests(t)
	var err error
	ctx := brest.ContextWithDb(context.Background(), db)
	err = brest.ExecuteWithPropagation(ctx, brest.Current, func(ctx context.Context, tx *bun.Tx) error {
		_, err := tx.NewInsert().Model(&Todo{Text: "ok"}).Exec(ctx)
		assert.Nil(t, err)
		return brest.ExecuteWithPropagation(ctx, brest.Savepoint, func(ctx context.Context, tx *bun.Tx) error {
			_, err := tx.NewInsert().Model(&Todo{Text: "ko"}).Exec(ctx)
			assert.Nil(t, err)
			return errors.New("ko")
		})
	})
	assert.Nil(t, err)
	count, err := db.NewSelect().Model(&Todo{}).Count(ctx)
	assert.Nil(t, err)
	assert.Equal(t, 1, count)
}

func TestTransactionalCurrentOKSavepointOK(t *testing.T) {
	db, _ := initTests(t)
	var err error
	ctx := brest.ContextWithDb(context.Background(), db)
	err = brest.ExecuteWithPropagation(ctx, brest.Current, func(ctx context.Context, tx *bun.Tx) error {
		_, err := tx.NewInsert().Model(&Todo{Text: "ok"}).Exec(ctx)
		assert.Nil(t, err)
		return brest.ExecuteWithPropagation(ctx, brest.Savepoint, func(ctx context.Context, tx *bun.Tx) error {
			_, err := tx.NewInsert().Model(&Todo{Text: "ok"}).Exec(ctx)
			return err
		})
	})
	assert.Nil(t, err)
	count, err := db.NewSelect().Model(&Todo{}).Count(ctx)
	assert.Nil(t, err)
	assert.Equal(t, 2, count)
}
