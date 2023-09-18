package orm

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
	return DeleteBuilder{builder: squirrel.Delete(f)}
}

type DeleteBuilder struct {
	DB      *gorm.DB
	builder squirrel.DeleteBuilder
}

func (this DeleteBuilder) Delete(from ...string) DeleteBuilder {
	var f string
	if len(from) > 0 {
		f = from[0]
	}
	this.builder = squirrel.Delete(f)
	return this
}

func (this DeleteBuilder) Prefix(sql string, args ...interface{}) DeleteBuilder {
	this.builder = this.builder.Prefix(sql, args...)
	return this
}

func (this DeleteBuilder) From(from string) DeleteBuilder {
	this.builder = this.builder.From(from)
	return this
}

func (this DeleteBuilder) Where(pred interface{}, args ...interface{}) DeleteBuilder {
	this.builder = this.builder.Where(pred, args...)
	return this
}

func (this DeleteBuilder) OrderBy(orderBys ...string) DeleteBuilder {
	this.builder = this.builder.OrderBy(orderBys...)
	return this
}

func (this DeleteBuilder) Limit(limit uint64) DeleteBuilder {
	this.builder = this.builder.Limit(limit)
	return this
}

func (this DeleteBuilder) Offset(offset uint64) DeleteBuilder {
	this.builder = this.builder.Offset(offset)
	return this
}

func (this DeleteBuilder) Suffix(sql string, args ...interface{}) DeleteBuilder {
	this.builder = this.builder.Suffix(sql, args...)
	return this
}

func (this DeleteBuilder) ToSql() (string, []interface{}, error) {
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
