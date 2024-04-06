package him

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/dunpju/higo-orm/event"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type Execer struct {
	sqlizer squirrel.Sqlizer
	gormDB  *gorm.DB
}

func newExecer(sqlizer squirrel.Sqlizer, gormDB *gorm.DB) *Execer {
	return &Execer{sqlizer: sqlizer, gormDB: gormDB}
}

func (this Execer) exec() (gormDB *gorm.DB, insertID int64, rowsAffected int64) {
	gormDB = this.gormDB
	sql, args, err := this.sqlizer.ToSql()

	this.eventBefore(sql, args, err)

	if err != nil {
		gormDB.Error = err
		return
	}

	var (
		curTime = time.Now()
		stmt    = &gorm.Statement{
			DB:       gormDB,
			ConnPool: gormDB.ConnPool,
			Context:  context.Background(),
			Clauses:  map[string]clause.Clause{},
		}
	)

	stmt.SQL.WriteString(sql)
	stmt.Vars = args

	result, err := gormDB.Statement.ConnPool.ExecContext(stmt.Context, sql, args...)

	this.eventAfter(sql, args, err, 0, 0)

	if err != nil {
		gormDB.Error = err
		return
	}

	gormDB.Logger.Trace(stmt.Context, curTime, func() (string, int64) {
		sqlStr, vars := stmt.SQL.String(), stmt.Vars
		if filter, ok := gormDB.Logger.(gorm.ParamsFilter); ok {
			sqlStr, vars = filter.ParamsFilter(stmt.Context, stmt.SQL.String(), stmt.Vars...)
		}
		affected, err1 := result.RowsAffected()
		if err1 != nil {
			return gormDB.Dialector.Explain(sqlStr, vars...), 0
		}
		return gormDB.Dialector.Explain(sqlStr, vars...), affected
	}, gormDB.Error)

	id, err := result.LastInsertId()

	this.eventAfter(sql, args, err, id, 0)

	if err != nil {
		gormDB.Error = err
		return
	}

	affected, err := result.RowsAffected()

	this.eventAfter(sql, args, err, id, affected)

	if err != nil {
		gormDB.Error = err
		return
	}

	insertID = id
	rowsAffected = affected
	return
}

func (this Execer) eventBefore(sql string, args []interface{}, err error) {
	switch this.sqlizer.(type) {
	case *InsertBuilder:
		event.Point(event.BeforeInsert, event.NewEventData(this.sqlizer.(*InsertBuilder).table, sql, args, err, 0, 0))
	case *UpdateBuilder:
		event.Point(event.BeforeUpdate, event.NewEventData(this.sqlizer.(*UpdateBuilder).table, sql, args, err, 0, 0))
	case *DeleteBuilder:
		event.Point(event.BeforeDelete, event.NewEventData(this.sqlizer.(*DeleteBuilder).table, sql, args, err, 0, 0))
	default:
	}
}

func (this Execer) eventAfter(sql string, args []interface{}, err error, lastInsertId int64, rowsAffected int64) {
	switch this.sqlizer.(type) {
	case *InsertBuilder:
		event.Point(event.AfterInsert, event.NewEventData(this.sqlizer.(*InsertBuilder).table, sql, args, err, lastInsertId, rowsAffected))
	case *UpdateBuilder:
		event.Point(event.AfterUpdate, event.NewEventData(this.sqlizer.(*UpdateBuilder).table, sql, args, err, lastInsertId, rowsAffected))
	case *DeleteBuilder:
		event.Point(event.AfterDelete, event.NewEventData(this.sqlizer.(*DeleteBuilder).table, sql, args, err, lastInsertId, rowsAffected))
	default:
	}
}
