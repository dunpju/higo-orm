package arm

import (
	"github.com/dunpju/higo-orm/him"
)

type IModel interface {
	IModel(properties ...him.IProperty) IModel
	DB() *him.DB
	Connection() string
	TableName() *TableName
	Apply(model *Model)
	Exist() bool
}
