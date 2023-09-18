package him

import (
	"context"
	"github.com/Masterminds/squirrel"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

func Delete(from ...string) DeleteBuilder {
	var f string
	if len(from) > 0 {
		f = from[0]
	}
	return DeleteBuilder{builder: squirrel.Delete(f), wheres: newWheres()}
}

type DeleteBuilder struct {
	DB      *gorm.DB
	builder squirrel.DeleteBuilder
	wheres  *wheres
}

func (this DeleteBuilder) Delete(from ...string) DeleteBuilder {
	builder := Delete(from...)
	this.builder = builder.builder
	this.wheres = builder.wheres
	return this
}

func (this DeleteBuilder) From(from string) DeleteBuilder {
	this.builder = this.builder.From(from)
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
	var db *gorm.DB
	if this.DB == nil {
		_db_, err := Gorm()
		if err != nil {
			_db_.Error = err
			return _db_, 0
		}
		db = _db_
	} else {
		db = this.DB
	}

	sql, args, err := this.ToSql()
	if err != nil {
		db.Error = err
		return db, 0
	}

	var (
		curTime = time.Now()
		stmt    = &gorm.Statement{
			DB:       db,
			ConnPool: db.ConnPool,
			Context:  context.Background(),
			Clauses:  map[string]clause.Clause{},
		}
		rowsAffected int64
	)

	stmt.SQL.WriteString(sql)
	stmt.Vars = args

	result, err := db.Statement.ConnPool.ExecContext(stmt.Context, sql, args...)

	db.Logger.Trace(stmt.Context, curTime, func() (string, int64) {
		sqlStr, vars := stmt.SQL.String(), stmt.Vars
		if filter, ok := db.Logger.(gorm.ParamsFilter); ok {
			sqlStr, vars = filter.ParamsFilter(stmt.Context, stmt.SQL.String(), stmt.Vars...)
		}
		affected, err1 := result.RowsAffected()
		if err1 != nil {
			return db.Dialector.Explain(sqlStr, vars...), 0
		}
		rowsAffected = affected
		return db.Dialector.Explain(sqlStr, vars...), affected
	}, db.Error)

	if err != nil {
		db.Error = err
		return db, 0
	}

	return db, rowsAffected
}