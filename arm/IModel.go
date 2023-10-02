package arm

import (
	"fmt"
	"github.com/dunpju/higo-orm/him"
)

type Property func(model IModel)
type Properties []Property

func (this Properties) Apply(model IModel) {
	for _, property := range this {
		property(model)
	}
}

type IModel interface {
	DB() *him.DB
	TableName() *TableName
	Exist() bool
}

type IConstructor[T any] interface {
	New(properties ...Property) T
}

type TableName struct {
	table, alias string
}

func NewTableName(table string) *TableName {
	return &TableName{table: table}
}

// Alias Table alias
func (this *TableName) Alias(alias string) *TableName {
	this.alias = alias
	return this
}

// GetAlias Get Table alias
func (this *TableName) GetAlias() string {
	return this.alias
}

func (this *TableName) String() string {
	if this.alias != "" {
		return fmt.Sprintf("%s AS %s", this.table, this.alias)
	}
	return this.table
}

type Fields string

func (this Fields) AS(as string) string {
	return fmt.Sprintf("%s AS %s", this, as)
}

func (this Fields) ASC() string {
	return fmt.Sprintf("%s ASC", this)
}

func (this Fields) DESC() string {
	return fmt.Sprintf("%s DESC", this)
}

func (this Fields) COUNT() string {
	return fmt.Sprintf("COUNT(%s)", this)
}

func (this Fields) SUM() string {
	return fmt.Sprintf("SUM(%s)", this)
}

func (this Fields) String() string {
	return string(this)
}
