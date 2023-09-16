package orm

type wheres struct {
	logic   Logic
	collect []where
}

func newWheres() *wheres {
	return &wheres{logic: AND, collect: make([]where, 0)}
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
