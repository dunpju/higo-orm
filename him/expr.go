package him

import "github.com/Masterminds/squirrel"

func Expr(sql string, args ...interface{}) squirrel.Sqlizer {
	return squirrel.Expr(sql, args...)
}
