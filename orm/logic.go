package orm

import "fmt"

type Logic string

const (
	AND Logic = "AND"
	OR  Logic = "OR"
)

func logic(w where, sql string, arg []interface{}, pred []string, args []interface{}) ([]string, []interface{}, error) {
	if w.logic == "AND" {
		if len(pred) == 0 {
			pred = append(pred, fmt.Sprintf("(%s)", sql))
			args = append(args, arg...)
		} else {
			pred = append(pred, "AND")
			pred = append(pred, fmt.Sprintf("(%s)", sql))
			args = append(args, arg...)
		}
	} else {
		if len(pred) == 0 {
			pred = append(pred, fmt.Sprintf("(%s)", sql))
			args = append(args, arg...)
		} else {
			pred = append(pred, "OR")
			pred = append(pred, fmt.Sprintf("(%s)", sql))
			args = append(args, arg...)
		}
	}
	return pred, args, nil
}
