package him

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/dunpju/higo-orm/event"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
	"time"
)

type Execer struct {
	sqlizer squirrel.Sqlizer
	gormDB  *gorm.DB
}

func newExecer(sqlizer squirrel.Sqlizer, gormDB *gorm.DB) *Execer {
	return &Execer{sqlizer: sqlizer, gormDB: gormDB}
}

func handleArgs(args []interface{}) ([]interface{}, error) {
	ret := make([]interface{}, 0)
	for _, arg := range args {
		if expr, ok := arg.(squirrel.Sqlizer); ok {
			exprSql, exprArgs, err := expr.ToSql()
			if err != nil {
				return ret, err
			}
			for _, exprArg := range exprArgs {
				exprSql = strings.Replace(exprSql, "?", fmt.Sprintf("%v", exprArg), 1)
			}
			ret = append(ret, exprSql)
		} else {
			ret = append(ret, arg)
		}
	}
	return ret, nil
}

func (this Execer) exec() (gormDB *gorm.DB, insertID int64, rowsAffected int64) {
	gormDB = this.gormDB
	sql, args, err := this.sqlizer.ToSql()
	if err != nil {
		gormDB.Error = err
		this.eventBefore(sql, args, err)
		return
	}

	args, err = handleArgs(args)
	if err != nil {
		gormDB.Error = err
		this.eventBefore(sql, args, err)
		return
	}
	this.eventBefore(sql, args, err)

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

	if err != nil {
		this.eventAfter(sql, args, err, 0, 0)
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

	if err != nil {
		this.eventAfter(sql, args, err, id, 0)
		gormDB.Error = err
		return
	}

	affected, err := result.RowsAffected()

	if err != nil {
		this.eventAfter(sql, args, err, id, affected)
		gormDB.Error = err
		return
	}

	this.eventAfter(sql, args, err, id, affected)

	insertID = id
	rowsAffected = affected
	return
}

func (this Execer) eventBefore(sql string, args []interface{}, err error) {
	switch this.sqlizer.(type) {
	case *InsertBuilder:
		event.Point(event.BeforeInsert, event.NewEventRecord(this.sqlizer.(*InsertBuilder).table, sql, args, err, 0, 0))
	case *UpdateBuilder:
		event.Point(event.BeforeUpdate, event.NewEventRecord(this.sqlizer.(*UpdateBuilder).table, sql, args, err, 0, 0))
	case *DeleteBuilder:
		event.Point(event.BeforeDelete, event.NewEventRecord(this.sqlizer.(*DeleteBuilder).table, sql, args, err, 0, 0))
	case RawBuilder:
		event.Point(event.BeforeRaw, event.NewEventRecord(this.sqlizer.(RawBuilder).table, sql, args, err, 0, 0))
	default:
	}
}

func (this Execer) eventAfter(sql string, args []interface{}, err error, lastInsertId int64, rowsAffected int64) {
	switch this.sqlizer.(type) {
	case *InsertBuilder:
		event.Point(event.AfterInsert, event.NewEventRecord(this.sqlizer.(*InsertBuilder).table, sql, args, err, lastInsertId, rowsAffected))
	case *UpdateBuilder:
		event.Point(event.AfterUpdate, event.NewEventRecord(this.sqlizer.(*UpdateBuilder).table, sql, args, err, lastInsertId, rowsAffected))
	case *DeleteBuilder:
		event.Point(event.AfterDelete, event.NewEventRecord(this.sqlizer.(*DeleteBuilder).table, sql, args, err, lastInsertId, rowsAffected))
	case RawBuilder:
		event.Point(event.AfterRaw, event.NewEventRecord(this.sqlizer.(RawBuilder).table, sql, args, err, lastInsertId, rowsAffected))
	default:
	}
}
