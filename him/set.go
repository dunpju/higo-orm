package him

type Sets []Set

func NewSets() Sets {
	return make([]Set, 0)
}

func (this Sets) Append(column any, value interface{}) {
	this = append(this, NewSet(column, value))
}

func (this Sets) ForEach(fn func(s Set) bool) {
	for _, set := range this {
		b := fn(set)
		if !b {
			break
		}
	}
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
