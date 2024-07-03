package him

import (
	"github.com/dunpju/higo-orm/event"
	"gorm.io/gorm"
	"math"
	"strings"
)

const (
	_count_ = "count_"
	_sum_   = "sum_"
)

type counter struct {
	Count_ int64
}

type sum struct {
	Sum_ string
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
	var (
		ok          bool
		paginate    IPaginate
		paginateSum IPaginateSum
	)
	if paginateSum, ok = dest.(IPaginateSum); ok {
		if paginate, ok = paginateSum.Dest().(IPaginate); !ok {
			paginate = NewPaginate(WithItems(paginateSum.Dest()))
		}
	} else if paginate, ok = dest.(IPaginate); !ok {
		paginate = NewPaginate(WithItems(dest))
	}

	paginate.SetPerPage(perPage).SetCurrentPage(page)

	if paginateSum != nil {
		sumSql, args, err := this.clone().sum(toStrings(paginateSum.Field()...)...).Limit(1).ToSql()
		this.eventBeforeSum(sumSql, args, err, nil)
		if err != nil {
			this.db.GormDB().Error = err
			return this.db.GormDB(), paginate
		}
		var sum_ interface{}
		if len(paginateSum.Field()) > 1 {
			sum_ = make(map[string]interface{})
		} else {
			sum_ = &sum{}
		}
		this.db.GormDB().Raw(sumSql, args...).Scan(sum_)

		this.eventAfterSum(sumSql, args, this.db.GormDB().Error, sum_)
		if s, ok := sum_.(*sum); ok {
			paginateSum.SetSum(s.Sum_)
		} else {
			paginateSum.SetSum(sum_)
		}
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

	var offset uint64
	if page > 0 {
		offset = (page - 1) * perPage
	}
	sql, args, err := this.Offset(offset).Limit(perPage).ToSql()

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

func (this *SelectBuilder) Sum(column string, more ...string) (*gorm.DB, interface{}) {
	columns := make([]string, 0)
	columns = append(columns, column)
	columns = append(columns, more...)
	sumStatement := this.sum(columns...)
	sumSql, args, err := sumStatement.Limit(1).ToSql()
	if err != nil {
		this.db.GormDB().Error = err
		return this.db.GormDB(), nil
	}
	var sum_ interface{}
	if len(this.sumColumn) > 1 {
		sum_ = make(map[string]interface{})
	} else {
		sum_ = &sum{}
	}
	this.db.GormDB().Raw(sumSql, args...).Scan(sum_)
	if this.db.GormDB().Error != nil {
		return this.db.GormDB(), nil
	}
	if s, ok := sum_.(*sum); ok {
		return this.db.GormDB(), s.Sum_
	}
	return this.db.GormDB(), sum_
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

func (this *SelectBuilder) eventBeforeSum(sql string, args []interface{}, err error, result interface{}) {
	event.Point(event.BeforeSum, event.NewEventRecordResult(this.from, sql, args, err, result))
}

func (this *SelectBuilder) eventAfterSum(sql string, args []interface{}, err error, result interface{}) {
	event.Point(event.AfterSum, event.NewEventRecordResult(this.from, sql, args, err, result))
}
