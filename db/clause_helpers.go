package db

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//goland:noinspection GoUnusedExportedFunction
func ColumnExpressions(columns []clause.Column, additionExprs ...clause.Expression) (exprs []clause.Expression) {
	exprs = make([]clause.Expression, 0, len(columns)+len(additionExprs))
	for _, column := range columns {
		exprs = append(exprs, clause.Expr{
			SQL:  "?",
			Vars: []interface{}{column},
		})
	}
	exprs = append(exprs, additionExprs...)
	return
}

//goland:noinspection GoUnusedExportedFunction
func ClauseSelect(iExprs any) (sel clause.Select) {
	switch exprs := iExprs.(type) {
	case []clause.Column:
		sel.Columns = exprs
	case []clause.Expression:
		sel.Expression = clause.CommaExpression{Exprs: exprs}
	}
	return
}

//goland:noinspection GoUnusedExportedFunction
func ClauseSelectDistinct(iExprs any) (sel clause.Select) {
	sel = ClauseSelect(iExprs)
	sel.Distinct = true
	return
}

//goland:noinspection GoUnusedExportedFunction
func ClauseScopeWhere(conditions []clause.Expression) func(*gorm.DB) *gorm.DB {
	return func(query *gorm.DB) *gorm.DB {
		if len(conditions) > 0 {
			return query.Clauses(clause.Where{Exprs: conditions})
		} else {
			return query
		}
	}
}

type LikeEscape struct {
	Column any
	Value  any
	Escape string
}

func (like LikeEscape) Build(builder clause.Builder) {
	builder.WriteQuoted(like.Column)
	_, _ = builder.WriteString(" LIKE ")
	builder.AddVar(builder, like.Value)
	if like.Escape != "" {
		_, _ = builder.WriteString(" ESCAPE ")
		builder.AddVar(builder, like.Escape)
	}
}

func (like LikeEscape) NegationBuild(builder clause.Builder) {
	builder.WriteQuoted(like.Column)
	_, _ = builder.WriteString(" NOT LIKE ")
	builder.AddVar(builder, like.Value)
	if like.Escape != "" {
		_, _ = builder.WriteString(" ESCAPE ")
		builder.AddVar(builder, like.Escape)
	}
}

//goland:noinspection GoUnusedExportedFunction
func ClauseExprForContain(tx *gorm.DB, column any, keyword string) clause.Expression {
	switch tx.Dialector.Name() {
	case "sqlite", "dm":
		return LikeEscape{
			Column: column,
			Value:  "%" + QuoteForLikeEscape(keyword, '\\') + "%",
			Escape: "\\",
		}
	default:
		return clause.Like{
			Column: column,
			Value:  "%" + QuoteForLikeEscape(keyword, '\\') + "%",
		}
	}
}

//goland:noinspection GoUnusedExportedFunction
func ClauseExprForStartsWith(tx *gorm.DB, column any, keyword string) clause.Expression {
	switch tx.Dialector.Name() {
	case "sqlite", "dm":
		return LikeEscape{
			Column: column,
			Value:  QuoteForLikeEscape(keyword, '\\') + "%",
			Escape: "\\",
		}
	default:
		return clause.Like{
			Column: column,
			Value:  QuoteForLikeEscape(keyword, '\\') + "%",
		}
	}
}

//goland:noinspection GoUnusedExportedFunction
func TableNameOf(tx *gorm.DB, tableModel any) (tableName string) {
	stmt := gorm.Statement{DB: tx}
	if err := stmt.Parse(tableModel); err == nil {
		tableName = stmt.Table
	}
	return
}
