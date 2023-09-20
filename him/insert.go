package him

import (
	"github.com/Masterminds/squirrel"
	"gorm.io/gorm"
	"sync"
)

type InsertBuilder struct {
	db         *gorm.DB
	connect    *connect
	setColumns *insertColumn
	setValues  []*insertValue
	affected   int64
	builder    squirrel.InsertBuilder
	Error      error
}

func newDBInsertBuilder(db *gorm.DB, connect *connect) InsertBuilder {
	return InsertBuilder{db: db, connect: connect, setColumns: newInsertColumn(), setValues: make([]*insertValue, 0)}
}

func newErrorInsertBuilder(err error) InsertBuilder {
	return InsertBuilder{Error: err}
}

func newInsertBuilder(connect string) InsertBuilder {
	if connect != "" {
		dbc, err := getConnect(connect)
		if err != nil {
			return newErrorInsertBuilder(err)
		}
		return newDBInsertBuilder(dbc.db.GormDB(), dbc)
	} else {
		dbc, err := getConnect(DefaultConnect)
		if err != nil {
			return newErrorInsertBuilder(err)
		}
		return newDBInsertBuilder(dbc.db.GormDB(), dbc)
	}
}

func (this InsertBuilder) DB() *gorm.DB {
	return this.db
}

func (this InsertBuilder) begin(db *gorm.DB) InsertBuilder {
	this.db = db
	return this
}

func (this InsertBuilder) insert(into string) InsertBuilder {
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
		this = this.values(value.values...)
	}
	this = this.columns(this.setColumns.columns...)
	return this
}

func (this InsertBuilder) ToSql() (string, []interface{}, error) {
	return this.toBuilder().builder.ToSql()
}

func (this InsertBuilder) Save() (*gorm.DB, int64) {
	builder, db, _ := this.save()
	return db, builder.affected
}

func (this InsertBuilder) LastInsertId() (*gorm.DB, int64) {
	_, db, id := this.save()
	return db, id
}

func (this InsertBuilder) save() (InsertBuilder, *gorm.DB, int64) {
	gormDB, insertID, rowsAffected := newExecer(this, this.db).exec()
	if gormDB.Error != nil {
		this.Error = gormDB.Error
		return this, gormDB, 0
	}
	this.db = gormDB
	this.affected = rowsAffected
	return this, this.db, insertID
}
