package orm

import (
	"github.com/Masterminds/squirrel"
)

type SelectBuilder struct {
	columns     []string
	from        string
	joins       []join
	wheres      []squirrel.Sqlizer
	hasOffset   bool
	offset      uint64
	hasLimit    bool
	limit       uint64
	hasOrderBy  bool
	orderBy     []string
	hasGroupBys bool
	groupBys    []string
	hasHaving   bool
	having      having
}

func Query() SelectBuilder {
	return SelectBuilder{
		columns: make([]string, 0),
		joins:   make([]join, 0),
		wheres:  make([]squirrel.Sqlizer, 0),
		orderBy: make([]string, 0),
	}
}

func (this SelectBuilder) Select(columns ...string) SelectBuilder {
	this.columns = make([]string, 0)
	this.columns = append(this.columns, columns...)
	return this
}

func (this SelectBuilder) From(from string) SelectBuilder {
	this.from = from
	return this
}

func (this SelectBuilder) Offset(limit uint64) SelectBuilder {
	this.hasOffset = true
	this.limit = limit
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
	this.groupBys = groupBys
	return this
}

type having struct {
	pred interface{}
	rest []interface{}
}

func (this SelectBuilder) Having(pred interface{}, rest ...interface{}) SelectBuilder {
	this.hasHaving = true
	this.having = having{pred: pred, rest: rest}
	return this
}

func (this SelectBuilder) ToSql() (string, []interface{}, error) {
	selectBuilder := squirrel.Select(this.columns...)
	selectBuilder = selectBuilder.From(this.from)
	selectBuilder = joins(selectBuilder, this.joins)
	selectBuilder = wheres(selectBuilder, this.wheres)
	if this.hasOrderBy {
		selectBuilder = selectBuilder.OrderBy(this.orderBy...)
	}
	if this.hasGroupBys {
		selectBuilder = selectBuilder.GroupBy(this.groupBys...)
	}
	if this.hasHaving {
		selectBuilder = selectBuilder.Having(this.having.pred, this.having.rest)
	}
	if this.hasOffset {
		selectBuilder = selectBuilder.Offset(this.offset)
	}
	if this.hasLimit {
		selectBuilder = selectBuilder.Limit(this.limit)
	}
	return selectBuilder.ToSql()
}

func joins(selectBuilder squirrel.SelectBuilder, joins []join) squirrel.SelectBuilder {
	for _, j := range joins {
		if j.jCase == LeftJoin {
			selectBuilder = selectBuilder.LeftJoin(j.join, j.rest)
		} else if j.jCase == RightJoin {
			selectBuilder = selectBuilder.RightJoin(j.join, j.rest)
		} else if j.jCase == InnerJoin {
			selectBuilder = selectBuilder.InnerJoin(j.join, j.rest)
		} else {
			selectBuilder = selectBuilder.Join(j.join, j.rest)
		}
	}
	return selectBuilder
}

func wheres(selectBuilder squirrel.SelectBuilder, wheres []squirrel.Sqlizer) squirrel.SelectBuilder {
	for _, sqlizer := range wheres {
		selectBuilder = selectBuilder.Where(sqlizer)
	}
	return selectBuilder
}
