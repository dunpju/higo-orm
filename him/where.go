package him

import (
	"github.com/Masterminds/squirrel"
)

func (this *SelectBuilder) whereHandle(selectBuilder squirrel.SelectBuilder, wheres *wheres) (squirrel.SelectBuilder, error) {
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

func (this *SelectBuilder) Raw(pred string, args ...interface{}) *SelectBuilder {
	this.isRaw = true
	this.wheres.and().raw(pred, args, nil)
	return this
}

func (this *SelectBuilder) OrRaw(pred string, args ...interface{}) *SelectBuilder {
	this.isRaw = true
	this.wheres.or().raw(pred, args, nil)
	return this
}

func (this *SelectBuilder) WhereRaw(fn func(builder WhereRawBuilder) WhereRawBuilder) *SelectBuilder {
	sql, args, err := fn(WhereRawBuilder{}).ToSql()
	this.wheres.and().whereRaw(sql, args, err)
	return this
}

func (this *SelectBuilder) OrWhereRaw(fn func(builder WhereRawBuilder) WhereRawBuilder) *SelectBuilder {
	sql, args, err := fn(WhereRawBuilder{}).ToSql()
	this.wheres.or().whereRaw(sql, args, err)
	return this
}

func (this *SelectBuilder) Where(column any, operator string, value interface{}) *SelectBuilder {
	this.wheres.and().where(columnToString(column), operator, value)
	return this
}

func (this *SelectBuilder) WhereIn(column any, value interface{}) *SelectBuilder {
	this.wheres.and().whereIn(columnToString(column), value)
	return this
}

func (this *SelectBuilder) WhereNotIn(column any, value interface{}) *SelectBuilder {
	this.wheres.and().whereNotIn(columnToString(column), value)
	return this
}

func (this *SelectBuilder) WhereNull(column any) *SelectBuilder {
	this.wheres.and().whereNull(columnToString(column))
	return this
}

func (this *SelectBuilder) WhereNotNull(column any) *SelectBuilder {
	this.wheres.and().whereNotNull(columnToString(column))
	return this
}

func (this *SelectBuilder) WhereLike(column any, value interface{}) *SelectBuilder {
	this.wheres.and().whereLike(columnToString(column), value)
	return this
}

func (this *SelectBuilder) NotLike(column any, value interface{}) *SelectBuilder {
	this.wheres.and().whereNotLike(columnToString(column), value)
	return this
}

func (this *SelectBuilder) WhereBetween(column any, first, second interface{}) *SelectBuilder {
	this.wheres.and().whereBetween(columnToString(column), first, second)
	return this
}

func (this *SelectBuilder) OrWhere(column any, operator string, value interface{}) *SelectBuilder {
	this.wheres.or().where(columnToString(column), operator, value)
	return this
}

func (this *SelectBuilder) OrWhereIn(column any, value interface{}) *SelectBuilder {
	this.wheres.or().whereIn(columnToString(column), value)
	return this
}

func (this *SelectBuilder) OrWhereNotIn(column any, value interface{}) *SelectBuilder {
	this.wheres.or().whereNotIn(columnToString(column), value)
	return this
}

func (this *SelectBuilder) OrWhereNull(column any) *SelectBuilder {
	this.wheres.or().whereNull(columnToString(column))
	return this
}

func (this *SelectBuilder) OrWhereNotNull(column any) *SelectBuilder {
	this.wheres.or().whereNotNull(columnToString(column))
	return this
}

func (this *SelectBuilder) OrLike(column any, value interface{}) *SelectBuilder {
	this.wheres.or().whereLike(columnToString(column), value)
	return this
}

func (this *SelectBuilder) OrNotLike(column any, value interface{}) *SelectBuilder {
	this.wheres.or().whereNotLike(columnToString(column), value)
	return this
}

func (this *SelectBuilder) OrWhereBetween(column any, first, second interface{}) *SelectBuilder {
	this.wheres.or().whereBetween(columnToString(column), first, second)
	return this
}
