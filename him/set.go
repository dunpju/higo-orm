package him

type Sets []Set

func NewSets() Sets {
	return make([]Set, 0)
}

func (this Sets) Reset() {
	this = NewSets()
}

func (this Sets) Append(column any, value interface{}) {
	this = append(this, NewSet(column, value))
}

func (this Sets) ForEach(fn func(s Set) bool) Sets {
	for _, set := range this {
		b := fn(set)
		if !b {
			break
		}
	}
	return this
}

type Set struct {
	column any
	value  interface{}
}

func NewSet(column any, value interface{}) Set {
	return Set{column: column, value: value}
}

func (s Set) Column() any {
	return s.column
}

func (s Set) Value() interface{} {
	return s.value
}
