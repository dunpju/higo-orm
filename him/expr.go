package him

import "github.com/Masterminds/squirrel"

type IExpression interface {
	Error() error
	ToSql() (string, []interface{}, error)
}

func Expr(sql string, args ...interface{}) squirrel.Sqlizer {
	return squirrel.Expr(sql, args...)
}

type Expression struct {
	err     error
	sqlizer squirrel.Sqlizer
}

func NewExpression(sql string, args []interface{}, err error) IExpression {
	return Expression{sqlizer: Expr(sql, args...), err: err}
}

func (e Expression) Error() error {
	return e.err
}

func (e Expression) ToSql() (string, []interface{}, error) {
	return e.sqlizer.ToSql()
}
