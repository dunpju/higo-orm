package orm

import (
	"github.com/Masterminds/squirrel"
)

type where struct {
	logic   string
	sqlizer squirrel.Sqlizer
}

func and(column, operator string, value interface{}) where {
	return where{"AND", condition(column, operator, value)}
}

func or(column, operator string, value interface{}) where {
	return where{"OR", condition(column, operator, value)}
}

func condition(column, operator string, value interface{}) squirrel.Sqlizer {
	if operator == ">" {
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
	} else {
		return squirrel.Eq{column: value}
	}
}

func (this SelectBuilder) Where(column, operator string, value interface{}) SelectBuilder {
	this.wheres = append(this.wheres, and(column, operator, value))
	return this
}

func (this SelectBuilder) WhereIn(column string, value interface{}) SelectBuilder {
	this.wheres = append(this.wheres, and(column, "=", value))
	return this
}

func (this SelectBuilder) WhereNotIn(column string, value interface{}) SelectBuilder {
	this.wheres = append(this.wheres, and(column, "<>", value))
	return this
}

func (this SelectBuilder) WhereNull(column string) SelectBuilder {
	this.wheres = append(this.wheres, and(column, "=", nil))
	return this
}

func (this SelectBuilder) WhereNotNull(column string) SelectBuilder {
	this.wheres = append(this.wheres, and(column, "<>", nil))
	return this
}

func (this SelectBuilder) Like(column string, value interface{}) SelectBuilder {
	this.wheres = append(this.wheres, and(column, "Like", value))
	return this
}

func (this SelectBuilder) NotLike(column string, value interface{}) SelectBuilder {
	this.wheres = append(this.wheres, and(column, "NotLike", value))
	return this
}

func (this SelectBuilder) WhereBetween(column string, first, second interface{}) SelectBuilder {
	this.wheres = append(this.wheres, and(column, "BETWEEN", between{column, first, second}))
	return this
}

func (this SelectBuilder) WhereRaw(fn func(builder SelectBuilder) squirrel.Sqlizer) SelectBuilder {
	selectBuilder := query()
	selectBuilder.isWhereRaw = true
	sql, args, err := fn(selectBuilder).ToSql()
	this.wheres = append(this.wheres, and("", "RAW", raw{sql, args, err}))
	return this
}

func (this SelectBuilder) OrWhere(column, operator string, value interface{}) SelectBuilder {
	this.wheres = append(this.wheres, or(column, operator, value))
	return this
}

func (this SelectBuilder) OrWhereNull(column string) SelectBuilder {
	this.wheres = append(this.wheres, or(column, "=", nil))
	return this
}

func (this SelectBuilder) OrWhereNotNull(column string) SelectBuilder {
	this.wheres = append(this.wheres, or(column, "<>", nil))
	return this
}

func (this SelectBuilder) OrWhereIn(column string, value interface{}) SelectBuilder {
	this.wheres = append(this.wheres, or(column, "=", value))
	return this
}

func (this SelectBuilder) OrWhereNotIn(column string, value interface{}) SelectBuilder {
	this.wheres = append(this.wheres, or(column, "<>", value))
	return this
}

func (this SelectBuilder) OrLike(column string, value interface{}) SelectBuilder {
	this.wheres = append(this.wheres, or(column, "Like", value))
	return this
}

func (this SelectBuilder) OrNotLike(column string, value interface{}) SelectBuilder {
	this.wheres = append(this.wheres, or(column, "NotLike", value))
	return this
}

func (this SelectBuilder) OrWhereBetween(column string, first, second interface{}) SelectBuilder {
	this.wheres = append(this.wheres, or(column, "BETWEEN", between{column, first, second}))
	return this
}

func (this SelectBuilder) OrWhereRaw(fn func(builder SelectBuilder) squirrel.Sqlizer) SelectBuilder {
	selectBuilder := query()
	selectBuilder.isWhereRaw = true
	sql, args, err := fn(selectBuilder).ToSql()
	this.wheres = append(this.wheres, or("", "RAW", raw{sql, args, err}))
	return this
}
