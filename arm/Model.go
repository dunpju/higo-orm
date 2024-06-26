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
	wheres  *him.Wheres
	sets    *him.Sets
	begin   bool
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
	return &Model{db: db, model: model, table: model.TableName(), wheres: him.NewWheres(), sets: him.NewSets()}
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

func (this *Model) ForceIndex(index string, more ...string) *Model {
	this.table.ForceIndex(index, more...)
	return this
}

func (this *Model) Select(columns ...any) *him.SelectBuilder {
	return this.db.Query().Select(columns...).From(this.table.String())
}

func (this *Model) Raw(pred string, args ...interface{}) him.RawBuilder {
	return him.NewRawBuilder(this.db, pred, args, this.table.String())
}

func (this *Model) Insert() *him.InsertBuilder {
	builder := this.db.Insert().Into(this.table.String())
	this.sets.ForEach(func(s him.Set) bool {
		builder.Set(s.Column(), s.Value())
		return true
	}).Reset()
	this.builder = builder
	return this.builder.(*him.InsertBuilder)
}

func (this *Model) Update() *him.UpdateBuilder {
	builder := this.db.Update().Table(this.table.String()).SetWheres(this.wheres)
	this.sets.ForEach(func(s him.Set) bool {
		builder.Set(s.Column(), s.Value())
		return true
	}).Reset()
	this.wheres.Reset()
	this.builder = builder
	return this.builder.(*him.UpdateBuilder)
}

func (this *Model) Delete() *him.DeleteBuilder {
	this.builder = this.db.Delete().From(this.table.String())
	return this.builder.(*him.DeleteBuilder)
}

func (this *Model) beginTX(opts ...*sql.TxOptions) *gorm.DB {
	db, err := him.DBConnect(this.db.Connect())
	if err != nil {
		this.err = err
		return nil
	}
	return db.GormDB().Begin(opts...)
}

func (this *Model) Begin(opts ...*sql.TxOptions) *him.TX {
	this.begin = true
	return him.NewTX(this.beginTX(opts...))
}

func (this *Model) TX(tx *gorm.DB) *Model {
	this.begin = true
	this.db.TX(tx)
	return this
}

func (this *Model) Error() error {
	return this.err
}

func (this *Model) GormDB() *gorm.DB {
	return this.db.GormDB()
}

func (this *Model) Builder(dao IDao, fn func()) IDao {
	model := this.model.Mutate()
	if this.begin {
		model.TX(this.db.GormDB())
	}
	dao.SetModel(model)
	fn()
	return dao
}

func (this *Model) Set(column any, value interface{}) *Model {
	this.sets.Append(column, value)
	if insertBuilder, insertOk := this.builder.(*him.InsertBuilder); insertOk {
		this.builder = insertBuilder.Column(column, value)
	} else if updateBuilder, updateOk := this.builder.(*him.UpdateBuilder); updateOk {
		this.builder = updateBuilder.Set(column, value)
	}
	return this
}

func (this *Model) CaseWhen(column him.CaseWhen) *Model {
	this.sets.Append(column.Field(), column.Builder())
	this.builder = this.builder.(*him.UpdateBuilder).CaseWhen(column)
	return this
}

func (this *Model) Where(column any, operator string, value interface{}) *Model {
	this.wheres.And().Where(him.ColumnToString(column), operator, value)
	return this
}

func (this *Model) IsEmpty(m IModel) bool {
	return IsEmpty(m)
}
