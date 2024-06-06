package him

type CaseWhen interface {
	Field() string
	Case() string
	WhenThen() []*WhenThen
	ELSE() *Else
}

type WhenThen struct {
	when, then string
}

func NewWhenThen(when, then any) *WhenThen {
	return &WhenThen{when: toString(when), then: toString(then)}
}

type Else struct {
	value string
}

func NewElse(value string) *Else {
	return &Else{value: value}
}

type Case struct {
	field, c string
	whens    []*WhenThen
	e        *Else
}

func (this *Case) Field() string {
	return this.field
}

func (this *Case) Case() string {
	return this.c
}

func (this *Case) WhenThen() []*WhenThen {
	return this.whens
}

func (this *Case) ELSE() *Else {
	return this.e
}

func NewCase(field string, c ...any) *Case {
	whens := make([]*WhenThen, 0)
	cc := ""
	if len(c) > 0 {
		cc = ToString(c[0])
	}
	return &Case{field: field, c: cc, whens: whens}
}

func (this *Case) When(when, then any) *Case {
	this.whens = append(this.whens, NewWhenThen(when, then))
	return this
}

func (this *Case) Else(value any) *Case {
	this.e = NewElse(ToString(value))
	return this
}
