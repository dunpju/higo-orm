package arm

import (
	"fmt"
	"github.com/dunpju/higo-orm/him"
	"regexp"
	"strings"
)

var (
	backQuoteReg = regexp.MustCompile("`")
)

type AsField string

func (this AsField) AS(as string) string {
	return asHandle(string(this), as)
}

func (this AsField) String() string {
	return string(this)
}

func asHandle(this string, as string) string {
	hasBackQuote := backQuoteReg.FindString(this)
	if hasBackQuote != "" {
		return fmt.Sprintf("%s AS `%s`", this, as)
	}
	return fmt.Sprintf("`%s` AS `%s`", this, as)
}

type ValueToStringInterface interface {
	string | int | int8 | int16 | int32 | int64 | float32 | float64
}

type Fields string

func (this Fields) string() string {
	return string(this)
}

func (this Fields) Case(field ...any) *Case {
	return newCase(this.string(), field...)
}

type Case struct {
	field, c string
	whens    []*him.WhenThen
	e        *him.Else
}

func (this *Case) Field() string {
	return this.field
}

func (this *Case) Case() string {
	return this.c
}

func (this *Case) WhenThen() []*him.WhenThen {
	return this.whens
}

func (this *Case) ELSE() *him.Else {
	return this.e
}

func newCase(field string, c ...any) *Case {
	whens := make([]*him.WhenThen, 0)
	cc := ""
	if len(c) > 0 {
		cc = him.ToString(c[0])
	}
	return &Case{field: field, c: cc, whens: whens}
}

func (this *Case) When(when, then any) *Case {
	this.whens = append(this.whens, him.NewWhenThen(when, then))
	return this
}

func (this *Case) Else(value any) *Case {
	this.e = him.NewElse(him.ToString(value))
	return this
}

func (this Fields) Eq(value interface{}) string {
	field := this.string()
	hasBackQuote := backQuoteReg.FindString(field)
	if hasBackQuote != "" {
		return fmt.Sprintf("%s = %s", field, fmt.Sprintf("'%v'", value))
	}
	return fmt.Sprintf("`%s` = %s", field, fmt.Sprintf("'%v'", value))
}

func (this Fields) VALUES() string {
	field := this.string()
	hasBackQuote := backQuoteReg.FindString(field)
	if hasBackQuote != "" {
		return fmt.Sprintf("%s = VALUES(%s)", field, field)
	}
	return fmt.Sprintf("`%s` = VALUES(%s)", field, field)
}

func (this Fields) FIELD(value string, moreValue ...interface{}) string {
	values := []string{fmt.Sprintf("'%s'", value)}
	for _, value := range moreValue {
		values = append(values, fmt.Sprintf("'%v'", value))
	}
	field := this.string()
	hasBackQuote := backQuoteReg.FindString(field)
	if hasBackQuote != "" {
		return fmt.Sprintf("FIELD(%s, %s)", field, strings.Join(values, ","))
	}
	return fmt.Sprintf("FIELD(`%s`, %s)", field, strings.Join(values, ","))
}

func (this Fields) IN(value interface{}, moreValue ...interface{}) string {
	values := []string{fmt.Sprintf("'%v'", value)}
	for _, value := range moreValue {
		values = append(values, fmt.Sprintf("'%v'", value))
	}
	field := this.string()
	hasBackQuote := backQuoteReg.FindString(field)
	if hasBackQuote != "" {
		return fmt.Sprintf("%s IN (%s)", field, strings.Join(values, ","))
	}
	return fmt.Sprintf("`%s` IN(%s)", field, strings.Join(values, ","))
}

func (this Fields) NotIn(value interface{}, moreValue ...interface{}) string {
	values := []string{fmt.Sprintf("'%v'", value)}
	for _, value := range moreValue {
		values = append(values, fmt.Sprintf("'%v'", value))
	}
	field := this.string()
	hasBackQuote := backQuoteReg.FindString(field)
	if hasBackQuote != "" {
		return fmt.Sprintf("%s NOT IN (%s)", field, strings.Join(values, ","))
	}
	return fmt.Sprintf("`%s` NOT IN(%s)", field, strings.Join(values, ","))
}

func (this Fields) Pre(pre string) Fields {
	hasBackQuote := backQuoteReg.FindString(pre)
	if hasBackQuote == "" {
		pre = fmt.Sprintf("`%s`", pre)
	}
	field := this.string()
	hasBackQuote = backQuoteReg.FindString(field)
	if hasBackQuote == "" {
		field = fmt.Sprintf("`%s`", field)
	}
	return Fields(fmt.Sprintf("%s.%s", pre, field))
}

func (this Fields) AS(as string) string {
	return asHandle(this.string(), as)
}

func (this Fields) ASC() string {
	field := this.string()
	hasBackQuote := backQuoteReg.FindString(field)
	if hasBackQuote != "" {
		return fmt.Sprintf("%s ASC", field)
	}
	return fmt.Sprintf("`%s` ASC", field)
}

func (this Fields) DESC() string {
	field := this.string()
	hasBackQuote := backQuoteReg.FindString(field)
	if hasBackQuote != "" {
		return fmt.Sprintf("%s DESC", field)
	}
	return fmt.Sprintf("`%s` DESC", field)
}

func (this Fields) COUNT() AsField {
	field := this.string()
	hasBackQuote := backQuoteReg.FindString(field)
	if hasBackQuote != "" {
		return AsField(fmt.Sprintf("COUNT(%s)", field))
	}
	return AsField(fmt.Sprintf("COUNT(`%s`)", field))
}

func (this Fields) SUM() AsField {
	field := this.string()
	hasBackQuote := backQuoteReg.FindString(field)
	if hasBackQuote != "" {
		return AsField(fmt.Sprintf("SUM(%s)", field))
	}
	return AsField(fmt.Sprintf("SUM(`%s`)", field))
}

func (this Fields) String() string {
	field := this.string()
	hasBackQuote := backQuoteReg.FindString(field)
	if hasBackQuote == "" {
		field = fmt.Sprintf("`%s`", field)
	}
	return field
}
