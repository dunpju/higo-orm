package him

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

func (this WhereRawBuilder) Where(column any, operator string, value interface{}) WhereRawBuilder {
	this.wheres.and().where(columnToString(column), operator, value)
	return this
}

func (this WhereRawBuilder) WhereIn(column any, value interface{}) WhereRawBuilder {
	this.wheres.and().whereIn(columnToString(column), value)
	return this
}

func (this WhereRawBuilder) WhereNotIn(column any, value interface{}) WhereRawBuilder {
	this.wheres.and().whereNotIn(columnToString(column), value)
	return this
}

func (this WhereRawBuilder) WhereNull(column any) WhereRawBuilder {
	this.wheres.and().whereNull(columnToString(column))
	return this
}

func (this WhereRawBuilder) WhereNotNull(column any) WhereRawBuilder {
	this.wheres.and().whereNotNull(columnToString(column))
	return this
}

func (this WhereRawBuilder) WhereLike(column any, value interface{}) WhereRawBuilder {
	this.wheres.and().whereLike(columnToString(column), value)
	return this
}

func (this WhereRawBuilder) NotLike(column any, value interface{}) WhereRawBuilder {
	this.wheres.and().whereNotLike(columnToString(column), value)
	return this
}

func (this WhereRawBuilder) WhereBetween(column any, first, second interface{}) WhereRawBuilder {
	this.wheres.and().whereBetween(columnToString(column), first, second)
	return this
}

func (this WhereRawBuilder) OrWhere(column any, operator string, value interface{}) WhereRawBuilder {
	this.wheres.or().where(columnToString(column), operator, value)
	return this
}

func (this WhereRawBuilder) OrWhereIn(column any, value interface{}) WhereRawBuilder {
	this.wheres.or().whereIn(columnToString(column), value)
	return this
}

func (this WhereRawBuilder) OrWhereNotIn(column any, value interface{}) WhereRawBuilder {
	this.wheres.or().whereNotIn(columnToString(column), value)
	return this
}

func (this WhereRawBuilder) OrWhereNull(column any) WhereRawBuilder {
	this.wheres.or().whereNull(columnToString(column))
	return this
}

func (this WhereRawBuilder) OrWhereNotNull(column any) WhereRawBuilder {
	this.wheres.or().whereNotNull(columnToString(column))
	return this
}

func (this WhereRawBuilder) OrLike(column any, value interface{}) WhereRawBuilder {
	this.wheres.or().whereLike(columnToString(column), value)
	return this
}

func (this WhereRawBuilder) OrNotLike(column any, value interface{}) WhereRawBuilder {
	this.wheres.or().whereNotLike(columnToString(column), value)
	return this
}

func (this WhereRawBuilder) OrWhereBetween(column any, first, second interface{}) WhereRawBuilder {
	this.wheres.or().whereBetween(columnToString(column), first, second)
	return this
}

func (this WhereRawBuilder) Raw(pred string, args ...interface{}) WhereRawBuilder {
	this.wheres.and().raw(pred, args, nil)
	return this
}
