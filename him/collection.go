package him

import (
	"github.com/dunpju/higo-orm/event"
	"gorm.io/gorm"
	"math"
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

	this.eventBefore(sql, args, err, nil)

	if err != nil {
		this.db.GormDB().Error = err
		return this.db.GormDB()
	}

	this.db.GormDB().Raw(sql, args...).Scan(dest)

	this.eventAfter(sql, args, err, dest)

	return this.db.GormDB()
}

func (this *SelectBuilder) Get(dest interface{}) *gorm.DB {
	sql, args, err := this.ToSql()

	this.eventBefore(sql, args, err, nil)

	if err != nil {
		this.db.GormDB().Error = err
		return this.db.GormDB()
	}

	this.db.GormDB().Raw(sql, args...).Scan(dest)

	this.eventAfter(sql, args, err, dest)

	return this.db.GormDB()
}

func (this *SelectBuilder) Paginate(page, perPage uint64, dest interface{}) (*gorm.DB, IPaginate) {
	paginate, ok := dest.(IPaginate)
	if !ok {
		paginate = NewPaginate(WithItems(dest))
	}
	countSql, args, err := this.clone().count().Limit(1).ToSql()

	this.eventBeforeCount(countSql, args, err, nil)

	if err != nil {
		this.db.GormDB().Error = err
		return this.db.GormDB(), paginate
	}
	count_ := counter{}
	this.db.GormDB().Raw(countSql, args...).Scan(&count_)

	this.eventAfterCount(countSql, args, this.db.GormDB().Error, count_)

	if this.db.GormDB().Error != nil {
		return this.db.GormDB(), paginate
	}
	if count_.Count_ == 0 {
		return this.db.GormDB(), paginate
	}
	sql, args, err := this.Offset((page - 1) * perPage).Limit(perPage).ToSql()

	this.eventBefore(countSql, args, err, nil)

	if err != nil {
		this.db.GormDB().Error = err
		return this.db.GormDB(), paginate
	}
	this.db.GormDB().Raw(sql, args...).Scan(paginate.GetItems())
	if this.db.GormDB().Error != nil {
		return this.db.GormDB(), paginate
	}

	paginate.SetTotal(uint64(count_.Count_)).
		SetPerPage(perPage).
		SetCurrentPage(page).
		SetLastPage(uint64(math.Ceil(float64(count_.Count_) / float64(perPage))))

	this.eventAfter(sql, args, this.db.GormDB().Error, paginate)

	return this.db.GormDB(), paginate
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
	sumStatement := this.sum(column)
	sumSql, args, err := sumStatement.ToSql()
	if err != nil {
		this.db.GormDB().Error = err
		return this.db.GormDB(), 0
	}
	count_ := counter{}
	this.db.GormDB().Raw(sumSql, args...).Scan(&count_)
	if this.db.GormDB().Error != nil {
		return this.db.GormDB(), 0
	}
	return this.db.GormDB(), uint64(count_.Count_)
}

func (this *SelectBuilder) eventBefore(sql string, args []interface{}, err error, result interface{}) {
	event.Point(event.BeforeSelect, event.NewEventRecordResult(this.from, sql, args, err, result))
}

func (this *SelectBuilder) eventAfter(sql string, args []interface{}, err error, result interface{}) {
	event.Point(event.AfterSelect, event.NewEventRecordResult(this.from, sql, args, err, result))
}

func (this *SelectBuilder) eventBeforeCount(sql string, args []interface{}, err error, result interface{}) {
	event.Point(event.BeforeCount, event.NewEventRecordResult(this.from, sql, args, err, result))
}

func (this *SelectBuilder) eventAfterCount(sql string, args []interface{}, err error, result interface{}) {
	event.Point(event.AfterCount, event.NewEventRecordResult(this.from, sql, args, err, result))
}
