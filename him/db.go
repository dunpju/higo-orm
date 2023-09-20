package him

import (
	"gorm.io/gorm"
)

type DB struct {
	gormDB  *gorm.DB
	connect string
	builder interface{}
}

func newDB(db *gorm.DB, connect string) *DB {
	return &DB{gormDB: db, connect: connect}
}

func (this *DB) Query() SelectBuilder {
	this.builder = newSelectBuilder(this.connect)
	return this.builder.(SelectBuilder)
}

func (this *DB) Insert(into string) InsertBuilder {
	this.builder = newInsertBuilder(this.connect).Insert(into)
	return this.builder.(InsertBuilder)
}

func (this *DB) Update(table string) UpdateBuilder {
	this.builder = newUpdateBuilder(this.connect).Update(table)
	return this.builder.(UpdateBuilder)
}

func (this *DB) Delete(from string) DeleteBuilder {
	this.builder = newDeleteBuilder(this.connect).Delete(from)
	return this.builder.(DeleteBuilder)
}

func (this *DB) Begin(tx ...*gorm.DB) *Transaction {
	return begin(this.connect, tx...)
}

func (this *DB) TX(tx ...*gorm.DB) *Transaction {
	return this.Begin(tx...)
}

func (this *DB) GormDB() *gorm.DB {
	return this.gormDB
}

func (this *DB) First(dest interface{}) *gorm.DB {
	return this.builder.(SelectBuilder).First(dest)
}

func (this *DB) Get(dest interface{}) *gorm.DB {
	return this.builder.(SelectBuilder).Get(dest)
}

func (this *DB) Paginate(page, perPage uint64, dest interface{}) (*gorm.DB, Paginate) {
	return this.builder.(SelectBuilder).Paginate(page, perPage, dest)
}

func (this *DB) Count() (*gorm.DB, int64) {
	return this.builder.(SelectBuilder).Count()
}

func (this *DB) Sum(column string) (*gorm.DB, uint64) {
	return this.builder.(SelectBuilder).Sum(column)
}

func (this *DB) LastInsertId() (*gorm.DB, int64) {
	return this.builder.(InsertBuilder).LastInsertId()
}

func (this *DB) Exec() (*gorm.DB, int64) {
	if update, ok := this.builder.(UpdateBuilder); ok {
		return update.Exec()
	} else if del, ok1 := this.builder.(DeleteBuilder); ok1 {
		return del.Exec()
	}
	return nil, 0
}
