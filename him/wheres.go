package him

import (
	"strings"
)

type Wheres struct {
	logic   Logic
	collect []where
}

func NewWheres() *Wheres {
	return &Wheres{logic: AND, collect: make([]where, 0)}
}

func (w *Wheres) clone(c *Wheres) {
	w.logic = c.logic
	w.collect = c.collect
}

func (w *Wheres) pred() (string, []interface{}, error) {
	pred := make([]string, 0)
	args := make([]interface{}, 0)
	err := w.forEach(func(w where) (bool, error) {
		sql, arg, err := w.sqlizer.ToSql()
		if err != nil {
			return false, err
		}
		pred, args, err = logic(w, sql, arg, pred, args)
		if err != nil {
			return false, err
		}
		return true, nil
	})
	return strings.Join(pred, " "), args, err
}

func (w *Wheres) len() int {
	return len(w.collect)
}

func (w *Wheres) forEach(fn func(w where) (bool, error)) error {
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

func (w *Wheres) And() *Wheres {
	return w.and()
}

func (w *Wheres) and() *Wheres {
	w.logic = AND
	return w
}

func (w *Wheres) Or() *Wheres {
	return w.or()
}

func (w *Wheres) or() *Wheres {
	w.logic = OR
	return w
}

func (w *Wheres) raw(sql string, args []interface{}, err error) {
	w.collect = append(w.collect, and("", "RAW", raw{sql, args, err}))
}

func (w *Wheres) whereRaw(sql string, args []interface{}, err error) {
	if w.logic == AND {
		w.collect = append(w.collect, and("", "whereRaw", whereRaw{sql, args, err}))
	} else {
		w.collect = append(w.collect, or("", "whereRaw", whereRaw{sql, args, err}))
	}
}

func (w *Wheres) Where(column, operator string, value interface{}) {
	w.where(column, operator, value)
}

func (w *Wheres) where(column, operator string, value interface{}) {
	if w.logic == AND {
		w.collect = append(w.collect, and(column, operator, value))
	} else {
		w.collect = append(w.collect, or(column, operator, value))
	}
}

func (w *Wheres) whereIn(column string, value interface{}) {
	if w.logic == AND {
		w.collect = append(w.collect, and(column, "IN", value))
	} else {
		w.collect = append(w.collect, or(column, "IN", value))
	}
}

func (w *Wheres) whereNotIn(column string, value interface{}) {
	if w.logic == AND {
		w.collect = append(w.collect, and(column, "NotIn", value))
	} else {
		w.collect = append(w.collect, or(column, "NotIn", value))
	}
}

func (w *Wheres) whereNull(column string) {
	if w.logic == AND {
		w.collect = append(w.collect, and(column, "Null", nil))
	} else {
		w.collect = append(w.collect, or(column, "Null", nil))
	}
}

func (w *Wheres) whereNotNull(column string) {
	if w.logic == AND {
		w.collect = append(w.collect, and(column, "NotNull", nil))
	} else {
		w.collect = append(w.collect, or(column, "NotNull", nil))
	}
}

func (w *Wheres) whereLike(column string, value interface{}) {
	if w.logic == AND {
		w.collect = append(w.collect, and(column, "Like", value))
	} else {
		w.collect = append(w.collect, or(column, "Like", value))
	}
}

func (w *Wheres) whereNotLike(column string, value interface{}) {
	if w.logic == AND {
		w.collect = append(w.collect, and(column, "NotLike", value))
	} else {
		w.collect = append(w.collect, or(column, "NotLike", value))
	}
}

func (w *Wheres) whereBetween(column string, first, second interface{}) {
	if w.logic == AND {
		w.collect = append(w.collect, and(column, "BETWEEN", between{column, first, second}))
	} else {
		w.collect = append(w.collect, or(column, "BETWEEN", between{column, first, second}))
	}
}

func (w *Wheres) Reset() {
	w = NewWheres()
}
