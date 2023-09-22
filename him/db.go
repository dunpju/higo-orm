package him

import (
	"gorm.io/gorm"
)

type DB struct {
	gormDB  *gorm.DB
	connect string
	begin   bool
	Builder interface{}
}

func newDB(db *gorm.DB, connect string) *DB {
	return &DB{gormDB: db, connect: connect}
}

func (this *DB) Query() Select {
	return newSelect(this, this.gormDB)
}

type Select struct {
	db     *DB
	gormDB *gorm.DB
}

func newSelect(db *DB, gormDB *gorm.DB) Select {
	return Select{db: db, gormDB: gormDB}
}

func (this Select) selectBuilder() SelectBuilder {
	if this.db.begin {
		this.db.Builder = newSelectBuilder(this.db.connect).begin(this.gormDB)
	} else {
		this.db.Builder = newSelectBuilder(this.db.connect)
	}
	return this.db.Builder.(SelectBuilder)
}
func (this Select) Distinct() Select {
	this.db.Builder = this.selectBuilder().Distinct()
	return this
}
func (this Select) Select(columns ...string) SelectFrom {
	if len(columns) == 0 {
		columns = append(columns, "*")
	}
	this.db.Builder = this.selectBuilder()._select(columns...)
	return newSelectFrom(this.db, this.gormDB)
}

func (this Select) Raw(pred string, args ...interface{}) SelectRaw {
	this.db.Builder = this.selectBuilder().Raw(pred, args...)
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
	return this.db.Builder.(SelectBuilder).Get(dest)
}

type SelectFrom struct {
	db     *DB
	gormDB *gorm.DB
}

func newSelectFrom(db *DB, gormDB *gorm.DB) SelectFrom {
	return SelectFrom{db: db, gormDB: gormDB}
}

func (this SelectFrom) From(from string) SelectBuilder {
	this.db.Builder = this.db.Builder.(SelectBuilder)._from(from)
	return this.db.Builder.(SelectBuilder)
}

func (this *DB) Insert() InsertInto {
	return newInsertInto(this, this.gormDB)
}

type InsertInto struct {
	db     *DB
	gormDB *gorm.DB
}

func newInsertInto(db *DB, gormDB *gorm.DB) InsertInto {
	return InsertInto{db: db, gormDB: gormDB}
}

func (this InsertInto) Into(from string) InsertBuilder {
	if this.db.begin {
		this.db.Builder = newInsertBuilder(this.db.connect).begin(this.gormDB).insert(from)
	} else {
		this.db.Builder = newInsertBuilder(this.db.connect).insert(from)
	}
	return this.db.Builder.(InsertBuilder)
}

func (this *DB) Update() UpdateTable {
	return newUpdateFrom(this, this.gormDB)
}

type UpdateTable struct {
	db     *DB
	gormDB *gorm.DB
}

func newUpdateFrom(db *DB, gormDB *gorm.DB) UpdateTable {
	return UpdateTable{db: db, gormDB: gormDB}
}

func (this UpdateTable) Table(from string) UpdateBuilder {
	if this.db.begin {
		this.db.Builder = newUpdateBuilder(this.db.connect).begin(this.gormDB).update(from)
	} else {
		this.db.Builder = newUpdateBuilder(this.db.connect).update(from)
	}
	return this.db.Builder.(UpdateBuilder)
}

func (this *DB) Delete() DeleteFrom {
	return newDeleteFrom(this, this.gormDB)
}

type DeleteFrom struct {
	db     *DB
	gormDB *gorm.DB
}

func newDeleteFrom(db *DB, gormDB *gorm.DB) DeleteFrom {
	return DeleteFrom{db: db, gormDB: gormDB}
}

func (this DeleteFrom) From(from string) DeleteBuilder {
	if this.db.begin {
		this.db.Builder = newDeleteBuilder(this.db.connect).begin(this.gormDB).delete(from)
	} else {
		this.db.Builder = newDeleteBuilder(this.db.connect).delete(from)
	}
	return this.db.Builder.(DeleteBuilder)
}

func (this *DB) Begin(tx ...*gorm.DB) *Transaction {
	this.begin = true
	if len(tx) > 0 {
		this.gormDB = tx[0]
	}
	return begin(this, this.gormDB)
}

func (this *DB) TX(tx *gorm.DB) *Transaction {
	return this.Begin(tx)
}

func (this *DB) Set(column string, value interface{}) *DB {
	this.Builder = this.Builder.(UpdateBuilder).Set(column, value)
	return this
}

func (this *DB) GormDB() *gorm.DB {
	return this.gormDB
}

func (this *DB) LastInsertId() (*gorm.DB, int64) {
	if this.begin {
		return this.Builder.(InsertBuilder).begin(this.gormDB).LastInsertId()
	}
	return this.Builder.(InsertBuilder).LastInsertId()
}

func (this *DB) Save() (*gorm.DB, int64) {
	if this.begin {
		return this.Builder.(InsertBuilder).begin(this.gormDB).Save()
	}
	return this.Builder.(InsertBuilder).Save()
}

func (this *DB) Exec() (*gorm.DB, int64) {
	if update, ok := this.Builder.(UpdateBuilder); ok {
		if this.begin {
			return update.begin(this.gormDB).Exec()
		}
		return update.Exec()
	} else if del, ok1 := this.Builder.(DeleteBuilder); ok1 {
		if this.begin {
			return del.begin(this.gormDB).Exec()
		}
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
