package him

import "gorm.io/gorm"

type RawBuilder struct {
	db    *DB
	pred  string
	args  []interface{}
	table string
}

func NewRawBuilder(db *DB, pred string, args []interface{}, table string) RawBuilder {
	return RawBuilder{db: db, pred: pred, args: args, table: table}
}

func (this RawBuilder) ToSql() (string, []interface{}, error) {
	return this.pred, this.args, nil
}

func (this RawBuilder) Get(dest interface{}) *gorm.DB {
	return this.db.Query().Raw(this.pred, this.args...).Get(dest)
}

func (this RawBuilder) Exec() (gormDB *gorm.DB, insertID int64, affected int64) {
	gormDB, insertID, affected = newExecer(this, this.db.gormDB).exec()
	this.db.gormDB = gormDB
	return this.db.gormDB, insertID, affected
}
