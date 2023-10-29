package him

import (
	"github.com/Masterminds/squirrel"
	"gorm.io/gorm"
)

type DeleteBuilder struct {
	db      *gorm.DB
	connect *connect
	builder squirrel.DeleteBuilder
	wheres  *Wheres
	Error   error
}

func newDeleteBuilder(connect string) *DeleteBuilder {
	if connect != "" {
		dbc, err := getConnect(connect)
		if err != nil {
			return &DeleteBuilder{Error: err}
		}
		return &DeleteBuilder{db: dbc.db.GormDB(), connect: dbc, wheres: NewWheres()}
	} else {
		dbc, err := getConnect(DefaultConnect)
		if err != nil {
			return &DeleteBuilder{Error: err}
		}
		return &DeleteBuilder{db: dbc.db.GormDB(), connect: dbc, wheres: NewWheres()}
	}
}

func (this *DeleteBuilder) delete(from string) *DeleteBuilder {
	this.builder = squirrel.Delete(from)
	return this
}

func (this *DeleteBuilder) DB() *gorm.DB {
	return this.db
}

func (this *DeleteBuilder) TX(tx *gorm.DB) *DeleteBuilder {
	this = this.begin(tx)
	return this
}

func (this *DeleteBuilder) begin(db *gorm.DB) *DeleteBuilder {
	this.db = db
	return this
}

func (this *DeleteBuilder) WhereRaw(fn func(builder WhereRawBuilder) WhereRawBuilder) *DeleteBuilder {
	sql, args, err := fn(WhereRawBuilder{}).ToSql()
	this.wheres.and().whereRaw(sql, args, err)
	return this
}

func (this *DeleteBuilder) OrWhereRaw(fn func(builder WhereRawBuilder) WhereRawBuilder) *DeleteBuilder {
	sql, args, err := fn(WhereRawBuilder{}).ToSql()
	this.wheres.or().whereRaw(sql, args, err)
	return this
}

func (this *DeleteBuilder) Where(column any, operator string, value interface{}) *DeleteBuilder {
	this.wheres.and().where(columnToString(column), operator, value)
	return this
}

func (this *DeleteBuilder) WhereIn(column any, value interface{}) *DeleteBuilder {
	this.wheres.and().whereIn(columnToString(column), value)
	return this
}

func (this *DeleteBuilder) WhereNotIn(column any, value interface{}) *DeleteBuilder {
	this.wheres.and().whereNotIn(columnToString(column), value)
	return this
}

func (this *DeleteBuilder) WhereNull(column any) *DeleteBuilder {
	this.wheres.and().whereNull(columnToString(column))
	return this
}

func (this *DeleteBuilder) WhereNotNull(column any) *DeleteBuilder {
	this.wheres.and().whereNotNull(columnToString(column))
	return this
}

func (this *DeleteBuilder) WhereLike(column any, value interface{}) *DeleteBuilder {
	this.wheres.and().whereLike(columnToString(column), value)
	return this
}

func (this *DeleteBuilder) NotLike(column any, value interface{}) *DeleteBuilder {
	this.wheres.and().whereNotLike(columnToString(column), value)
	return this
}

func (this *DeleteBuilder) WhereBetween(column any, first, second interface{}) *DeleteBuilder {
	this.wheres.and().whereBetween(columnToString(column), first, second)
	return this
}

func (this *DeleteBuilder) OrWhere(column any, operator string, value interface{}) *DeleteBuilder {
	this.wheres.or().where(columnToString(column), operator, value)
	return this
}

func (this *DeleteBuilder) OrWhereIn(column any, value interface{}) *DeleteBuilder {
	this.wheres.or().whereIn(columnToString(column), value)
	return this
}

func (this *DeleteBuilder) OrWhereNotIn(column any, value interface{}) *DeleteBuilder {
	this.wheres.or().whereNotIn(columnToString(column), value)
	return this
}

func (this *DeleteBuilder) OrWhereNull(column any) *DeleteBuilder {
	this.wheres.or().whereNull(columnToString(column))
	return this
}

func (this *DeleteBuilder) OrWhereNotNull(column any) *DeleteBuilder {
	this.wheres.or().whereNotNull(columnToString(column))
	return this
}

func (this *DeleteBuilder) OrLike(column any, value interface{}) *DeleteBuilder {
	this.wheres.or().whereLike(columnToString(column), value)
	return this
}

func (this *DeleteBuilder) OrNotLike(column any, value interface{}) *DeleteBuilder {
	this.wheres.or().whereNotLike(columnToString(column), value)
	return this
}

func (this *DeleteBuilder) OrWhereBetween(column any, first, second interface{}) *DeleteBuilder {
	this.wheres.or().whereBetween(columnToString(column), first, second)
	return this
}

func (this *DeleteBuilder) whereHandle() (*DeleteBuilder, error) {
	pred, args, err := this.wheres.pred()
	if err != nil {
		return this, err
	}
	this.builder = this.builder.Where(pred, args...)
	return this, nil
}

func (this *DeleteBuilder) ToSql() (string, []interface{}, error) {
	builder, err := this.whereHandle()
	if err != nil {
		return "", nil, err
	}
	this = builder
	return this.builder.ToSql()
}

func (this *DeleteBuilder) Exec() (*gorm.DB, int64) {
	gormDB, _, rowsAffected := newExecer(this, this.db).exec()
	this.db = gormDB
	return this.db, rowsAffected
}

func (this *DeleteBuilder) SetWheres(wheres *Wheres) *DeleteBuilder {
	this.wheres = wheres
	return this
}
