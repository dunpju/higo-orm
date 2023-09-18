package him

import "fmt"

type Logic string

const (
	AND Logic = "AND"
	OR  Logic = "OR"
)

func logic(w where, sql string, arg []interface{}, pred []string, args []interface{}) ([]string, []interface{}, error) {
	if w.logic == AND {
		if len(pred) == 0 {
			if _, ok := w.sqlizer.(raw); ok {
				pred = append(pred, fmt.Sprintf("%s", sql))
				args = append(args, arg...)
			} else {
				pred = append(pred, fmt.Sprintf("(%s)", sql))
				args = append(args, arg...)
			}
		} else {
			pred = append(pred, string(AND))
			pred = append(pred, fmt.Sprintf("(%s)", sql))
			args = append(args, arg...)
		}
	} else {
		if len(pred) == 0 {
			pred = append(pred, fmt.Sprintf("(%s)", sql))
			args = append(args, arg...)
		} else {
			pred = append(pred, string(OR))
			pred = append(pred, fmt.Sprintf("(%s)", sql))
			args = append(args, arg...)
		}
	}
	return pred, args, nil
}
