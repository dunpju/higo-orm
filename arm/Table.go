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
	if this.alias != "" {
		return fmt.Sprintf("%s AS %s", this.table, this.alias)
	}
	return this.table
}
