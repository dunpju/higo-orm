package arm

import (
	"database/sql"
	"github.com/dunpju/higo-orm/him"
	"gorm.io/gorm"
)

type IModel interface {
	DB() *him.DB
	Connection() string
	TableName() *TableName
	Apply(model *Model)
	Exist() bool
	IModel(properties ...him.IProperty) IModel
	Begin(opts ...*sql.TxOptions) *him.TX
	TX(tx *gorm.DB) *Model
}
