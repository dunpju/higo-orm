package him

import "fmt"

func ColumnToString(column any) string {
	return columnToString(column)
}

func columnToString(column any) string {
	if c, ok := column.(string); ok {
		return c
	} else if c, ok := column.(fmt.Stringer); ok {
		return c.String()
	} else {
		return fmt.Errorf("column cannot convert to string").Error()
	}
}

func columnsToString(columns ...any) []string {
	ret := make([]string, 0)
	for _, column := range columns {
		if c, ok := column.(string); ok {
			ret = append(ret, c)
		} else if c, ok := column.(fmt.Stringer); ok {
			ret = append(ret, c.String())
		} else {
			ret = append(ret, fmt.Errorf("column cannot convert to string").Error())
			break
		}
	}
	return ret
}
