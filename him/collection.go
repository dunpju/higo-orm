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

func (this *SelectBuilder) First(dest interface{}) *gorm.DB {
	this = this.Limit(1)
	sql, args, err := this.ToSql()
	if err != nil {
		this.db.GormDB().Error = err
		return this.db.GormDB()
	}
	return this.db.GormDB().Raw(sql, args...).Scan(dest)
}

func (this *SelectBuilder) Get(dest interface{}) *gorm.DB {
	sql, args, err := this.ToSql()
	if err != nil {
		this.db.GormDB().Error = err
		return this.db.GormDB()
	}
	return this.db.GormDB().Raw(sql, args...).Scan(dest)
}

func (this *SelectBuilder) Paginate(page, perPage uint64, dest interface{}) (*gorm.DB, Paginate) {
	countStatement := this.count().Limit(1)
	countSql, args, err := countStatement.ToSql()
	if err != nil {
		this.db.GormDB().Error = err
		return this.db.GormDB(), Paginate{}
	}
	count_ := counter{}
	this.db.GormDB().Raw(countSql, args...).Scan(&count_)
	if this.db.GormDB().Error != nil {
		return this.db.GormDB(), Paginate{}
	}
	offset := (page - 1) * perPage
	sql, args, err := this.Offset(offset).Limit(perPage).ToSql()
	if err != nil {
		this.db.GormDB().Error = err
		return this.db.GormDB(), Paginate{}
	}
	this.db.GormDB().Raw(sql, args...).Scan(dest)
	if this.db.GormDB().Error != nil {
		return this.db.GormDB(), Paginate{}
	}
	return this.db.GormDB(), Paginate{Total: uint64(count_.Count_), PerPage: perPage, CurrentPage: page, Items: dest}
}

func (this *SelectBuilder) Count() (*gorm.DB, int64) {
	var count_ int64
	this.db.gormDB = this.db.GormDB().Table(this.from)
	if len(this.columns) > 0 {
		this.db.gormDB = this.db.GormDB().Select(strings.Join(this.columns, ","))
	}
	if this.wheres.len() > 0 {
		pred, args, err := this.wheres.pred()
		if err != nil {
			this.db.GormDB().Error = err
			return this.db.GormDB(), count_
		}
		this.db.gormDB = this.db.GormDB().Where(pred, args...)
	}
	if this.hasGroupBys {
		for _, by := range this.groupBys {
			this.db.gormDB = this.db.GormDB().Group(by)
		}
	}
	this.db.gormDB = this.db.GormDB().Count(&count_)
	return this.db.GormDB(), count_
}

func (this *SelectBuilder) Sum(column string) (*gorm.DB, uint64) {
	countStatement := this.sum(column)
	countSql, args, err := countStatement.ToSql()
	if err != nil {
		this.db.GormDB().Error = err
		return this.db.GormDB(), 0
	}
	count_ := counter{}
	this.db.GormDB().Raw(countSql, args...).Scan(&count_)
	if this.db.GormDB().Error != nil {
		return this.db.GormDB(), 0
	}
	return this.db.GormDB(), uint64(count_.Count_)
}
