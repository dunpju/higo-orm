package him

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"regexp"
)

type CaseWhen interface {
	Field() string
	SetField(field string)
	Case() string
	WhenThen() []*WhenThen
	ELSE() *Else
	Builder() squirrel.CaseBuilder
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

type CaseBuilder struct {
	field, c string
	whens    []*WhenThen
	e        *Else
}

func (this *CaseBuilder) Field() string {
	return this.field
}

func (this *CaseBuilder) SetField(field string) {
	hasBackQuote := regexp.MustCompile("`").FindString(field)
	if hasBackQuote != "" {
		this.field = field
	} else {
		this.field = fmt.Sprintf("`%s`", field)
	}
}

func (this *CaseBuilder) Case() string {
	return this.c
}

func (this *CaseBuilder) WhenThen() []*WhenThen {
	return this.whens
}

func (this *CaseBuilder) ELSE() *Else {
	return this.e
}

func (this *CaseBuilder) Builder() squirrel.CaseBuilder {
	builder := squirrel.Case(this.Case())
	for _, w := range this.WhenThen() {
		builder = builder.When(w.when, w.then)
	}
	if this.ELSE() != nil {
		builder = builder.Else(this.ELSE().value)
	}
	return builder
}

func NewCaseBuilder(field string, c ...any) *CaseBuilder {
	return newCaseBuilder(c...).setField(field)
}

func newCaseBuilder(c ...any) *CaseBuilder {
	whens := make([]*WhenThen, 0)
	cc := ""
	if len(c) > 0 {
		cc = ToString(c[0])
	}
	return &CaseBuilder{c: cc, whens: whens}
}

func (this *CaseBuilder) When(when, then any) *CaseBuilder {
	this.whens = append(this.whens, NewWhenThen(when, then))
	return this
}

func (this *CaseBuilder) Else(value any) *CaseBuilder {
	this.e = NewElse(ToString(value))
	return this
}

func (this *CaseBuilder) setField(field string) *CaseBuilder {
	this.SetField(field)
	return this
}

func Case(field ...any) *CaseBuilder {
	return newCaseBuilder(field...)
}
