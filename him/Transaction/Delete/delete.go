package Delete

import (
	"github.com/dunpju/higo-orm/orm"
	"gorm.io/gorm"
)

func Transaction(db *gorm.DB) orm.DeleteBuilder {
	return orm.DeleteBuilder{DB: db}
}
