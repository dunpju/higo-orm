package orm

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

func Update(table ...string) UpdateBuilder {
	var t string
	if len(table) > 0 {
		t = table[0]
	}
	return UpdateBuilder{builder: squirrel.Update(t)}
}

type UpdateBuilder struct {
	DB      *gorm.DB
	builder squirrel.UpdateBuilder
}

func (this UpdateBuilder) Update(table ...string) UpdateBuilder {
	return Update(table...)
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

func (this UpdateBuilder) Where(pred interface{}, args ...interface{}) UpdateBuilder {
	this.builder = this.builder.Where(pred, args...)
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

func (this UpdateBuilder) ToSql() (string, []interface{}, error) {
	return this.builder.ToSql()
}

func (this UpdateBuilder) Save() (*gorm.DB, int64) {
	var db *gorm.DB
	if this.DB == nil {
		_db_, err := Gorm()
		if err != nil {
			_db_.Error = err
			return _db_, 0
		}
		db = _db_
	} else {
		fmt.Println(this.DB)
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
