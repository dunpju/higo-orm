package him

import (
	"context"
	"github.com/Masterminds/squirrel"
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

	if err != nil {
		gormDB.Error = err
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		gormDB.Error = err
		return
	}

	affected, err := result.RowsAffected()
	if err != nil {
		gormDB.Error = err
		return
	}
	insertID = id
	rowsAffected = affected
	return
}
