package him

import (
	"gorm.io/gorm"
)

type DB struct {
	gormDB  *gorm.DB
	connect string
	begin   bool
	builder interface{}
}

func newDB(db *gorm.DB, connect string) *DB {
	return &DB{gormDB: db, connect: connect}
}

func (this *DB) Query() SelectBuilder {
	this.builder = newSelectBuilder(this.connect)
	return this.builder.(SelectBuilder)
}

func (this *DB) Insert() InsertInto {
	return newInsertInto(this, this.gormDB)
}

type InsertInto struct {
	db      *DB
	gormDB  *gorm.DB
	builder InsertBuilder
}

func newInsertInto(db *DB, gormDB *gorm.DB) InsertInto {
	return InsertInto{db: db, gormDB: gormDB}
}

func (this InsertInto) Into(from string) InsertBuilder {
	if this.db.begin {
		this.builder = newInsertBuilder(this.db.connect).begin(this.gormDB).insert(from)
	} else {
		this.builder = newInsertBuilder(this.db.connect).insert(from)
	}
	return this.builder
}

func (this *DB) Update() UpdateTable {
	return newUpdateFrom(this, this.gormDB)
}

type UpdateTable struct {
	db      *DB
	gormDB  *gorm.DB
	builder UpdateBuilder
}

func newUpdateFrom(db *DB, gormDB *gorm.DB) UpdateTable {
	return UpdateTable{db: db, gormDB: gormDB}
}

func (this UpdateTable) Table(from string) UpdateBuilder {
	if this.db.begin {
		this.db.builder = newUpdateBuilder(this.db.connect).begin(this.gormDB).update(from)
	} else {
		this.db.builder = newUpdateBuilder(this.db.connect).update(from)
	}
	return this.db.builder.(UpdateBuilder)
}

func (this *DB) Delete() DeleteFrom {
	return newDeleteFrom(this, this.gormDB)
}

type DeleteFrom struct {
	db      *DB
	gormDB  *gorm.DB
	builder DeleteBuilder
}

func newDeleteFrom(db *DB, gormDB *gorm.DB) DeleteFrom {
	return DeleteFrom{db: db, gormDB: gormDB}
}

func (this DeleteFrom) From(from string) DeleteBuilder {
	if this.db.begin {
		this.db.builder = newDeleteBuilder(this.db.connect).begin(this.gormDB).delete(from)
	} else {
		this.db.builder = newDeleteBuilder(this.db.connect).delete(from)
	}
	return this.db.builder.(DeleteBuilder)
}

func (this *DB) Begin(tx ...*gorm.DB) *Transaction {
	this.begin = true
	return begin(this, tx...)
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

func (this *DB) Save() (*gorm.DB, int64) {
	return this.builder.(InsertBuilder).Save()
}

func (this *DB) Exec() (*gorm.DB, int64) {
	if update, ok := this.builder.(UpdateBuilder); ok {
		return update.Exec()
	} else if del, ok1 := this.builder.(DeleteBuilder); ok1 {
		return del.Exec()
	}
	return nil, 0
}
