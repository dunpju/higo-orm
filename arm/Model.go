package arm

import (
	"github.com/dunpju/higo-orm/him"
	"gorm.io/gorm"
)

type Model struct {
	db      *him.DB
	table   *TableName
	builder any
}

func NewModel(db *him.DB, table *TableName) *Model {
	return &Model{db: db, table: table}
}

func (this *Model) Select(columns ...string) him.SelectBuilder {
	return this.db.Query().Select(columns...).From(this.table.String())
}

func (this *Model) Raw(pred string, args ...interface{}) him.SelectRaw {
	return this.db.Query().Raw(pred, args...)
}

func (this *Model) Set(column string, value interface{}) *Model {
	this.db.Set(column, value)
	return this
}

func (this *Model) Update() (*gorm.DB, int64) {
	return this.db.Update().Table(this.table.String()).Exec()
}

func (this *Model) Delete() (*gorm.DB, int64) {
	return this.db.Delete().From(this.table.String()).Exec()
}
