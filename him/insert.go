package him

import (
	"context"
	"github.com/Masterminds/squirrel"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

func Insert(into string) InsertBuilder {
	return InsertBuilder{builder: squirrel.Insert(into)}
}

type InsertBuilder struct {
	DB      *gorm.DB
	builder squirrel.InsertBuilder
}

func (this InsertBuilder) Insert(into string) InsertBuilder {
	builder := Insert(into)
	this.builder = builder.builder
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
