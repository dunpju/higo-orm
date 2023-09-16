package orm

func whereRawHandle(wheres wheres) (string, []interface{}, error) {
	pred, args, err := wheres.pred()
	if err != nil {
		return "", nil, err
	}
	return pred, args, nil
}

type WhereRawBuilder struct {
	wheres wheres
}

func (this WhereRawBuilder) ToSql() (string, []interface{}, error) {
	return whereRawHandle(this.wheres)
}

func (this WhereRawBuilder) Where(column, operator string, value interface{}) WhereRawBuilder {
	this.wheres.and().where(column, operator, value)
	return this
}

func (this WhereRawBuilder) WhereIn(column string, value interface{}) WhereRawBuilder {
	this.wheres.and().whereIn(column, value)
	return this
}

func (this WhereRawBuilder) WhereNotIn(column string, value interface{}) WhereRawBuilder {
	this.wheres.and().whereNotIn(column, value)
	return this
}

func (this WhereRawBuilder) WhereNull(column string) WhereRawBuilder {
	this.wheres.and().whereNull(column)
	return this
}

func (this WhereRawBuilder) WhereNotNull(column string) WhereRawBuilder {
	this.wheres.and().whereNotNull(column)
	return this
}

func (this WhereRawBuilder) WhereLike(column string, value interface{}) WhereRawBuilder {
	this.wheres.and().whereLike(column, value)
	return this
}

func (this WhereRawBuilder) NotLike(column string, value interface{}) WhereRawBuilder {
	this.wheres.and().whereNotLike(column, value)
	return this
}

func (this WhereRawBuilder) WhereBetween(column string, first, second interface{}) WhereRawBuilder {
	this.wheres.and().whereBetween(column, first, second)
	return this
}

func (this WhereRawBuilder) OrWhere(column, operator string, value interface{}) WhereRawBuilder {
	this.wheres.or().where(column, operator, value)
	return this
}

func (this WhereRawBuilder) OrWhereIn(column string, value interface{}) WhereRawBuilder {
	this.wheres.or().whereIn(column, value)
	return this
}

func (this WhereRawBuilder) OrWhereNotIn(column string, value interface{}) WhereRawBuilder {
	this.wheres.or().whereNotIn(column, value)
	return this
}

func (this WhereRawBuilder) OrWhereNull(column string) WhereRawBuilder {
	this.wheres.or().whereNull(column)
	return this
}

func (this WhereRawBuilder) OrWhereNotNull(column string) WhereRawBuilder {
	this.wheres.or().whereNotNull(column)
	return this
}

func (this WhereRawBuilder) OrLike(column string, value interface{}) WhereRawBuilder {
	this.wheres.or().whereLike(column, value)
	return this
}

func (this WhereRawBuilder) OrNotLike(column string, value interface{}) WhereRawBuilder {
	this.wheres.or().whereNotLike(column, value)
	return this
}

func (this WhereRawBuilder) OrWhereBetween(column string, first, second interface{}) WhereRawBuilder {
	this.wheres.or().whereBetween(column, first, second)
	return this
}
