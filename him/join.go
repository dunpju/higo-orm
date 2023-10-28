package him

import (
	"fmt"
	"github.com/Masterminds/squirrel"
)

type joinCase int

const (
	Join joinCase = iota + 1
	LeftJoin
	RightJoin
	InnerJoin
)

type join struct {
	jCase joinCase
	join  string
	rest  []interface{}
}

func (this *SelectBuilder) Join(table, first, operator, second string, rest ...interface{}) *SelectBuilder {
	this.joins = append(this.joins, join{Join, fmt.Sprintf("%s ON %s %s %s", table, first, operator, second), rest})
	return this
}

func (this *SelectBuilder) LeftJoin(table, first, operator, second string, rest ...interface{}) *SelectBuilder {
	this.joins = append(this.joins, join{LeftJoin, fmt.Sprintf("%s ON %s %s %s", table, first, operator, second), rest})
	return this
}

func (this *SelectBuilder) RightJoin(table, first, operator, second string, rest ...interface{}) *SelectBuilder {
	this.joins = append(this.joins, join{RightJoin, fmt.Sprintf("%s ON %s %s %s", table, first, operator, second), rest})
	return this
}

func (this *SelectBuilder) InnerJoin(table, first, operator, second string, rest ...interface{}) *SelectBuilder {
	this.joins = append(this.joins, join{InnerJoin, fmt.Sprintf("%s ON %s %s %s", table, first, operator, second), rest})
	return this
}

func joins(selectBuilder squirrel.SelectBuilder, joins []join) squirrel.SelectBuilder {
	for _, j := range joins {
		if j.jCase == LeftJoin {
			selectBuilder = selectBuilder.LeftJoin(j.join, j.rest...)
		} else if j.jCase == RightJoin {
			selectBuilder = selectBuilder.RightJoin(j.join, j.rest...)
		} else if j.jCase == InnerJoin {
			selectBuilder = selectBuilder.InnerJoin(j.join, j.rest...)
		} else {
			selectBuilder = selectBuilder.Join(j.join, j.rest...)
		}
	}
	return selectBuilder
}
