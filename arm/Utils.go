package arm

import (
	"reflect"
)

func IsEmpty(m IModel) bool {
	v := reflect.ValueOf(m).Elem()
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Kind() != reflect.Ptr && !reflect.DeepEqual(v.Field(i).Interface(), reflect.Zero(v.Field(i).Type()).Interface()) {
			return true
		}
	}
	return false
}
