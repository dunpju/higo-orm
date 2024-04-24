package utils

import (
	"github.com/dunpju/higo-orm/arm"
	"reflect"
)

func IsEmpty(m arm.IModel) bool {
	v := reflect.ValueOf(m).Elem()
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Kind() != reflect.Ptr && !reflect.DeepEqual(v.Field(i).Interface(), reflect.Zero(v.Field(i).Type()).Interface()) {
			return true
		}
	}
	return false
}
