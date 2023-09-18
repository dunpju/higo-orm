package him

import (
	"fmt"
	"github.com/Masterminds/squirrel"
)

type SelectBuilder struct {
	isCount     bool
	countColumn string
	isSum       bool
	sumColumn   string
	isRaw       bool
	isWhereRaw  bool
	columns     []string
	from        string
	joins       []join
	wheres      *wheres
	hasOffset   bool
	offset      uint64
	hasLimit    bool
	limit       uint64
	hasOrderBy  bool
	orderBy     []string
	hasGroupBys bool
	groupBys    []string
	hasHaving   bool
	havings     []having
	hasColumn   bool
	column      []column
	hasDistinct bool
}

func Query() SelectBuilder {
	return query()
}

func query() SelectBuilder {
	return SelectBuilder{
		columns:  make([]string, 0),
		joins:    make([]join, 0),
		wheres:   newWheres(),
		orderBy:  make([]string, 0),
		groupBys: make([]string, 0),
		havings:  make([]having, 0),
		column:   make([]column, 0),
	}
}

func (this SelectBuilder) Select(columns ...string) SelectBuilder {
	this.columns = append(this.columns, columns...)
	return this
}

func (this SelectBuilder) From(from string) SelectBuilder {
	this.from = from
	return this
}

func (this SelectBuilder) Offset(offset uint64) SelectBuilder {
	this.hasOffset = true
	this.offset = offset
	return this
}

func (this SelectBuilder) Limit(limit uint64) SelectBuilder {
	this.hasLimit = true
	this.limit = limit
	return this
}

func (this SelectBuilder) OrderBy(orderBys ...string) SelectBuilder {
	this.hasOrderBy = true
	this.orderBy = orderBys
	return this
}

func (this SelectBuilder) GroupBy(groupBys ...string) SelectBuilder {
	this.hasGroupBys = true
	this.groupBys = append(this.groupBys, groupBys...)
	return this
}

func (this SelectBuilder) Column(col interface{}, args ...interface{}) SelectBuilder {
	this.hasColumn = true
	this.column = append(this.column, column{column: col, args: args})
	return this
}

func (this SelectBuilder) Distinct() SelectBuilder {
	this.hasDistinct = true
	return this
}

func (this SelectBuilder) count() SelectBuilder {
	this.isCount = true
	this.countColumn = fmt.Sprintf("COUNT(*) %s", _count_)
	return this
}

func (this SelectBuilder) sum(column string) SelectBuilder {
	this.isSum = true
	this.sumColumn = fmt.Sprintf("SUM(%s) %s", column, _count_)
	return this
}

func (this SelectBuilder) ToSql() (string, []interface{}, error) {
	if this.isWhereRaw || this.isRaw {
		return whereRawHandle(*this.wheres)
	}
	if this.isCount {
		this.columns = make([]string, 0)
		this.columns = append(this.columns, this.countColumn)
	}
	if this.isSum {
		this.columns = make([]string, 0)
		this.columns = append(this.columns, this.sumColumn)
	}
	selectBuilder := squirrel.Select(this.columns...)
	selectBuilder = selectBuilder.From(this.from)
	selectBuilder = joins(selectBuilder, this.joins)
	selectBuilder, err := this.whereHandle(selectBuilder, this.wheres)
	if err != nil {
		return "", nil, err
	}
	if this.hasOrderBy {
		selectBuilder = selectBuilder.OrderBy(this.orderBy...)
	}
	if this.hasGroupBys {
		selectBuilder = selectBuilder.GroupBy(this.groupBys...)
	}
	if this.hasHaving {
		for _, h := range this.havings {
			selectBuilder = selectBuilder.Having(h.pred, h.rest...)
		}
	}
	if this.hasOffset {
		selectBuilder = selectBuilder.Offset(this.offset)
	}
	if this.hasLimit {
		selectBuilder = selectBuilder.Limit(this.limit)
	}
	if this.hasColumn {
		for _, c := range this.column {
			selectBuilder = selectBuilder.Column(c.column, c.args...)
		}
	}
	if this.hasDistinct {
		selectBuilder = selectBuilder.Distinct()
	}
	return selectBuilder.ToSql()
}
