package him

import (
	"github.com/Masterminds/squirrel"
	"gorm.io/gorm"
)

type UpdateBuilder struct {
	db      *gorm.DB
	connect *connect
	builder squirrel.UpdateBuilder
	wheres  *wheres
	Error   error
}

func newUpdateBuilder(connect string) UpdateBuilder {
	if connect != "" {
		dbc, err := getConnect(connect)
		if err != nil {
			return UpdateBuilder{Error: err}
		}
		return UpdateBuilder{db: dbc.db.GormDB(), connect: dbc, wheres: newWheres()}
	} else {
		dbc, err := getConnect(DefaultConnect)
		if err != nil {
			return UpdateBuilder{Error: err}
		}
		return UpdateBuilder{db: dbc.db.GormDB(), connect: dbc, wheres: newWheres()}
	}
}

func (this UpdateBuilder) update(table string) UpdateBuilder {
	this.builder = squirrel.Update(table)
	return this
}

func (this UpdateBuilder) DB() *gorm.DB {
	return this.db
}

func (this UpdateBuilder) TX(tx *gorm.DB) UpdateBuilder {
	this = this.begin(tx)
	return this
}

func (this UpdateBuilder) begin(db *gorm.DB) UpdateBuilder {
	this.db = db
	return this
}

func (this UpdateBuilder) Prefix(sql string, args ...interface{}) UpdateBuilder {
	this.builder = this.builder.Prefix(sql, args...)
	return this
}

func (this UpdateBuilder) Set(column any, value interface{}) UpdateBuilder {
	this.builder = this.builder.Set(columnToString(column), value)
	return this
}

func (this UpdateBuilder) SetMap(clauses map[string]interface{}) UpdateBuilder {
	this.builder = this.builder.SetMap(clauses)
	return this
}

func (this UpdateBuilder) From(from string) UpdateBuilder {
	this.builder = this.builder.From(from)
	return this
}

func (this UpdateBuilder) OrderBy(orderBys ...string) UpdateBuilder {
	this.builder = this.builder.OrderBy(orderBys...)
	return this
}

func (this UpdateBuilder) Limit(limit uint64) UpdateBuilder {
	this.builder = this.builder.Limit(limit)
	return this
}

func (this UpdateBuilder) Offset(offset uint64) UpdateBuilder {
	this.builder = this.builder.Offset(offset)
	return this
}

func (this UpdateBuilder) Suffix(sql string, args ...interface{}) UpdateBuilder {
	this.builder = this.builder.Suffix(sql, args...)
	return this
}

func (this UpdateBuilder) WhereRaw(fn func(builder WhereRawBuilder) WhereRawBuilder) UpdateBuilder {
	sql, args, err := fn(WhereRawBuilder{}).ToSql()
	this.wheres.and().whereRaw(sql, args, err)
	return this
}

func (this UpdateBuilder) OrWhereRaw(fn func(builder WhereRawBuilder) WhereRawBuilder) UpdateBuilder {
	sql, args, err := fn(WhereRawBuilder{}).ToSql()
	this.wheres.or().whereRaw(sql, args, err)
	return this
}

func (this UpdateBuilder) Where(column any, operator string, value interface{}) UpdateBuilder {
	this.wheres.and().where(columnToString(column), operator, value)
	return this
}

func (this UpdateBuilder) WhereIn(column any, value interface{}) UpdateBuilder {
	this.wheres.and().whereIn(columnToString(column), value)
	return this
}

func (this UpdateBuilder) WhereNotIn(column any, value interface{}) UpdateBuilder {
	this.wheres.and().whereNotIn(columnToString(column), value)
	return this
}

func (this UpdateBuilder) WhereNull(column any) UpdateBuilder {
	this.wheres.and().whereNull(columnToString(column))
	return this
}

func (this UpdateBuilder) WhereNotNull(column any) UpdateBuilder {
	this.wheres.and().whereNotNull(columnToString(column))
	return this
}

func (this UpdateBuilder) WhereLike(column any, value interface{}) UpdateBuilder {
	this.wheres.and().whereLike(columnToString(column), value)
	return this
}

func (this UpdateBuilder) NotLike(column any, value interface{}) UpdateBuilder {
	this.wheres.and().whereNotLike(columnToString(column), value)
	return this
}

func (this UpdateBuilder) WhereBetween(column any, first, second interface{}) UpdateBuilder {
	this.wheres.and().whereBetween(columnToString(column), first, second)
	return this
}

func (this UpdateBuilder) OrWhere(column any, operator string, value interface{}) UpdateBuilder {
	this.wheres.or().where(columnToString(column), operator, value)
	return this
}

func (this UpdateBuilder) OrWhereIn(column any, value interface{}) UpdateBuilder {
	this.wheres.or().whereIn(columnToString(column), value)
	return this
}

func (this UpdateBuilder) OrWhereNotIn(column any, value interface{}) UpdateBuilder {
	this.wheres.or().whereNotIn(columnToString(column), value)
	return this
}

func (this UpdateBuilder) OrWhereNull(column any) UpdateBuilder {
	this.wheres.or().whereNull(columnToString(column))
	return this
}

func (this UpdateBuilder) OrWhereNotNull(column any) UpdateBuilder {
	this.wheres.or().whereNotNull(columnToString(column))
	return this
}

func (this UpdateBuilder) OrLike(column any, value interface{}) UpdateBuilder {
	this.wheres.or().whereLike(columnToString(column), value)
	return this
}

func (this UpdateBuilder) OrNotLike(column any, value interface{}) UpdateBuilder {
	this.wheres.or().whereNotLike(columnToString(column), value)
	return this
}

func (this UpdateBuilder) OrWhereBetween(column any, first, second interface{}) UpdateBuilder {
	this.wheres.or().whereBetween(columnToString(column), first, second)
	return this
}

func (this UpdateBuilder) whereHandle() (UpdateBuilder, error) {
	pred, args, err := this.wheres.pred()
	if err != nil {
		return this, err
	}
	this.builder = this.builder.Where(pred, args...)
	return this, nil
}

func (this UpdateBuilder) ToSql() (string, []interface{}, error) {
	builder, err := this.whereHandle()
	if err != nil {
		return "", nil, err
	}
	this = builder
	return this.builder.ToSql()
}

func (this UpdateBuilder) Exec() (gormDB *gorm.DB, affected int64) {
	db, _, rowsAffected := newExecer(this, this.db).exec()
	this.db = db
	return this.db, rowsAffected
}
