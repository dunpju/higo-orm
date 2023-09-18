package him

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
	sql, args, err := this.ToSql()
	if err != nil {
		this.db.Error = err
		return this.db
	}
	return this.db.Raw(sql, args...).Scan(dest)
}

func (this SelectBuilder) Get(dest interface{}) *gorm.DB {
	sql, args, err := this.ToSql()
	if err != nil {
		this.db.Error = err
		return this.db
	}
	return this.db.Raw(sql, args...).Scan(dest)
}

func (this SelectBuilder) Paginate(page, perPage uint64, dest interface{}) (*gorm.DB, Paginate) {
	countStatement := this.count().Limit(1)
	countSql, args, err := countStatement.ToSql()
	if err != nil {
		this.db.Error = err
		return this.db, Paginate{}
	}
	count_ := counter{}
	this.db.Raw(countSql, args...).Scan(&count_)
	if this.db.Error != nil {
		return this.db, Paginate{}
	}
	offset := (page - 1) * perPage
	sql, args, err := this.Offset(offset).Limit(perPage).ToSql()
	if err != nil {
		this.db.Error = err
		return this.db, Paginate{}
	}
	this.db.Raw(sql, args...).Scan(dest)
	if this.db.Error != nil {
		return this.db, Paginate{}
	}
	return this.db, Paginate{Total: uint64(count_.Count_), PerPage: perPage, CurrentPage: page, Items: dest}
}

func (this SelectBuilder) Count() int64 {
	var count_ int64
	this.db = this.db.Table(this.from)
	if len(this.columns) > 0 {
		this.db = this.db.Select(strings.Join(this.columns, ","))
	}
	if this.wheres.len() > 0 {
		pred, args, err := this.wheres.pred()
		if err != nil {
			this.db.Error = err
			return count_
		}
		this.db = this.db.Where(pred, args...)
	}
	if this.hasGroupBys {
		for _, by := range this.groupBys {
			this.db = this.db.Group(by)
		}
	}
	this.db = this.db.Count(&count_)
	return count_
}

func (this SelectBuilder) Sum(column string) uint64 {
	countStatement := this.sum(column)
	countSql, args, err := countStatement.ToSql()
	if err != nil {
		this.db.Error = err
		return 0
	}
	count_ := counter{}
	this.db.Raw(countSql, args...).Scan(&count_)
	if this.db.Error != nil {
		return 0
	}
	return uint64(count_.Count_)
}
