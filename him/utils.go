package him

import "fmt"

func columnToString(column any) string {
	if c, ok := column.(string); ok {
		return c
	} else if c, ok := column.(fmt.Stringer); ok {
		return c.String()
	} else {
		return fmt.Errorf("column cannot convert to string").Error()
	}
}
