package arm

import (
	"database/sql"
	"encoding/json"
	"github.com/dunpju/higo-orm/him"
	"gorm.io/gorm"
)

type IModel interface {
	DB() *him.DB
	Connection() string
	TableName() *TableName
	Apply(model *Model)
	Exist() bool
	NewModel(properties ...him.IProperty) IModel
	Begin(opts ...*sql.TxOptions) *him.TX
	TX(tx *gorm.DB) *Model
}

type Models string

func MakeModels(v interface{}) Models {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return Models(b)
}
