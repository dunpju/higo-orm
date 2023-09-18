package him

import "fmt"

type between struct {
	column        string
	first, second interface{}
}

func (this between) ToSql() (string, []interface{}, error) {
	args := make([]interface{}, 0)
	args = append(args, this.first, this.second)
	return fmt.Sprintf("%s BETWEEN ? AND ?", this.column), args, nil
}
