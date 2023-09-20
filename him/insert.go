package him

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"sync"
	"time"
)

type InsertBuilder struct {
	db         *gorm.DB
	connect    *connect
	setColumns *insertColumn
	setValues  []*insertValue
	builder    squirrel.InsertBuilder
	Error      error
}

func newInsertBuilder(db *gorm.DB, connect *connect) InsertBuilder {
	return InsertBuilder{db: db, connect: connect, setColumns: newInsertColumn(), setValues: make([]*insertValue, 0)}
}

func newErrorInsertBuilder(err error) InsertBuilder {
	return InsertBuilder{Error: err}
}

func NewInsertBuilder(connect ...string) InsertBuilder {
	if len(connect) > 0 {
		dbc, err := getConnect(connect[0])
		if err != nil {
			return newErrorInsertBuilder(err)
		}
		return newInsertBuilder(dbc.db.GormDB(), dbc)
	} else {
		dbc, err := getConnect(DefaultConnect)
		if err != nil {
			return newErrorInsertBuilder(err)
		}
		return newInsertBuilder(dbc.db.GormDB(), dbc)
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

type insertColumn struct {
	columns   []string
	columnMap sync.Map
}

func newInsertColumn() *insertColumn {
	return &insertColumn{columns: []string{}}
}

func (this *insertColumn) add(columns ...string) {
	if len(columns) > 0 {
		for _, c := range columns {
			if _, ok := this.columnMap.Load(c); !ok {
				this.columnMap.Store(c, c)
				this.columns = append(this.columns, c)
			}
		}
	}
}

func (this *insertColumn) len() int {
	return len(this.columns)
}

func (this InsertBuilder) Columns(columns ...string) InsertBuilder {
	this.setColumns.add(columns...)
	return this
}

func (this InsertBuilder) columns(columns ...string) InsertBuilder {
	this.builder = this.builder.Columns(columns...)
	return this
}

type insertValue struct {
	values []interface{}
}

func newInsertValue(values ...interface{}) *insertValue {
	return &insertValue{values: values}
}

func (this *insertValue) value(values ...interface{}) {
	if len(values) > 0 {
		this.values = append(this.values, values...)
	}
}

func (this *insertValue) len() int {
	return len(this.values)
}

func (this InsertBuilder) Values(values ...interface{}) InsertBuilder {
	this.setValues = append(this.setValues, newInsertValue(values...))
	return this
}

func (this InsertBuilder) values(values ...interface{}) InsertBuilder {
	this.builder = this.builder.Values(values...)
	return this
}

func (this InsertBuilder) Set(column string, value interface{}) InsertBuilder {
	this.setColumns.add(column)
	if len(this.setValues) > 0 {
		this.setValues[0].value(value)
	} else {
		this.setValues = append(this.setValues, newInsertValue(value))
	}
	return this
}

func (this InsertBuilder) toBuilder() InsertBuilder {
	for _, value := range this.setValues {
		fmt.Println(value.values...)
		this = this.values(value.values...)
	}
	this = this.columns(this.setColumns.columns...)
	return this
}

func (this InsertBuilder) ToSql() (string, []interface{}, error) {
	return this.toBuilder().builder.ToSql()
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
