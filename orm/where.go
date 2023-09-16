package orm

import (
	"github.com/Masterminds/squirrel"
)

type where struct {
	logic   Logic
	sqlizer squirrel.Sqlizer
}

type wheres struct {
	logic   Logic
	collect []where
}

func newWheres() *wheres {
	return &wheres{collect: make([]where, 0)}
}

func (w *wheres) forEach(fn func(w where) (bool, error)) error {
	for _, c := range w.collect {
		b, err := fn(c)
		if err != nil {
			return err
		}
		if !b {
			break
		}
	}
	return nil
}

func (w *wheres) and() *wheres {
	w.logic = AND
	return w
}

func (w *wheres) or() *wheres {
	w.logic = OR
	return w
}

func (w *wheres) raw(sql string, args []interface{}, err error) {
	if w.logic == AND {
		w.collect = append(w.collect, and("", "RAW", raw{sql, args, err}))
	} else {
		w.collect = append(w.collect, or("", "RAW", raw{sql, args, err}))
	}
}

func (w *wheres) where(column, operator string, value interface{}) {
	if w.logic == AND {
		w.collect = append(w.collect, and(column, operator, value))
	} else {
		w.collect = append(w.collect, or(column, operator, value))
	}
}

func (w *wheres) whereIn(column string, value interface{}) {
	if w.logic == AND {
		w.collect = append(w.collect, and(column, "IN", value))
	} else {
		w.collect = append(w.collect, or(column, "IN", value))
	}
}

func (w *wheres) whereNotIn(column string, value interface{}) {
	if w.logic == AND {
		w.collect = append(w.collect, and(column, "NotIn", value))
	} else {
		w.collect = append(w.collect, or(column, "NotIn", value))
	}
}

func (w *wheres) whereNull(column string) {
	if w.logic == AND {
		w.collect = append(w.collect, and(column, "Null", nil))
	} else {
		w.collect = append(w.collect, or(column, "Null", nil))
	}
}

func (w *wheres) whereNotNull(column string) {
	if w.logic == AND {
		w.collect = append(w.collect, and(column, "NotNull", nil))
	} else {
		w.collect = append(w.collect, or(column, "NotNull", nil))
	}
}

func (w *wheres) whereLike(column string, value interface{}) {
	if w.logic == AND {
		w.collect = append(w.collect, and(column, "Like", value))
	} else {
		w.collect = append(w.collect, or(column, "Like", value))
	}
}

func (w *wheres) whereNotLike(column string, value interface{}) {
	if w.logic == AND {
		w.collect = append(w.collect, and(column, "NotLike", value))
	} else {
		w.collect = append(w.collect, or(column, "NotLike", value))
	}
}

func (w *wheres) whereBetween(column string, first, second interface{}) {
	if w.logic == AND {
		w.collect = append(w.collect, and(column, "BETWEEN", between{column, first, second}))
	} else {
		w.collect = append(w.collect, or(column, "BETWEEN", between{column, first, second}))
	}
}

func (this SelectBuilder) WhereRaw(fn func(builder WhereRawBuilder) squirrel.Sqlizer) SelectBuilder {
	sql, args, err := fn(WhereRawBuilder{}).ToSql()
	this.wheres.and().raw(sql, args, err)
	return this
}

func (this SelectBuilder) OrWhereRaw(fn func(builder WhereRawBuilder) squirrel.Sqlizer) SelectBuilder {
	sql, args, err := fn(WhereRawBuilder{}).ToSql()
	this.wheres.or().raw(sql, args, err)
	return this
}

func (this SelectBuilder) Where(column, operator string, value interface{}) SelectBuilder {
	this.wheres.and().where(column, operator, value)
	return this
}

func (this SelectBuilder) WhereIn(column string, value interface{}) SelectBuilder {
	this.wheres.and().whereIn(column, value)
	return this
}

func (this SelectBuilder) WhereNotIn(column string, value interface{}) SelectBuilder {
	this.wheres.and().whereNotIn(column, value)
	return this
}

func (this SelectBuilder) WhereNull(column string) SelectBuilder {
	this.wheres.and().whereNull(column)
	return this
}

func (this SelectBuilder) WhereNotNull(column string) SelectBuilder {
	this.wheres.and().whereNotNull(column)
	return this
}

func (this SelectBuilder) WhereLike(column string, value interface{}) SelectBuilder {
	this.wheres.and().whereLike(column, value)
	return this
}

func (this SelectBuilder) NotLike(column string, value interface{}) SelectBuilder {
	this.wheres.and().whereNotLike(column, value)
	return this
}

func (this SelectBuilder) WhereBetween(column string, first, second interface{}) SelectBuilder {
	this.wheres.and().whereBetween(column, first, second)
	return this
}

func (this SelectBuilder) OrWhere(column, operator string, value interface{}) SelectBuilder {
	this.wheres.or().where(column, operator, value)
	return this
}

func (this SelectBuilder) OrWhereIn(column string, value interface{}) SelectBuilder {
	this.wheres.or().whereIn(column, value)
	return this
}

func (this SelectBuilder) OrWhereNotIn(column string, value interface{}) SelectBuilder {
	this.wheres.or().whereNotIn(column, value)
	return this
}

func (this SelectBuilder) OrWhereNull(column string) SelectBuilder {
	this.wheres.or().whereNull(column)
	return this
}

func (this SelectBuilder) OrWhereNotNull(column string) SelectBuilder {
	this.wheres.or().whereNotNull(column)
	return this
}

func (this SelectBuilder) OrLike(column string, value interface{}) SelectBuilder {
	this.wheres.or().whereLike(column, value)
	return this
}

func (this SelectBuilder) OrNotLike(column string, value interface{}) SelectBuilder {
	this.wheres.or().whereNotLike(column, value)
	return this
}

func (this SelectBuilder) OrWhereBetween(column string, first, second interface{}) SelectBuilder {
	this.wheres.or().whereBetween(column, first, second)
	return this
}
