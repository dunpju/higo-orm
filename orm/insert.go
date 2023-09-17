package orm

import (
	"context"
	"github.com/Masterminds/squirrel"
	"gorm.io/gorm"
)

func Insert(into string) InsertBuilder {
	return insert(into)
}

func insert(into string) InsertBuilder {
	return InsertBuilder{builder: squirrel.Insert(into)}
}

type InsertBuilder struct {
	builder squirrel.InsertBuilder
}

func (this InsertBuilder) Columns(columns ...string) InsertBuilder {
	this.builder = this.builder.Columns(columns...)
	return this
}

func (this InsertBuilder) Values(values ...interface{}) InsertBuilder {
	this.builder = this.builder.Values(values...)
	return this
}

func (this InsertBuilder) ToSql() (string, []interface{}, error) {
	return this.builder.ToSql()
}

// Deprecated
func (this InsertBuilder) LastInsertId() (*gorm.DB, int64) {
	db, err := Gorm()
	if err != nil {
		db.Error = err
		return db, 0
	}
	sql, args, err := this.ToSql()
	if err != nil {
		db.Error = err
		return db, 0
	}
	ctx := context.Background()
	result, err := db.ConnPool.ExecContext(ctx, sql, args...)
	if err != nil {
		db.Error = err
		return db, 0
	}
	db.Exec(sql, args...)
	if db.Error != nil {
		return db, 0
	}
	id, err := result.LastInsertId()
	if err != nil {
		db.Error = err
		return db, 0
	}
	return db, id
}
