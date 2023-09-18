package him

import (
	"context"
	"github.com/Masterminds/squirrel"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type InsertBuilder struct {
	db      *gorm.DB
	connect *connect
	builder squirrel.InsertBuilder
	Error   error
}

func NewInsertBuilder(connect ...string) InsertBuilder {
	if len(connect) > 0 {
		dbc, err := getConnect(connect[0])
		if err != nil {
			return InsertBuilder{Error: err}
		}
		return InsertBuilder{db: dbc.db.GormDB(), connect: dbc}
	} else {
		dbc, err := getConnect(DefaultConnect)
		if err != nil {
			return InsertBuilder{Error: err}
		}
		return InsertBuilder{db: dbc.db.GormDB(), connect: dbc}
	}
}

func (this InsertBuilder) DB() *gorm.DB {
	return this.db
}

func (this InsertBuilder) Transaction(db *gorm.DB) InsertBuilder {
	this.db = db
	return this
}

func (this InsertBuilder) Insert(into string) InsertBuilder {
	this.builder = squirrel.Insert(into)
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
		return this.db.Dialector.Explain(sqlStr, vars...), affected
	}, this.db.Error)

	if err != nil {
		this.db.Error = err
		return this.db, 0
	}

	insertID, err := result.LastInsertId()
	insertOk := err == nil && insertID > 0

	if !insertOk {
		this.db.Error = err
		return this.db, 0
	}

	return this.db, insertID
}
