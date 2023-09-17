package orm

import (
	"github.com/Masterminds/squirrel"
	"gorm.io/gorm"
	"time"
)

func Insert(into string) InsertBuilder {
	return insert(into)
}

func insert(into string) InsertBuilder {
	return InsertBuilder{builder: squirrel.Insert(into)}
}

type InsertBuilder struct {
	builder squirrel.InsertBuilder
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
	db, err := Gorm()
	if err != nil {
		db.Error = err
		return db, 0
	}
	sql, args, err := this.ToSql()
	if err != nil {
		db.Error = err
		return db, 0
	}

	var (
		curTime = time.Now()
		stmt    = db.Statement
	)

	stmt.SQL.WriteString(sql)
	stmt.Vars = args

	db.Logger.Trace(stmt.Context, curTime, func() (string, int64) {
		sqlStr, vars := stmt.SQL.String(), stmt.Vars
		if filter, ok := db.Logger.(gorm.ParamsFilter); ok {
			sqlStr, vars = filter.ParamsFilter(stmt.Context, stmt.SQL.String(), stmt.Vars...)
		}
		return db.Dialector.Explain(sqlStr, vars...), db.RowsAffected
	}, db.Error)

	result, err := db.Statement.ConnPool.ExecContext(stmt.Context, sql, args...)
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
