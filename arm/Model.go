package arm

import (
	"database/sql"
	"github.com/dunpju/higo-orm/him"
	"gorm.io/gorm"
)

type Model struct {
	db      *him.DB
	table   *TableName
	builder any
	err     error
}

func NewModel(db *him.DB, table *TableName) *Model {
	return &Model{db: db, table: table}
}

func (this *Model) Select(columns ...string) him.SelectBuilder {
	return this.db.Query().Select(columns...).From(this.table.String())
}

func (this *Model) Raw(pred string, args ...interface{}) him.RawBuilder {
	return him.NewRawBuilder(this.db, pred, args)
}

func (this *Model) Set(column string, value interface{}) *Model {
	this.db.Set(column, value)
	return this
}

func (this *Model) Insert() (*gorm.DB, int64) {
	return this.db.Insert().Into(this.table.String()).LastInsertId()
}

func (this *Model) Update() (*gorm.DB, int64) {
	return this.db.Update().Table(this.table.String()).Exec()
}

func (this *Model) Delete() (*gorm.DB, int64) {
	return this.db.Delete().From(this.table.String()).Exec()
}

func (this *Model) BeginTX(opts ...*sql.TxOptions) *gorm.DB {
	db, err := him.DBConnect(this.db.Connect())
	if err != nil {
		this.err = err
		return nil
	}
	return db.GormDB().Begin(opts...)
}

func (this *Model) Begin(opts ...*sql.TxOptions) *TX {
	return newTX(this.BeginTX(opts...))
}

func (this *Model) TX(tx *gorm.DB) *Model {
	this.db.TX(tx)
	return this
}

func (this *Model) Error() error {
	return this.err
}

func (this *Model) GormDB() *gorm.DB {
	return this.db.GormDB()
}
