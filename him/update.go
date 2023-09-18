package him

import (
	"context"
	"github.com/Masterminds/squirrel"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type UpdateBuilder struct {
	db      *gorm.DB
	connect *connect
	builder squirrel.UpdateBuilder
	wheres  *wheres
	Error   error
}

func NewUpdateBuilder(connect ...string) UpdateBuilder {
	if len(connect) > 0 {
		dbc, err := getConnect(connect[0])
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

func (this UpdateBuilder) DB() *gorm.DB {
	return this.db
}

func (this UpdateBuilder) Transaction(db *gorm.DB) UpdateBuilder {
	this.db = db
	return this
}

func (this UpdateBuilder) Update(table ...string) UpdateBuilder {
	var t string
	if len(table) > 0 {
		t = table[0]
	}
	this.builder = squirrel.Update(t)
	return this
}

func (this UpdateBuilder) Prefix(sql string, args ...interface{}) UpdateBuilder {
	this.builder = this.builder.Prefix(sql, args...)
	return this
}

func (this UpdateBuilder) Table(table string) UpdateBuilder {
	this.builder = this.builder.Table(table)
	return this
}

func (this UpdateBuilder) Set(column string, value interface{}) UpdateBuilder {
	this.builder = this.builder.Set(column, value)
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

func (this UpdateBuilder) Where(column, operator string, value interface{}) UpdateBuilder {
	this.wheres.and().where(column, operator, value)
	return this
}

func (this UpdateBuilder) WhereIn(column string, value interface{}) UpdateBuilder {
	this.wheres.and().whereIn(column, value)
	return this
}

func (this UpdateBuilder) WhereNotIn(column string, value interface{}) UpdateBuilder {
	this.wheres.and().whereNotIn(column, value)
	return this
}

func (this UpdateBuilder) WhereNull(column string) UpdateBuilder {
	this.wheres.and().whereNull(column)
	return this
}

func (this UpdateBuilder) WhereNotNull(column string) UpdateBuilder {
	this.wheres.and().whereNotNull(column)
	return this
}

func (this UpdateBuilder) WhereLike(column string, value interface{}) UpdateBuilder {
	this.wheres.and().whereLike(column, value)
	return this
}

func (this UpdateBuilder) NotLike(column string, value interface{}) UpdateBuilder {
	this.wheres.and().whereNotLike(column, value)
	return this
}

func (this UpdateBuilder) WhereBetween(column string, first, second interface{}) UpdateBuilder {
	this.wheres.and().whereBetween(column, first, second)
	return this
}

func (this UpdateBuilder) OrWhere(column, operator string, value interface{}) UpdateBuilder {
	this.wheres.or().where(column, operator, value)
	return this
}

func (this UpdateBuilder) OrWhereIn(column string, value interface{}) UpdateBuilder {
	this.wheres.or().whereIn(column, value)
	return this
}

func (this UpdateBuilder) OrWhereNotIn(column string, value interface{}) UpdateBuilder {
	this.wheres.or().whereNotIn(column, value)
	return this
}

func (this UpdateBuilder) OrWhereNull(column string) UpdateBuilder {
	this.wheres.or().whereNull(column)
	return this
}

func (this UpdateBuilder) OrWhereNotNull(column string) UpdateBuilder {
	this.wheres.or().whereNotNull(column)
	return this
}

func (this UpdateBuilder) OrLike(column string, value interface{}) UpdateBuilder {
	this.wheres.or().whereLike(column, value)
	return this
}

func (this UpdateBuilder) OrNotLike(column string, value interface{}) UpdateBuilder {
	this.wheres.or().whereNotLike(column, value)
	return this
}

func (this UpdateBuilder) OrWhereBetween(column string, first, second interface{}) UpdateBuilder {
	this.wheres.or().whereBetween(column, first, second)
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

func (this UpdateBuilder) Exec() (*gorm.DB, int64) {
	sql, args, err := this.ToSql()
	if err != nil {
		this.db.Error = err
		return this.db, 0
	}

	var (
		curTime = time.Now()
		stmt    = &gorm.Statement{
			DB:       this.db,
			ConnPool: this.db.ConnPool,
			Context:  context.Background(),
			Clauses:  map[string]clause.Clause{},
		}
		rowsAffected int64
	)

	stmt.SQL.WriteString(sql)
	stmt.Vars = args

	result, err := this.db.Statement.ConnPool.ExecContext(stmt.Context, sql, args...)

	this.db.Logger.Trace(stmt.Context, curTime, func() (string, int64) {
		sqlStr, vars := stmt.SQL.String(), stmt.Vars
		if filter, ok := this.db.Logger.(gorm.ParamsFilter); ok {
			sqlStr, vars = filter.ParamsFilter(stmt.Context, stmt.SQL.String(), stmt.Vars...)
		}
		affected, err1 := result.RowsAffected()
		if err1 != nil {
			return this.db.Dialector.Explain(sqlStr, vars...), 0
		}
		rowsAffected = affected
		return this.db.Dialector.Explain(sqlStr, vars...), affected
	}, this.db.Error)

	if err != nil {
		this.db.Error = err
		return this.db, 0
	}

	return this.db, rowsAffected
}
