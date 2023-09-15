package orm

import "fmt"

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

func (this SelectBuilder) Join(table, first, operator, second string, rest ...interface{}) SelectBuilder {
	this.joins = append(this.joins, join{Join, fmt.Sprintf("%s ON %s %s %s", table, first, operator, second), rest})
	return this
}

func (this SelectBuilder) LeftJoin(table, first, operator, second string, rest ...interface{}) SelectBuilder {
	this.joins = append(this.joins, join{LeftJoin, fmt.Sprintf("%s ON %s %s %s", table, first, operator, second), rest})
	return this
}

func (this SelectBuilder) RightJoin(table, first, operator, second string, rest ...interface{}) SelectBuilder {
	this.joins = append(this.joins, join{RightJoin, fmt.Sprintf("%s ON %s %s %s", table, first, operator, second), rest})
	return this
}

func (this SelectBuilder) InnerJoin(table, first, operator, second string, rest ...interface{}) SelectBuilder {
	this.joins = append(this.joins, join{InnerJoin, fmt.Sprintf("%s ON %s %s %s", table, first, operator, second), rest})
	return this
}
