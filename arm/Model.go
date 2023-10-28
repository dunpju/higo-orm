package arm

import (
	"database/sql"
	"github.com/dunpju/higo-orm/him"
	"gorm.io/gorm"
)

type Model struct {
	db      *him.DB
	table   *TableName
	model   IModel
	builder any
	err     error
}

func Connect(model IModel) error {
	db, err := him.DBConnect(model.Connection())
	if err != nil {
		return err
	}
	model.Apply(newModel(db, model))
	return nil
}

func newModel(db *him.DB, model IModel) *Model {
	return &Model{db: db, model: model, table: model.TableName()}
}

func (this *Model) DB() *him.DB {
	return this.db
}

func (this *Model) Property(properties ...him.IProperty) {
	him.Properties(properties).Apply(this.model)
}

func (this *Model) Connection() string {
	return him.DefaultConnect
}

func (this *Model) Alias(alias string) *Model {
	this.table.Alias(alias)
	return this
}

func (this *Model) Select(columns ...string) him.SelectBuilder {
	return this.db.Query().Select(columns...).From(this.table.String())
}

func (this *Model) Raw(pred string, args ...interface{}) him.RawBuilder {
	return him.NewRawBuilder(this.db, pred, args)
}

func (this *Model) Insert() him.InsertBuilder {
	this.builder = this.db.Insert().Into(this.table.String())
	return this.builder.(him.InsertBuilder)
}

func (this *Model) Update() him.UpdateBuilder {
	this.builder = this.db.Update().Table(this.table.String())
	return this.builder.(him.UpdateBuilder)
}

func (this *Model) Delete() him.DeleteBuilder {
	this.builder = this.db.Delete().From(this.table.String())
	return this.builder.(him.DeleteBuilder)
}

func (this *Model) Set(column any, value interface{}) *Model {
	if insertBuilder, ok := this.builder.(him.InsertBuilder); ok {
		this.builder = insertBuilder.Set(column, value)
	} else if updateBuilder, ok := this.builder.(him.UpdateBuilder); ok {
		this.builder = updateBuilder.Set(column, value)
	}
	return this
}

func (this *Model) BeginTX(opts ...*sql.TxOptions) *gorm.DB {
	db, err := him.DBConnect(this.db.Connect())
	if err != nil {
		this.err = err
		return nil
	}
	return db.GormDB().Begin(opts...)
}

func (this *Model) Begin(opts ...*sql.TxOptions) *him.TX {
	return him.NewTX(this.BeginTX(opts...))
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
