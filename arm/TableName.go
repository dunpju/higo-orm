package arm

import (
	"fmt"
	"strings"
)

type TableName struct {
	table, alias string
	forceIndex   []string
}

func NewTableName(table string) *TableName {
	return &TableName{table: table, forceIndex: make([]string, 0)}
}

// Alias Table alias
func (this *TableName) Alias(alias string) *TableName {
	this.alias = alias
	return this
}

func (this *TableName) ForceIndex(index string, more ...string) *TableName {
	this.forceIndex = append(this.forceIndex, index)
	this.forceIndex = append(this.forceIndex, more...)
	return this
}

// GetAlias Get Table alias
func (this *TableName) GetAlias() string {
	return this.alias
}

func (this *TableName) String() string {
	indexBackQuote := func() {
		for i, index := range this.forceIndex {
			if !backQuoteReg.Match([]byte(index)) {
				this.forceIndex[i] = fmt.Sprintf("`%s`", index)
			}
		}
	}
	isMatch := backQuoteReg.Match([]byte(this.table))
	if this.alias != "" {
		table := fmt.Sprintf("%s AS %s", this.table, this.alias)
		if !isMatch {
			table = fmt.Sprintf("`%s` AS `%s`", this.table, this.alias)
		}
		if len(this.forceIndex) > 0 {
			indexBackQuote()
			table = fmt.Sprintf("%s FORCE INDEX (%s)", table, strings.Join(this.forceIndex, ","))
		}
		return table
	}
	if !isMatch {
		table := fmt.Sprintf("`%s`", this.table)
		if len(this.forceIndex) > 0 {
			indexBackQuote()
			table = fmt.Sprintf("%s FORCE INDEX (%s)", table, strings.Join(this.forceIndex, ","))
		}
		return table
	}
	return this.table
}
