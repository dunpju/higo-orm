package him

import (
	"gorm.io/gorm"
)

type DB struct {
	gormDB  *gorm.DB
	connect string
	dbc     *DBConfig
	begin   bool
	db      *DB
	Builder interface{}
	Error   error
}

func newDB(gormDB *gorm.DB, connect string, dbc *DBConfig, begin bool) *DB {
	return &DB{gormDB: gormDB, connect: connect, dbc: dbc, begin: begin}
}

func (this *DB) Connect() string {
	return this.connect
}

func (this *DB) Rest() *DB {
	this.gormDB.Error = nil
	this.Error = nil
	return this
}

func (this *DB) DBC() *DBConfig {
	return this.dbc
}

func (this *DB) Begin(tx ...*gorm.DB) *Transaction {
	this.begin = true
	if len(tx) > 0 {
		this.gormDB = tx[0]
		return begin(this, this.gormDB)
	}
	return begin(this, nil)
}

func (this *DB) TX(tx *gorm.DB) *Transaction {
	return this.Begin(tx)
}

func (this *DB) Set(column string, value interface{}) *DB {
	if updateBuilder, ok := this.db.Builder.(*UpdateBuilder); ok {
		this.db.Builder = updateBuilder.Set(column, value)
	} else if insertBuilder, ok := this.db.Builder.(*InsertBuilder); ok {
		this.db.Builder = insertBuilder.Column(column, value)
	}
	return this
}

func (this *DB) GormDB() *gorm.DB {
	return this.gormDB
}

func (this *DB) LastInsertId() (*gorm.DB, int64) {
	if this.begin {
		return this.db.Builder.(*InsertBuilder).begin(this.gormDB).LastInsertId()
	}
	return this.db.Builder.(*InsertBuilder).LastInsertId()
}

func (this *DB) Save() (*gorm.DB, int64) {
	if this.begin {
		return this.db.Builder.(*InsertBuilder).begin(this.gormDB).Save()
	}
	return this.db.Builder.(*InsertBuilder).Save()
}

func (this *DB) Exec() (*gorm.DB, int64) {
	if update, updateOk := this.db.Builder.(*UpdateBuilder); updateOk {
		if this.begin {
			return update.begin(this.gormDB).Exec()
		}
		return update.Exec()
	} else if del, deleteOk := this.db.Builder.(*DeleteBuilder); deleteOk {
		if this.begin {
			return del.begin(this.gormDB).Exec()
		}
		return del.Exec()
	}
	return nil, 0
}

func (this *DB) Raw(pred string, args ...interface{}) Raw {
	return newRaw(this, this.gormDB, pred, args)
}

type Raw struct {
	db     *DB
	gormDB *gorm.DB
	pred   string
	args   []interface{}
}

func newRaw(db *DB, gormDB *gorm.DB, pred string, args []interface{}) Raw {
	return Raw{db: db, gormDB: gormDB, pred: pred, args: args}
}

func (this Raw) Exec() (gormDB *gorm.DB, insertID int64, rowsAffected int64) {
	if this.db.begin {
		gormDB, insertID, rowsAffected = newExecer(newSelectBuilder(this.db.connect).begin(this.gormDB).Raw(this.pred, this.args...), this.gormDB).exec()
	} else {
		gormDB, insertID, rowsAffected = newExecer(newSelectBuilder(this.db.connect).Raw(this.pred, this.args...), this.gormDB).exec()
	}
	this.db.gormDB = gormDB
	return
}

func (this Raw) Get(dest interface{}) *gorm.DB {
	return newSelectBuilder(this.db.connect).Raw(this.pred, this.args...).Get(dest)
}
