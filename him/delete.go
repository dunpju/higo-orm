package him

import (
	"github.com/Masterminds/squirrel"
	"gorm.io/gorm"
)

func Delete(from ...string) DeleteBuilder {
	var f string
	if len(from) > 0 {
		f = from[0]
	}
	return DeleteBuilder{builder: squirrel.Delete(f), wheres: newWheres()}
}

type DeleteBuilder struct {
	db      *gorm.DB
	connect *connect
	builder squirrel.DeleteBuilder
	wheres  *wheres
	Error   error
}

func newDeleteBuilder(connect string) DeleteBuilder {
	if connect != "" {
		dbc, err := getConnect(connect)
		if err != nil {
			return DeleteBuilder{Error: err}
		}
		return DeleteBuilder{db: dbc.db.GormDB(), connect: dbc, wheres: newWheres()}
	} else {
		dbc, err := getConnect(DefaultConnect)
		if err != nil {
			return DeleteBuilder{Error: err}
		}
		return DeleteBuilder{db: dbc.db.GormDB(), connect: dbc, wheres: newWheres()}
	}
}

func (this DeleteBuilder) DB() *gorm.DB {
	return this.db
}

func (this DeleteBuilder) begin(db *gorm.DB) DeleteBuilder {
	this.db = db
	return this
}

func (this DeleteBuilder) delete(from string) DeleteBuilder {
	this.builder = squirrel.Delete(from)
	return this
}

func (this DeleteBuilder) WhereRaw(fn func(builder WhereRawBuilder) WhereRawBuilder) DeleteBuilder {
	sql, args, err := fn(WhereRawBuilder{}).ToSql()
	this.wheres.and().whereRaw(sql, args, err)
	return this
}

func (this DeleteBuilder) OrWhereRaw(fn func(builder WhereRawBuilder) WhereRawBuilder) DeleteBuilder {
	sql, args, err := fn(WhereRawBuilder{}).ToSql()
	this.wheres.or().whereRaw(sql, args, err)
	return this
}

func (this DeleteBuilder) Where(column, operator string, value interface{}) DeleteBuilder {
	this.wheres.and().where(column, operator, value)
	return this
}

func (this DeleteBuilder) WhereIn(column string, value interface{}) DeleteBuilder {
	this.wheres.and().whereIn(column, value)
	return this
}

func (this DeleteBuilder) WhereNotIn(column string, value interface{}) DeleteBuilder {
	this.wheres.and().whereNotIn(column, value)
	return this
}

func (this DeleteBuilder) WhereNull(column string) DeleteBuilder {
	this.wheres.and().whereNull(column)
	return this
}

func (this DeleteBuilder) WhereNotNull(column string) DeleteBuilder {
	this.wheres.and().whereNotNull(column)
	return this
}

func (this DeleteBuilder) WhereLike(column string, value interface{}) DeleteBuilder {
	this.wheres.and().whereLike(column, value)
	return this
}

func (this DeleteBuilder) NotLike(column string, value interface{}) DeleteBuilder {
	this.wheres.and().whereNotLike(column, value)
	return this
}

func (this DeleteBuilder) WhereBetween(column string, first, second interface{}) DeleteBuilder {
	this.wheres.and().whereBetween(column, first, second)
	return this
}

func (this DeleteBuilder) OrWhere(column, operator string, value interface{}) DeleteBuilder {
	this.wheres.or().where(column, operator, value)
	return this
}

func (this DeleteBuilder) OrWhereIn(column string, value interface{}) DeleteBuilder {
	this.wheres.or().whereIn(column, value)
	return this
}

func (this DeleteBuilder) OrWhereNotIn(column string, value interface{}) DeleteBuilder {
	this.wheres.or().whereNotIn(column, value)
	return this
}

func (this DeleteBuilder) OrWhereNull(column string) DeleteBuilder {
	this.wheres.or().whereNull(column)
	return this
}

func (this DeleteBuilder) OrWhereNotNull(column string) DeleteBuilder {
	this.wheres.or().whereNotNull(column)
	return this
}

func (this DeleteBuilder) OrLike(column string, value interface{}) DeleteBuilder {
	this.wheres.or().whereLike(column, value)
	return this
}

func (this DeleteBuilder) OrNotLike(column string, value interface{}) DeleteBuilder {
	this.wheres.or().whereNotLike(column, value)
	return this
}

func (this DeleteBuilder) OrWhereBetween(column string, first, second interface{}) DeleteBuilder {
	this.wheres.or().whereBetween(column, first, second)
	return this
}

func (this DeleteBuilder) whereHandle() (DeleteBuilder, error) {
	pred, args, err := this.wheres.pred()
	if err != nil {
		return this, err
	}
	this.builder = this.builder.Where(pred, args...)
	return this, nil
}

func (this DeleteBuilder) ToSql() (string, []interface{}, error) {
	builder, err := this.whereHandle()
	if err != nil {
		return "", nil, err
	}
	this = builder
	return this.builder.ToSql()
}

func (this DeleteBuilder) Exec() (*gorm.DB, int64) {
	gormDB, _, rowsAffected := newExecer(this, this.db).exec()
	return gormDB, rowsAffected
}
