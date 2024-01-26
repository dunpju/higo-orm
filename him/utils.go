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

func ColumnsToString(columns ...any) []string {
	return columnsToString(columns...)
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

type ValueToStringInterface interface {
	string | int | int8 | int16 | int32 | int64 | float32 | float64
}

func valuesToString[T ValueToStringInterface](values ...T) []string {
	ret := make([]string, 0)
	for _, value := range values {
		ret = append(ret, fmt.Sprintf("%v", value))
	}
	return ret
}
