package orm

import (
	"github.com/Masterminds/squirrel"
)

func (this SelectBuilder) whereHandle(selectBuilder squirrel.SelectBuilder, wheres *wheres) (squirrel.SelectBuilder, error) {
	pred, args, err := wheres.pred()
	if err != nil {
		return squirrel.SelectBuilder{}, err
	}
	selectBuilder = selectBuilder.Where(pred, args...)
	return selectBuilder, nil
}

type where struct {
	logic   Logic
	sqlizer squirrel.Sqlizer
}

func (this SelectBuilder) Raw(pred string, args ...interface{}) SelectBuilder {
	this.isRaw = true
	this.wheres.and().raw(pred, args, nil)
	return this
}

func (this SelectBuilder) OrRaw(pred string, args ...interface{}) SelectBuilder {
	this.isRaw = true
	this.wheres.or().raw(pred, args, nil)
	return this
}

func (this SelectBuilder) WhereRaw(fn func(builder WhereRawBuilder) WhereRawBuilder) SelectBuilder {
	sql, args, err := fn(WhereRawBuilder{}).ToSql()
	this.wheres.and().whereRaw(sql, args, err)
	return this
}

func (this SelectBuilder) OrWhereRaw(fn func(builder WhereRawBuilder) WhereRawBuilder) SelectBuilder {
	sql, args, err := fn(WhereRawBuilder{}).ToSql()
	this.wheres.or().whereRaw(sql, args, err)
	return this
}

func (this SelectBuilder) Where(column, operator string, value interface{}) SelectBuilder {
	this.wheres.and().where(column, operator, value)
	return this
}

func (this SelectBuilder) WhereIn(column string, value interface{}) SelectBuilder {
	this.wheres.and().whereIn(column, value)
	return this
}

func (this SelectBuilder) WhereNotIn(column string, value interface{}) SelectBuilder {
	this.wheres.and().whereNotIn(column, value)
	return this
}

func (this SelectBuilder) WhereNull(column string) SelectBuilder {
	this.wheres.and().whereNull(column)
	return this
}

func (this SelectBuilder) WhereNotNull(column string) SelectBuilder {
	this.wheres.and().whereNotNull(column)
	return this
}

func (this SelectBuilder) WhereLike(column string, value interface{}) SelectBuilder {
	this.wheres.and().whereLike(column, value)
	return this
}

func (this SelectBuilder) NotLike(column string, value interface{}) SelectBuilder {
	this.wheres.and().whereNotLike(column, value)
	return this
}

func (this SelectBuilder) WhereBetween(column string, first, second interface{}) SelectBuilder {
	this.wheres.and().whereBetween(column, first, second)
	return this
}

func (this SelectBuilder) OrWhere(column, operator string, value interface{}) SelectBuilder {
	this.wheres.or().where(column, operator, value)
	return this
}

func (this SelectBuilder) OrWhereIn(column string, value interface{}) SelectBuilder {
	this.wheres.or().whereIn(column, value)
	return this
}

func (this SelectBuilder) OrWhereNotIn(column string, value interface{}) SelectBuilder {
	this.wheres.or().whereNotIn(column, value)
	return this
}

func (this SelectBuilder) OrWhereNull(column string) SelectBuilder {
	this.wheres.or().whereNull(column)
	return this
}

func (this SelectBuilder) OrWhereNotNull(column string) SelectBuilder {
	this.wheres.or().whereNotNull(column)
	return this
}

func (this SelectBuilder) OrLike(column string, value interface{}) SelectBuilder {
	this.wheres.or().whereLike(column, value)
	return this
}

func (this SelectBuilder) OrNotLike(column string, value interface{}) SelectBuilder {
	this.wheres.or().whereNotLike(column, value)
	return this
}

func (this SelectBuilder) OrWhereBetween(column string, first, second interface{}) SelectBuilder {
	this.wheres.or().whereBetween(column, first, second)
	return this
}
