/**
 * 此代码基于 gorm.io/gorm/clause.OrderBy，用于 ORDER BY 语句中的 NULL 值问题
 *
 * MySQL 和 SQLite 的行为相当于认为 NULL 小于所有非 NULL 值
 * PostgreSQL 的行为则相当于认为 NULL 大于所有非 NULL 值，但可以通过 NULLS FIRST 或 NULLS LAST 修改此行为
 * @see: https://dev.mysql.com/doc/refman/5.7/en/problems-with-null.html
 * @see: https://www.postgresql.org/docs/13/indexes-ordering.html
 * @see: https://www.sqlite.org/lang_createindex.html#nulls_first_and_nulls_last
 */
package db

import (
	"gorm.io/gorm/clause"
)

type NullsType string

//goland:noinspection GoUnusedConst
const (
	NullsFirst NullsType = " NULLS FIRST"
	NullsLast  NullsType = " NULLS LAST"
)

type OrderByColumn struct {
	Column  clause.Column
	Desc    bool
	Reorder bool
	Nulls   NullsType
}

type OrderBy struct {
	Columns    []OrderByColumn
	Expression clause.Expression
}

// Name where clause name
func (orderBy OrderBy) Name() string {
	return "ORDER BY"
}

// Build build where clause
func (orderBy OrderBy) Build(builder clause.Builder) {
	if orderBy.Expression != nil {
		orderBy.Expression.Build(builder)
	} else {
		for idx, column := range orderBy.Columns {
			if idx > 0 {
				_ = builder.WriteByte(',')
			}

			builder.WriteQuoted(column.Column)
			if column.Desc {
				_, _ = builder.WriteString(" DESC")
			}
			if column.Nulls != "" {
				_, _ = builder.WriteString(string(column.Nulls))
			}
		}
	}
}

// MergeClause merge order by clauses
func (orderBy OrderBy) MergeClause(c *clause.Clause) {
	if v, ok := c.Expression.(OrderBy); ok {
		for i := len(orderBy.Columns) - 1; i >= 0; i-- {
			if orderBy.Columns[i].Reorder {
				orderBy.Columns = orderBy.Columns[i:]
				c.Expression = orderBy
				return
			}
		}

		copiedColumns := make([]OrderByColumn, len(v.Columns))
		copy(copiedColumns, v.Columns)
		orderBy.Columns = append(copiedColumns, orderBy.Columns...)
	}

	c.Expression = orderBy
}
