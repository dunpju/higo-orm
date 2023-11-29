package arm

import (
	"fmt"
	"regexp"
)

var (
	backQuoteReg = regexp.MustCompile("`")
)

type Fields string

func (this Fields) Pre(pre string) Fields {
	hasBackQuote := backQuoteReg.FindString(string(this))
	if hasBackQuote != "" {
		return Fields(fmt.Sprintf("%s.%s", pre, string(this)))
	}
	return Fields(fmt.Sprintf("`%s`.`%s`", pre, string(this)))
}

func (this Fields) AS(as string) string {
	hasBackQuote := backQuoteReg.FindString(string(this))
	if hasBackQuote != "" {
		return fmt.Sprintf("%s AS %s", string(this), as)
	}
	return fmt.Sprintf("`%s` AS `%s`", string(this), as)
}

func (this Fields) ASC() string {
	hasBackQuote := backQuoteReg.FindString(string(this))
	if hasBackQuote != "" {
		return fmt.Sprintf("%s ASC", string(this))
	}
	return fmt.Sprintf("`%s` ASC", string(this))
}

func (this Fields) DESC() string {
	hasBackQuote := backQuoteReg.FindString(string(this))
	if hasBackQuote != "" {
		return fmt.Sprintf("%s DESC", string(this))
	}
	return fmt.Sprintf("`%s` DESC", string(this))
}

func (this Fields) COUNT() Fields {
	hasBackQuote := backQuoteReg.FindString(string(this))
	if hasBackQuote != "" {
		return Fields(fmt.Sprintf("COUNT(%s)", string(this)))
	}
	return Fields(fmt.Sprintf("COUNT(`%s`)", string(this)))
}

func (this Fields) SUM() Fields {
	hasBackQuote := backQuoteReg.FindString(string(this))
	if hasBackQuote != "" {
		return Fields(fmt.Sprintf("SUM(%s)", string(this)))
	}
	return Fields(fmt.Sprintf("SUM(`%s`)", string(this)))
}

func (this Fields) String() string {
	hasBackQuote := backQuoteReg.FindString(string(this))
	if hasBackQuote != "" {
		return string(this)
	}
	return fmt.Sprintf("`%s`", string(this))
}
