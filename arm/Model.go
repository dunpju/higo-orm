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
	return this.db.Query().Select(columns...).From(this.table)
}

func (this *Model) Raw(pred string, args ...interface{}) him.RawBuilder {
	return him.NewRawBuilder(this.db, pred, args)
}

func (this *Model) Insert() him.InsertBuilder {
	return this.db.Insert().Into(this.table.String())
}

func (this *Model) Update() him.UpdateBuilder {
	return this.db.Update().Table(this.table.String())
}

func (this *Model) Delete() him.DeleteBuilder {
	return this.db.Delete().From(this.table.String())
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
