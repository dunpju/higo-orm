package Update

import (
	"github.com/dunpju/higo-orm/orm"
	"gorm.io/gorm"
)

func Transaction(db *gorm.DB) orm.UpdateBuilder {
	return orm.UpdateBuilder{DB: db}
}
