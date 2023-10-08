package arm

import "fmt"

type Fields string

func (this Fields) Pre(pre string) Fields {
	return Fields(fmt.Sprintf("%s.%s", pre, this))
}

func (this Fields) AS(as string) string {
	return fmt.Sprintf("%s AS %s", this, as)
}

func (this Fields) ASC() string {
	return fmt.Sprintf("%s ASC", this)
}

func (this Fields) DESC() string {
	return fmt.Sprintf("%s DESC", this)
}

func (this Fields) COUNT() Fields {
	return Fields(fmt.Sprintf("COUNT(%s)", this))
}

func (this Fields) SUM() Fields {
	return Fields(fmt.Sprintf("SUM(%s)", this))
}

func (this Fields) String() string {
	return string(this)
}
