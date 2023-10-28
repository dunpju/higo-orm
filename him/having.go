package him

type having struct {
	pred interface{}
	rest []interface{}
}

func (this *SelectBuilder) Having(pred interface{}, args ...interface{}) *SelectBuilder {
	this.hasHaving = true
	this.havings = append(this.havings, having{pred: pred, rest: args})
	return this
}
