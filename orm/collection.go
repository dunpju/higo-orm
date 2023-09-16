package orm

import (
	"gorm.io/gorm"
	"strings"
)

const (
	_count_ = "count_"
)

type counter struct {
	Count_ int64
}

func (this SelectBuilder) First(dest interface{}) *gorm.DB {
	this = this.Limit(1)
	db, err := Gorm()
	if err != nil {
		db.Error = err
		return db
	}
	sql, args, err := this.ToSql()
	if err != nil {
		db.Error = err
		return db
	}
	return db.Raw(sql, args...).Scan(dest)
}

func (this SelectBuilder) Get(dest interface{}) *gorm.DB {
	db, err := Gorm()
	if err != nil {
		db.Error = err
		return db
	}
	sql, args, err := this.ToSql()
	if err != nil {
		db.Error = err
		return db
	}
	return db.Raw(sql, args...).Scan(dest)
}

func (this SelectBuilder) Paginate(page, perPage uint64, dest interface{}) (*gorm.DB, Paginate) {
	db, err := Gorm()
	if err != nil {
		db.Error = err
		return db, Paginate{}
	}
	countStatement := this.count().Limit(1)
	countSql, args, err := countStatement.ToSql()
	if err != nil {
		db.Error = err
		return db, Paginate{}
	}
	count_ := counter{}
	db.Raw(countSql, args...).Scan(&count_)
	if db.Error != nil {
		return db, Paginate{}
	}
	offset := (page - 1) * perPage
	sql, args, err := this.Offset(offset).Limit(perPage).ToSql()
	if err != nil {
		db.Error = err
		return db, Paginate{}
	}
	db.Raw(sql, args...).Scan(dest)
	if db.Error != nil {
		return db, Paginate{}
	}
	return db, Paginate{Total: uint64(count_.Count_), PerPage: perPage, CurrentPage: page, Items: dest}
}

func (this SelectBuilder) Count() (*gorm.DB, int64) {
	db, err := Gorm()
	if err != nil {
		db.Error = err
		return db, 0
	}
	var count_ int64
	db = db.Table(this.from)
	if len(this.columns) > 0 {
		db = db.Select(strings.Join(this.columns, ","))
	}
	if this.wheres.len() > 0 {
		pred, args, err := this.wheres.pred()
		if err != nil {
			db.Error = err
			return db, count_
		}
		db = db.Where(pred, args...)
	}
	db = db.Count(&count_)
	return db, count_
}

func (this SelectBuilder) Sum(column string) (*gorm.DB, uint64) {
	db, err := Gorm()
	if err != nil {
		db.Error = err
		return db, 0
	}
	countStatement := this.sum(column)
	countSql, args, err := countStatement.ToSql()
	if err != nil {
		db.Error = err
		return db, 0
	}
	count_ := counter{}
	db.Raw(countSql, args...).Scan(&count_)
	if db.Error != nil {
		return db, 0
	}
	return db, uint64(count_.Count_)
}
