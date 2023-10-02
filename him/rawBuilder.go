package him

import "gorm.io/gorm"

type RawBuilder struct {
	db   *DB
	pred string
	args []interface{}
}

func NewRawBuilder(db *DB, pred string, args []interface{}) RawBuilder {
	return RawBuilder{db: db, pred: pred, args: args}
}

func (this RawBuilder) ToSql() (string, []interface{}, error) {
	return this.pred, this.args, nil
}

func (this RawBuilder) Get(dest interface{}) *gorm.DB {
	return this.db.Query().Raw(this.pred, this.args...).Get(dest)
}

func (this RawBuilder) Exec() (gormDB *gorm.DB, insertID int64, rowsAffected int64) {
	gormDB, insertID, rowsAffected = newExecer(this, this.db.gormDB).exec()
	this.db.gormDB = gormDB
	return
}
