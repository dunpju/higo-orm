package arm

import (
	"github.com/dunpju/higo-orm/him"
)

type IModel interface {
	DB() *him.DB
	TableName() *TableName
	Apply(model *Model)
	Exist() bool
}
