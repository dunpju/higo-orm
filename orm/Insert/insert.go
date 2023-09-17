package Insert

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/dunpju/higo-orm/orm"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

func Into(table string) InsertBuilder {
	return into(table)
}

func into(table string) InsertBuilder {
	return InsertBuilder{builder: squirrel.Insert(table)}
}

func Transaction(db *gorm.DB) InsertBuilder {
	return InsertBuilder{db: db}
}

type InsertBuilder struct {
	db      *gorm.DB
	builder squirrel.InsertBuilder
}

func (this InsertBuilder) Into(table string) InsertBuilder {
	this.builder = squirrel.Insert(table)
	return this
}

func (this InsertBuilder) Columns(columns ...string) InsertBuilder {
	this.builder = this.builder.Columns(columns...)
	return this
}

func (this InsertBuilder) Values(values ...interface{}) InsertBuilder {
	this.builder = this.builder.Values(values...)
	return this
}

func (this InsertBuilder) ToSql() (string, []interface{}, error) {
	return this.builder.ToSql()
}

func (this InsertBuilder) LastInsertId() (*gorm.DB, int64) {
	var db *gorm.DB
	if this.db == nil {
		_db, err := orm.Gorm()
		if err != nil {
			_db.Error = err
			return _db, 0
		}
		db = _db
	} else {
		db = this.db
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
	)

	stmt.SQL.WriteString(sql)
	stmt.Vars = args

	result, err := db.Statement.ConnPool.ExecContext(stmt.Context, sql, args...)

	db.Logger.Trace(stmt.Context, curTime, func() (string, int64) {
		sqlStr, vars := stmt.SQL.String(), stmt.Vars
		if filter, ok := db.Logger.(gorm.ParamsFilter); ok {
			sqlStr, vars = filter.ParamsFilter(stmt.Context, stmt.SQL.String(), stmt.Vars...)
		}
		affected, err := result.RowsAffected()
		if err != nil {
			return db.Dialector.Explain(sqlStr, vars...), 0
		}
		return db.Dialector.Explain(sqlStr, vars...), affected
	}, db.Error)

	if err != nil {
		db.Error = err
		return db, 0
	}

	insertID, err := result.LastInsertId()
	insertOk := err == nil && insertID > 0

	if !insertOk {
		db.Error = err
		return db, 0
	}

	return db, insertID
}
