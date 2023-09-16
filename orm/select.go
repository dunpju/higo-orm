package orm

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"strings"
)

type SelectBuilder struct {
	isWhereRaw  bool
	columns     []string
	from        string
	joins       []join
	wheres      []where
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
	return query()
}

func query() SelectBuilder {
	return SelectBuilder{
		columns:  make([]string, 0),
		joins:    make([]join, 0),
		wheres:   make([]where, 0),
		orderBy:  make([]string, 0),
		groupBys: make([]string, 0),
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
	if this.isWhereRaw {
		return whereRaw(this.wheres)
	}
	selectBuilder := squirrel.Select(this.columns...)
	if this.from != "" {
		selectBuilder = selectBuilder.From(this.from)
	}
	selectBuilder = joins(selectBuilder, this.joins)
	selectBuilder, err := wheres(selectBuilder, this.wheres)
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

func wheres(selectBuilder squirrel.SelectBuilder, wheres []where) (squirrel.SelectBuilder, error) {
	pred := make([]string, 0)
	args := make([]interface{}, 0)
	for _, w := range wheres {
		sql, arg, err := w.sqlizer.ToSql()
		if err != nil {
			return selectBuilder, err
		}
		pred, args, err = logic(w, sql, arg, pred, args)
		if err != nil {
			return selectBuilder, err
		}
	}
	selectBuilder = selectBuilder.Where(strings.Join(pred, " "), args...)
	return selectBuilder, nil
}

func whereRaw(wheres []where) (string, []interface{}, error) {
	pred := make([]string, 0)
	args := make([]interface{}, 0)
	for _, w := range wheres {
		sql, arg, err := w.sqlizer.ToSql()
		if err != nil {
			return "", nil, err
		}
		pred, args, err = logic(w, sql, arg, pred, args)
		if err != nil {
			return "", nil, err
		}
	}
	return strings.Join(pred, " "), args, nil
}

func logic(w where, sql string, arg []interface{}, pred []string, args []interface{}) ([]string, []interface{}, error) {
	if w.logic == "AND" {
		if len(pred) == 0 {
			pred = append(pred, fmt.Sprintf("(%s)", sql))
			args = append(args, arg...)
		} else {
			pred = append(pred, "AND")
			pred = append(pred, fmt.Sprintf("(%s)", sql))
			args = append(args, arg...)
		}
	} else {
		if len(pred) == 0 {
			pred = append(pred, fmt.Sprintf("(%s)", sql))
			args = append(args, arg...)
		} else {
			pred = append(pred, "OR")
			pred = append(pred, fmt.Sprintf("(%s)", sql))
			args = append(args, arg...)
		}
	}
	return pred, args, nil
}
