package orm

import (
	"fmt"
	"github.com/Masterminds/squirrel"
)

type between struct {
	column        string
	first, second interface{}
}

func (this between) ToSql() (string, []interface{}, error) {
	args := make([]interface{}, 0)
	args = append(args, this.first, this.second)
	return fmt.Sprintf("%s BETWEEN ? ?", this.column), args, nil
}

func where(column, operator string, value interface{}) squirrel.Sqlizer {
	if operator == ">" {
		return squirrel.Gt{column: value}
	} else if operator == ">=" {
		return squirrel.GtOrEq{column: value}
	} else if operator == "<" {
		return squirrel.Lt{column: value}
	} else if operator == "<=" {
		return squirrel.LtOrEq{column: value}
	} else if operator == "<>" || operator == "!=" {
		return squirrel.NotEq{column: value}
	} else {
		return squirrel.Eq{column: value}
	}
}

func (this SelectBuilder) Where(column, operator string, value interface{}) SelectBuilder {
	this.wheres = append(this.wheres, where(column, operator, value))
	return this
}

func (this SelectBuilder) OrWhere(column, operator string, value interface{}) SelectBuilder {
	this.wheres = append(this.wheres, squirrel.Or{where(column, operator, value)})
	return this
}

func (this SelectBuilder) WhereIn(column string, value interface{}) SelectBuilder {
	this.wheres = append(this.wheres, squirrel.Eq{column: value})
	return this
}

func (this SelectBuilder) WhereNotIn(column string, value interface{}) SelectBuilder {
	this.wheres = append(this.wheres, squirrel.NotEq{column: value})
	return this
}

func (this SelectBuilder) OrWhereIn(column string, value interface{}) SelectBuilder {
	this.wheres = append(this.wheres, squirrel.Or{squirrel.Eq{column: value}})
	return this
}

func (this SelectBuilder) OrWhereNotIn(column string, value interface{}) SelectBuilder {
	this.wheres = append(this.wheres, squirrel.Or{squirrel.NotEq{column: value}})
	return this
}

func (this SelectBuilder) Like(column string, value interface{}) SelectBuilder {
	this.wheres = append(this.wheres, squirrel.Like{column: value})
	return this
}

func (this SelectBuilder) NotLike(column string, value interface{}) SelectBuilder {
	this.wheres = append(this.wheres, squirrel.NotLike{column: value})
	return this
}

func (this SelectBuilder) OrLike(column string, value interface{}) SelectBuilder {
	this.wheres = append(this.wheres, squirrel.Or{squirrel.Like{column: value}})
	return this
}

func (this SelectBuilder) OrNotLike(column string, value interface{}) SelectBuilder {
	this.wheres = append(this.wheres, squirrel.Or{squirrel.NotLike{column: value}})
	return this
}

func (this SelectBuilder) Between(column string, value1, value2 interface{}) SelectBuilder {
	this.wheres = append(this.wheres, between{column, value1, value2})
	return this
}
