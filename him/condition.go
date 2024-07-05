package him

import (
	"fmt"
	"github.com/Masterminds/squirrel"
)

func and(column, operator string, value interface{}) where {
	return where{AND, conditionHandle(column, operator, value)}
}

func or(column, operator string, value interface{}) where {
	return where{OR, conditionHandle(column, operator, value)}
}

func conditionHandle(column, operator string, value interface{}) squirrel.Sqlizer {
	if expr, ok := value.(squirrel.Sqlizer); ok {
		exprSql, exprArgs, err := expr.ToSql()
		sql := fmt.Sprintf("%s %s %s", column, operator, exprSql)
		for _, _ = range exprArgs {
			sql += " ?"
		}
		return NewExpression(sql, exprArgs, err)
	} else if operator == ">" {
		return squirrel.Gt{column: value}
	} else if operator == ">=" {
		return squirrel.GtOrEq{column: value}
	} else if operator == "<" {
		return squirrel.Lt{column: value}
	} else if operator == "<=" {
		return squirrel.LtOrEq{column: value}
	} else if operator == "<>" || operator == "!=" {
		return squirrel.NotEq{column: value}
	} else if operator == "Like" {
		return squirrel.Like{column: value}
	} else if operator == "NotLike" {
		return squirrel.NotLike{column: value}
	} else if operator == "BETWEEN" {
		return value.(between)
	} else if operator == "RAW" {
		return value.(raw)
	} else if operator == "whereRaw" {
		return value.(whereRaw)
	} else if operator == "IN" {
		return squirrel.Eq{column: value}
	} else if operator == "NotIn" {
		return squirrel.NotEq{column: value}
	} else if operator == "Null" {
		return squirrel.Eq{column: value}
	} else if operator == "NotNull" {
		return squirrel.NotEq{column: value}
	} else {
		return squirrel.Eq{column: value}
	}
}
