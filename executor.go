package brest

import (
	"context"

	"github.com/uptrace/bun"
)

// Executor structure
type Executor struct {
	config    *Config
	restQuery *RestQuery
	entity    interface{}
	count     int
	total     int
}

// NewExecutor constructs Executor
func NewExecutor(config *Config, restQuery *RestQuery, entity interface{}) *Executor {
	e := new(Executor)
	e.config = config
	e.restQuery = restQuery
	e.entity = entity
	e.count = 0
	return e
}

// Execute executes query
func (e *Executor) Execute(ctx context.Context, execFunc ExecFunc) error {
	var err error
	err = Execute(ctx, func(ctx context.Context, tx *bun.Tx) error {
		return execFunc(ctx, tx)
	})
	return err
}

// GetOneExecFunc gets one execution function
func (e *Executor) GetOneExecFunc() ExecFunc {
	return func(ctx context.Context, tx *bun.Tx) error {
		q := tx.NewSelect().Model(e.entity).WherePK()
		q = addQueryFields(q, e.restQuery.Fields)
		q = addQueryRelations(q, e.restQuery.Relations)
		count, err := q.ScanAndCount(ctx)
		if err != nil {
			return NewErrorFromCause(err)
		}
		e.count = count
		return nil
	}
}

// GetSliceExecFunc gets slice execution function
func (e *Executor) GetSliceExecFunc() ExecFunc {
	return func(ctx context.Context, tx *bun.Tx) error {
		var err error
		q := tx.NewSelect().Model(e.entity)
		q = addQueryLimit(q, e.restQuery.Limit)
		q = addQueryOffset(q, e.restQuery.Offset)
		q = addQueryFields(q, e.restQuery.Fields)
		q = addQuerySorts(q, e.restQuery.Sorts)
		q = addQueryFilter(q, e.restQuery.Filter, And)
		e.count, err = q.ScanAndCount(ctx)
		if err != nil {
			return NewErrorFromCause(err)
		}
		return nil
	}
}

// InsertExecFunc inserts execution function
func (e *Executor) InsertExecFunc() ExecFunc {
	return func(ctx context.Context, tx *bun.Tx) error {
		q := tx.NewInsert().Model(e.entity)
		if _, err := q.Exec(ctx); err != nil {
			return NewErrorFromCause(err)
		}
		e.count = 1
		return nil
	}
}

// UpdateExecFunc updates execution function
func (e *Executor) UpdateExecFunc() ExecFunc {
	return func(ctx context.Context, tx *bun.Tx) error {
		q := tx.NewUpdate().Model(e.entity).WherePK()
		if _, err := q.Exec(ctx); err != nil {
			return NewErrorFromCause(err)
		}
		e.count = 1
		return nil
	}
}

// DeleteExecFunc deletes execution function
func (e *Executor) DeleteExecFunc() ExecFunc {
	return func(ctx context.Context, tx *bun.Tx) error {
		q := tx.NewDelete().Model(e.entity).WherePK()
		if _, err := q.Exec(ctx); err != nil {
			return NewErrorFromCause(err)
		}
		e.count = 1
		return nil
	}
}
