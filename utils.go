package brest

import (
	"fmt"
	"reflect"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/schema"
)

func setPk(db *bun.DB, resourceType reflect.Type, elem reflect.Value, key string) error {
	table := db.Table(resourceType)
	if len(table.PKs) == 1 {
		pk := table.PKs[0]
		err := pk.ScanValue(elem, key)
		return err
	}
	return NewErrorBadRequest(fmt.Sprintf("only single pk is permitted for resource '%v'", resourceType))
}

func addQueryLimit(query *bun.SelectQuery, limit int) *bun.SelectQuery {
	if limit == 0 {
		return query
	}
	return query.Limit(int(limit))
}

func addQueryOffset(query *bun.SelectQuery, offset int) *bun.SelectQuery {
	if offset == 0 {
		return query
	}
	return query.Offset(int(offset))
}

func addQueryFields(query *bun.SelectQuery, fields []*Field) *bun.SelectQuery {
	if fields == nil {
		return query
	}
	q := query
	if len(fields) > 0 {
		for _, field := range fields {
			q = q.Column(field.Name)
		}
	}
	return q
}

func addQueryRelations(query *bun.SelectQuery, relations []*Relation) *bun.SelectQuery {
	if relations == nil {
		return query
	}
	q := query
	if len(relations) > 0 {
		for _, relation := range relations {
			q = q.Relation(relation.Name)
		}
	}
	return q
}

func addQuerySorts(query *bun.SelectQuery, sorts []*Sort) *bun.SelectQuery {
	if sorts == nil {
		return query
	}
	q := query
	if len(sorts) > 0 {
		for _, sort := range sorts {
			var orderExpr string
			if sort.Asc {
				orderExpr = sort.Name + " ASC"
			} else {
				orderExpr = sort.Name + " DESC"
			}
			q = q.Order(orderExpr)
		}
	}
	return q
}

func addQueryFilter(query *bun.SelectQuery, filter *Filter, parentGroupOp Op) *bun.SelectQuery {
	if filter == nil {
		return query
	}

	if filter.Op == And || filter.Op == Or {
		return addWhereGroup(query,
			func(query *bun.SelectQuery) *bun.SelectQuery {
				q := query
				for _, subfilter := range filter.Filters {
					q = addQueryFilter(query, subfilter, filter.Op)
				}
				return q
			})
	}

	switch filter.Op {
	case Eq:
		return addWhere(query, "? = ?", filter.Attr, filter.Value, parentGroupOp)
	case Neq:
		return addWhere(query, "? != ?", filter.Attr, filter.Value, parentGroupOp)
	case In:
		return addWhere(query, "? IN (?)", filter.Attr, bun.In(filter.Value), parentGroupOp)
	case Nin:
		return addWhere(query, "? NOT IN (?)", filter.Attr, bun.In(filter.Value), parentGroupOp)
	case Gt:
		return addWhere(query, "? > ?", filter.Attr, filter.Value, parentGroupOp)
	case Gte:
		return addWhere(query, "? >= ?", filter.Attr, filter.Value, parentGroupOp)
	case Lt:
		return addWhere(query, "? < ?", filter.Attr, filter.Value, parentGroupOp)
	case Lte:
		return addWhere(query, "? <= ?", filter.Attr, filter.Value, parentGroupOp)
	case Lk:
		return addWhere(query, "? LIKE ?", filter.Attr, filter.Value, parentGroupOp)
	case Nlk:
		return addWhere(query, "? NOT LIKE ?", filter.Attr, filter.Value, parentGroupOp)
	case Llk:
		return addWhere(query, "lower(?) LIKE lower(?)", filter.Attr, filter.Value, parentGroupOp)
	case Nllk:
		return addWhere(query, "lower(?) NOT LIKE lower(?)", filter.Attr, filter.Value, parentGroupOp)
	case Sim:
		return addWhere(query, "? SIMILAR TO ?", filter.Attr, filter.Value, parentGroupOp)
	case Nsim:
		return addWhere(query, "? NOT SIMILAR TO ?", filter.Attr, filter.Value, parentGroupOp)
	case Lulk:
		return addWhere(query, "lower(unaccent(?)) LIKE lower(unaccent(?))", filter.Attr, filter.Value, parentGroupOp)
	case Nlulk:
		return addWhere(query, "lower(unaccent(?)) NOT LIKE lower(unaccent(?))", filter.Attr, filter.Value, parentGroupOp)
	case Null:
		return addWhere(query, "? IS NULL", filter.Attr, "", parentGroupOp)
	case Nnull:
		return addWhere(query, "? IS NOT NULL", filter.Attr, "", parentGroupOp)
	default:
		return query
	}
}

func addWhere(query *bun.SelectQuery, condition string, attribute string, value interface{}, parentGroupOp Op) *bun.SelectQuery {
	if parentGroupOp == Or {
		return query.WhereOr(condition, schema.Ident(attribute), value)
	}
	return query.Where(condition, schema.Ident(attribute), value)
}

func addWhereGroup(query *bun.SelectQuery, fnGroup func(query *bun.SelectQuery) *bun.SelectQuery) *bun.SelectQuery {
	return query.WhereGroup("", fnGroup)
}
