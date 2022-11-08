package transactional_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/aptogeo/brest/transactional"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

type Todo struct {
	ID   int `bun:"id,pk,autoincrement"`
	Text string
}

func initTests(t *testing.T) *bun.DB {
	sqldb, err := sql.Open(sqliteshim.ShimName, "file::memory:?cache=shared")
	if err != nil {
		panic(err)
	}
	db := bun.NewDB(sqldb, sqlitedialect.New())
	db.ResetModel(context.Background(), (*Todo)(nil))
	return db
}

func TestTransactionalCurrentKO(t *testing.T) {
	db := initTests(t)
	ctx := transactional.ContextWithDb(context.Background(), db)
	err := transactional.Execute(ctx, func(ctx context.Context, tx *bun.Tx) error {
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
	db := initTests(t)
	ctx := transactional.ContextWithDb(context.Background(), db)
	err := transactional.Execute(ctx, func(ctx context.Context, tx *bun.Tx) error {
		_, err := tx.NewInsert().Model(&Todo{Text: "ok"}).Exec(ctx)
		return err
	})
	assert.Nil(t, err)
	count, err := db.NewSelect().Model(&Todo{}).Count(ctx)
	assert.Nil(t, err)
	assert.Equal(t, 1, count)
}

func TestTransactionalCurrentOKCurrentKO(t *testing.T) {
	db := initTests(t)
	ctx := transactional.ContextWithDb(context.Background(), db)
	err := transactional.Execute(ctx, func(ctx context.Context, tx *bun.Tx) error {
		_, err := tx.NewInsert().Model(&Todo{Text: "ok"}).Exec(ctx)
		assert.Nil(t, err)
		return transactional.Execute(ctx, func(ctx context.Context, tx *bun.Tx) error {
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
	db := initTests(t)
	ctx := transactional.ContextWithDb(context.Background(), db)
	err := transactional.Execute(ctx, func(ctx context.Context, tx *bun.Tx) error {
		_, err := tx.NewInsert().Model(&Todo{Text: "ok"}).Exec(ctx)
		assert.Nil(t, err)
		return transactional.Execute(ctx, func(ctx context.Context, tx *bun.Tx) error {
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
	db := initTests(t)
	var err error
	ctx := transactional.ContextWithDb(context.Background(), db)
	err = transactional.ExecuteWithPropagation(ctx, transactional.Mandatory, func(ctx context.Context, tx *bun.Tx) error {
		_, err := tx.NewInsert().Model(&Todo{Text: "ok"}).Exec(ctx)
		return err
	})
	assert.NotNil(t, err)
	count, err := db.NewSelect().Model(&Todo{}).Count(ctx)
	assert.Nil(t, err)
	assert.Equal(t, 0, count)

	err = transactional.Execute(ctx, func(ctx context.Context, tx *bun.Tx) error {
		return transactional.ExecuteWithPropagation(ctx, transactional.Current, func(ctx context.Context, tx *bun.Tx) error {
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
	db := initTests(t)
	var err error
	ctx := transactional.ContextWithDb(context.Background(), db)
	err = transactional.ExecuteWithPropagation(ctx, transactional.Savepoint, func(ctx context.Context, tx *bun.Tx) error {
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
	db := initTests(t)
	var err error
	ctx := transactional.ContextWithDb(context.Background(), db)
	err = transactional.ExecuteWithPropagation(ctx, transactional.Savepoint, func(ctx context.Context, tx *bun.Tx) error {
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
	db := initTests(t)
	var err error
	ctx := transactional.ContextWithDb(context.Background(), db)
	err = transactional.ExecuteWithPropagation(ctx, transactional.Savepoint, func(ctx context.Context, tx *bun.Tx) error {
		_, err := tx.NewInsert().Model(&Todo{Text: "ok"}).Exec(ctx)
		assert.Nil(t, err)
		return transactional.ExecuteWithPropagation(ctx, transactional.Savepoint, func(ctx context.Context, tx *bun.Tx) error {
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
	db := initTests(t)
	var err error
	ctx := transactional.ContextWithDb(context.Background(), db)
	err = transactional.ExecuteWithPropagation(ctx, transactional.Savepoint, func(ctx context.Context, tx *bun.Tx) error {
		_, err := tx.NewInsert().Model(&Todo{Text: "ok"}).Exec(ctx)
		assert.Nil(t, err)
		return transactional.ExecuteWithPropagation(ctx, transactional.Savepoint, func(ctx context.Context, tx *bun.Tx) error {
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
	db := initTests(t)
	var err error
	ctx := transactional.ContextWithDb(context.Background(), db)
	err = transactional.ExecuteWithPropagation(ctx, transactional.Current, func(ctx context.Context, tx *bun.Tx) error {
		_, err := tx.NewInsert().Model(&Todo{Text: "ok"}).Exec(ctx)
		assert.Nil(t, err)
		return transactional.ExecuteWithPropagation(ctx, transactional.Savepoint, func(ctx context.Context, tx *bun.Tx) error {
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
	db := initTests(t)
	var err error
	ctx := transactional.ContextWithDb(context.Background(), db)
	err = transactional.ExecuteWithPropagation(ctx, transactional.Current, func(ctx context.Context, tx *bun.Tx) error {
		_, err := tx.NewInsert().Model(&Todo{Text: "ok"}).Exec(ctx)
		assert.Nil(t, err)
		return transactional.ExecuteWithPropagation(ctx, transactional.Savepoint, func(ctx context.Context, tx *bun.Tx) error {
			_, err := tx.NewInsert().Model(&Todo{Text: "ok"}).Exec(ctx)
			return err
		})
	})
	assert.Nil(t, err)
	count, err := db.NewSelect().Model(&Todo{}).Count(ctx)
	assert.Nil(t, err)
	assert.Equal(t, 2, count)
}
