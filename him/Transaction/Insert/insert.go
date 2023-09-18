package Insert

import (
	"github.com/dunpju/higo-orm/orm"
	"gorm.io/gorm"
)

func Transaction(db *gorm.DB) orm.InsertBuilder {
	return orm.InsertBuilder{DB: db}
}
