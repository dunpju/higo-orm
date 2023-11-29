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
	isMatch := backQuoteReg.Match([]byte(this.String()))
	if isMatch {
		return Fields(fmt.Sprintf("%s.%s", pre, this))
	}
	return Fields(fmt.Sprintf("`%s`.`%s`", pre, this))
}

func (this Fields) AS(as string) string {
	isMatch := backQuoteReg.Match([]byte(this.String()))
	if isMatch {
		return fmt.Sprintf("%s AS %s", this, as)
	}
	return fmt.Sprintf("`%s` AS `%s`", this, as)
}

func (this Fields) ASC() string {
	isMatch := backQuoteReg.Match([]byte(this.String()))
	if isMatch {
		return fmt.Sprintf("%s ASC", this)
	}
	return fmt.Sprintf("`%s` ASC", this)
}

func (this Fields) DESC() string {
	isMatch := backQuoteReg.Match([]byte(this.String()))
	if isMatch {
		return fmt.Sprintf("%s DESC", this)
	}
	return fmt.Sprintf("`%s` DESC", this)
}

func (this Fields) COUNT() Fields {
	isMatch := backQuoteReg.Match([]byte(this.String()))
	if isMatch {
		return Fields(fmt.Sprintf("COUNT(%s)", this))
	}
	return Fields(fmt.Sprintf("COUNT(`%s`)", this))
}

func (this Fields) SUM() Fields {
	isMatch := backQuoteReg.Match([]byte(this.String()))
	if isMatch {
		return Fields(fmt.Sprintf("SUM(%s)", this))
	}
	return Fields(fmt.Sprintf("SUM(`%s`)", this))
}

func (this Fields) String() string {
	isMatch := backQuoteReg.Match([]byte(this))
	if isMatch {
		return string(this)
	}
	return fmt.Sprintf("`%s`", string(this))
}
