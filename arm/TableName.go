package arm

import "fmt"

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
	isMatch := backQuoteReg.Match([]byte(this.table))
	if this.alias != "" {
		if !isMatch {
			return fmt.Sprintf("`%s` AS `%s`", this.table, this.alias)
		}
		return fmt.Sprintf("%s AS %s", this.table, this.alias)
	}
	if !isMatch {
		return fmt.Sprintf("`%s`", this.table)
	}
	return this.table
}
