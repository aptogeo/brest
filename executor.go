package brest

import (
	"context"

	"github.com/aptogeo/brest/transactional"
	"github.com/uptrace/bun"
)

// Executor structure
type Executor struct {
	config             *Config
	restQuery          *RestQuery
	entity             interface{}
	debug              bool
	count              int
	originalSearchPath string
}

// NewExecutor constructs Executor
func NewExecutor(config *Config, restQuery *RestQuery, entity interface{}, debug bool) *Executor {
	e := new(Executor)
	e.config = config
	e.restQuery = restQuery
	e.entity = entity
	e.debug = debug
	e.count = 0
	return e
}

// ExecuteWithSearchPath executes with search path
func (e *Executor) Execute(ctx context.Context, execFunc transactional.ExecFunc) error {
	var err error
	err = transactional.Execute(ctx, func(ctx context.Context, tx *bun.Tx) error {
		return execFunc(ctx, tx)
	})
	return err
}

// GetOneExecFunc gets one execution function
func (e *Executor) GetOneExecFunc() transactional.ExecFunc {
	return func(ctx context.Context, tx *bun.Tx) error {
		q := tx.NewSelect().Model(e.entity).WherePK()
		q = addQueryFields(q, e.restQuery.Fields)
		q = addQueryRelations(q, e.restQuery.Relations)
		if e.debug {
			// TODO
		}
		count, err := q.ScanAndCount(ctx)
		if err != nil {
			return NewErrorFromCause(e.restQuery, err)
		}
		e.count = count
		return nil
	}
}

// GetSliceExecFunc gets slice execution function
func (e *Executor) GetSliceExecFunc() transactional.ExecFunc {
	return func(ctx context.Context, tx *bun.Tx) error {
		var err error
		q := tx.NewSelect().Model(e.entity)
		q = addQueryLimit(q, e.restQuery.Limit)
		q = addQueryOffset(q, e.restQuery.Offset)
		q = addQueryFields(q, e.restQuery.Fields)
		q = addQuerySorts(q, e.restQuery.Sorts)
		q = addQueryFilter(q, e.restQuery.Filter, And)
		if e.debug {
			// TODO
		}
		e.count, err = q.ScanAndCount(ctx)
		if err != nil {
			return NewErrorFromCause(e.restQuery, err)
		}
		return nil
	}
}

// InsertExecFunc inserts execution function
func (e *Executor) InsertExecFunc() transactional.ExecFunc {
	return func(ctx context.Context, tx *bun.Tx) error {
		q := tx.NewInsert().Model(e.entity)
		if e.debug {
			// TODO
		}
		if _, err := q.Exec(ctx); err != nil {
			return NewErrorFromCause(e.restQuery, err)
		}
		e.count = 1
		return nil
	}
}

// UpdateExecFunc updates execution function
func (e *Executor) UpdateExecFunc() transactional.ExecFunc {
	return func(ctx context.Context, tx *bun.Tx) error {
		q := tx.NewUpdate().Model(e.entity).WherePK()
		if e.debug {
			// TODO
		}
		if _, err := q.Exec(ctx); err != nil {
			return NewErrorFromCause(e.restQuery, err)
		}
		e.count = 1
		return nil
	}
}

// DeleteExecFunc deletes execution function
func (e *Executor) DeleteExecFunc() transactional.ExecFunc {
	return func(ctx context.Context, tx *bun.Tx) error {
		q := tx.NewDelete().Model(e.entity).WherePK()
		if e.debug {
			// TODO
		}
		if _, err := q.Exec(ctx); err != nil {
			return NewErrorFromCause(e.restQuery, err)
		}
		e.count = 1
		return nil
	}
}
