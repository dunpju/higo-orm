package pagination

import "github.com/dunpju/higo-orm/arm"

type IPaginateSum interface {
	SetSum(sum interface{})
	Dest() interface{}
	Field() arm.Fields
}
