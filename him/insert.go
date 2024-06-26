package him

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"gorm.io/gorm"
	"strings"
	"sync"
)

type InsertBuilder struct {
	db                      *DB
	connect                 *connect
	setColumns              *insertColumn
	setValues               []*insertValue
	setOnDuplicateKeyUpdate []string
	affected                int64
	builder                 squirrel.InsertBuilder
	table                   string
	Error                   error
}

func newDBInsertBuilder(db *DB, connect *connect) *InsertBuilder {
	insertBuilder := &InsertBuilder{db: db, connect: connect, setColumns: newInsertColumn(), setValues: make([]*insertValue, 0)}
	insertBuilder.db.Builder = insertBuilder
	return insertBuilder
}

func newErrorInsertBuilder(err error) *InsertBuilder {
	return &InsertBuilder{Error: err}
}

func newInsertBuilder(db *DB) *InsertBuilder {
	if conn, ok := _connect.Load(db.connect); ok {
		return newDBInsertBuilder(db, conn.(*connect))
	} else {
		return newErrorInsertBuilder(fmt.Errorf("db connect nonexistent"))
	}
}

func (this *InsertBuilder) DB() *gorm.DB {
	return this.db.GormDB()
}

func (this *InsertBuilder) TX(tx *gorm.DB) *InsertBuilder {
	this = this.begin(tx)
	return this
}

func (this *InsertBuilder) begin(db *gorm.DB) *InsertBuilder {
	this.db.gormDB = db
	return this
}

func (this *InsertBuilder) insert(into string) *InsertBuilder {
	this.builder = squirrel.Insert(into)
	this.table = into
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

func (this *InsertBuilder) Columns(columns ...any) *ValuesBuilder {
	this.setColumns.add(columnsToString(columns...)...)
	return newValuesBuilder(this)
}

func (this *InsertBuilder) columns(columns ...string) *InsertBuilder {
	this.builder = this.builder.Columns(columns...)
	return this
}

type ValuesBuilder struct {
	insertBuilder *InsertBuilder
}

func newValuesBuilder(insertBuilder *InsertBuilder) *ValuesBuilder {
	return &ValuesBuilder{insertBuilder: insertBuilder}
}

func (this *ValuesBuilder) Values(values ...interface{}) *ValuesBuilder {
	this.insertBuilder.setValues = append(this.insertBuilder.setValues, newInsertValue(values...))
	return this
}

// OnDuplicateKeyUpdate age = values(age) 或者 age = 10 注意: values()括号里是字段名称
func (this *ValuesBuilder) OnDuplicateKeyUpdate(values ...interface{}) *ValuesBuilder {
	this.insertBuilder.OnDuplicateKeyUpdate(values...)
	return this
}

func (this *ValuesBuilder) Save() (gormDB *gorm.DB, affected int64) {
	return this.insertBuilder.Save()
}

func (this *ValuesBuilder) ToSql() (string, []interface{}, error) {
	return this.insertBuilder.ToSql()
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

func (this *InsertBuilder) values(values ...interface{}) *InsertBuilder {
	this.builder = this.builder.Values(values...)
	return this
}

func (this *InsertBuilder) Column(column any, value interface{}) *InsertBuilder {
	return this.Set(columnToString(column), value)
}

func (this *InsertBuilder) Set(column any, value interface{}) *InsertBuilder {
	this.setColumns.add(columnToString(column))
	if len(this.setValues) > 0 {
		this.setValues[0].value(value)
	} else {
		this.setValues = append(this.setValues, newInsertValue(value))
	}
	this.db.Builder = this
	return this
}

// OnDuplicateKeyUpdate age = values(age) 或者 age = 10 注意: values()括号里是字段名称
// 注意: 高并发情况可能会出现死锁
// 成功插入新增数据返回1
// 主键或唯一索引已经重复，插入失败，但更新的字段的数据和原数据不同，则表示更新成功，此时返回的结果是2
// 更新数据和原始数据相同，则表示更新也没有成功，则返回0
func (this *InsertBuilder) OnDuplicateKeyUpdate(values ...interface{}) *InsertBuilder {
	this.setOnDuplicateKeyUpdate = append(this.setOnDuplicateKeyUpdate, toStrings(values...)...)
	return this
}

func (this *InsertBuilder) toBuilder() *InsertBuilder {
	for _, value := range this.setValues {
		this = this.values(value.values...)
	}
	this = this.columns(this.setColumns.columns...)
	return this
}

func (this *InsertBuilder) ToSql() (string, []interface{}, error) {
	this.builder = this.toBuilder().builder
	if len(this.setOnDuplicateKeyUpdate) > 0 {
		sql, args, err := this.builder.ToSql()
		return fmt.Sprintf("%s ON DUPLICATE KEY UPDATE %s", sql, strings.Join(this.setOnDuplicateKeyUpdate, ",")), args, err
	}
	return this.builder.ToSql()
}

func (this *InsertBuilder) Save() (gormDB *gorm.DB, affected int64) {
	builder, db, _ := this.save()
	return db, builder.affected
}

func (this *InsertBuilder) LastInsertId() (*gorm.DB, int64) {
	_, db, id := this.save()
	return db, id
}

func (this *InsertBuilder) save() (*InsertBuilder, *gorm.DB, int64) {
	gormDB, insertID, rowsAffected := newExecer(this, this.db.GormDB()).exec()
	if gormDB.Error != nil {
		this.Error = gormDB.Error
		return this, gormDB, 0
	}
	this.db.gormDB = gormDB
	this.affected = rowsAffected
	return this, this.db.gormDB, insertID
}
