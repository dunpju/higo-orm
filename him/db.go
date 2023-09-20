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

func (this *DB) Query() Select {
	return newSelect(this, this.gormDB)
}

type Select struct {
	db         *DB
	gormDB     *gorm.DB
	selectFrom SelectFrom
}

func newSelect(db *DB, gormDB *gorm.DB) Select {
	return Select{db: db, gormDB: gormDB}
}

func (this Select) selectBuilder() SelectBuilder {
	if this.db.begin {
		this.db.builder = newSelectBuilder(this.db.connect).begin(this.gormDB)
	} else {
		this.db.builder = newSelectBuilder(this.db.connect)
	}
	return this.db.builder.(SelectBuilder)
}
func (this Select) Distinct() Select {
	this.db.builder = this.selectBuilder().Distinct()
	return this
}
func (this Select) Select(columns ...string) SelectFrom {
	this.db.builder = this.selectBuilder().Select(columns...)
	return newSelectFrom(this.db, this.gormDB)
}

func (this Select) Raw(pred string, args ...interface{}) SelectRaw {
	this.db.builder = this.selectBuilder().Raw(pred, args...)
	return newSelectRaw(this.db, this.gormDB)
}

type SelectRaw struct {
	db     *DB
	gormDB *gorm.DB
}

func newSelectRaw(db *DB, gormDB *gorm.DB) SelectRaw {
	return SelectRaw{db: db, gormDB: gormDB}
}

func (this SelectRaw) Get(dest interface{}) *gorm.DB {
	return this.db.builder.(SelectBuilder).Get(dest)
}

type SelectFrom struct {
	db     *DB
	gormDB *gorm.DB
}

func newSelectFrom(db *DB, gormDB *gorm.DB) SelectFrom {
	return SelectFrom{db: db, gormDB: gormDB}
}

func (this SelectFrom) From(from string) SelectBuilder {
	this.db.builder = this.db.builder.(SelectBuilder).From(from)
	return this.db.builder.(SelectBuilder)
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
	var t *gorm.DB
	if len(tx) > 0 {
		t = tx[0]
	}
	return begin(this, t)
}

func (this *DB) TX(tx *gorm.DB) *Transaction {
	return this.Begin(tx)
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

func (this *DB) Raw(pred string, args ...interface{}) ExecRaw {
	return newExecRaw(this, this.gormDB, pred, args)
}

type ExecRaw struct {
	db     *DB
	gormDB *gorm.DB
	pred   string
	args   []interface{}
}

func newExecRaw(db *DB, gormDB *gorm.DB, pred string, args []interface{}) ExecRaw {
	return ExecRaw{db: db, gormDB: gormDB, pred: pred, args: args}
}

func (this ExecRaw) Exec() (gormDB *gorm.DB, insertID int64, rowsAffected int64) {
	if this.db.begin {
		gormDB, insertID, rowsAffected = newExecer(newSelectBuilder(this.db.connect).begin(this.gormDB).Raw(this.pred, this.args...), this.gormDB).exec()
	} else {
		gormDB, insertID, rowsAffected = newExecer(newSelectBuilder(this.db.connect).Raw(this.pred, this.args...), this.gormDB).exec()
	}
	this.db.gormDB = gormDB
	return
}
