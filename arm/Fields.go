package arm

import (
	"fmt"
	"regexp"
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

type Fields string

func (this Fields) Pre(pre string) Fields {
	field := string(this)
	hasBackQuote := backQuoteReg.FindString(field)
	if hasBackQuote != "" {
		return Fields(fmt.Sprintf("%s.%s", pre, field))
	}
	return Fields(fmt.Sprintf("`%s`.`%s`", pre, field))
}

func (this Fields) AS(as string) string {
	return asHandle(string(this), as)
}

func (this Fields) ASC() string {
	field := string(this)
	hasBackQuote := backQuoteReg.FindString(field)
	if hasBackQuote != "" {
		return fmt.Sprintf("%s ASC", field)
	}
	return fmt.Sprintf("`%s` ASC", field)
}

func (this Fields) DESC() string {
	field := string(this)
	hasBackQuote := backQuoteReg.FindString(field)
	if hasBackQuote != "" {
		return fmt.Sprintf("%s DESC", field)
	}
	return fmt.Sprintf("`%s` DESC", field)
}

func (this Fields) COUNT() AsField {
	hasBackQuote := backQuoteReg.FindString(string(this))
	if hasBackQuote != "" {
		return AsField(fmt.Sprintf("COUNT(%s)", string(this)))
	}
	return AsField(fmt.Sprintf("COUNT(`%s`)", string(this)))
}

func (this Fields) SUM() AsField {
	hasBackQuote := backQuoteReg.FindString(string(this))
	if hasBackQuote != "" {
		return AsField(fmt.Sprintf("SUM(%s)", string(this)))
	}
	return AsField(fmt.Sprintf("SUM(`%s`)", string(this)))
}

func (this Fields) String() string {
	field := string(this)
	hasBackQuote := backQuoteReg.FindString(field)
	if hasBackQuote != "" {
		return field
	}
	return fmt.Sprintf("`%s`", field)
}
