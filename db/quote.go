package db

import (
	"strings"
)

// QuoteForLikeEscape 转义字符串用于 LIKE 查询
// 例：
//
//	db.First(&row, "column LIKE ? ESCAPE ?", "%" + types.QuoteForLikeEscape(keyword, '\\') + "%", "\\")
//
// @see: https://www.sqlite.org/lang_expr.html#like
// @see: https://dev.mysql.com/doc/refman/5.7/en/string-comparison-functions.html#operator_like
// @see: https://www.postgresql.org/docs/12/functions-matching.html#FUNCTIONS-LIKE
func QuoteForLikeEscape(s string, escapeChar rune) string {
	escape := string(escapeChar)
	s = strings.ReplaceAll(s, escape, escape+escape)
	if escapeChar != '_' {
		s = strings.ReplaceAll(s, "_", escape+"_")
	}
	if escapeChar != '%' {
		s = strings.ReplaceAll(s, "%", escape+"%")
	}
	return s
}
