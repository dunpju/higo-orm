package arm

import (
	"github.com/dunpju/higo-orm/him"
)

type IModel interface {
	DB() *him.DB
	Connection() string
	TableName() *TableName
	Apply(model *Model)
	Exist() bool
}
