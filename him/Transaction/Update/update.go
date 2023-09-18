package Update

import (
	"github.com/dunpju/higo-orm/him"
	"gorm.io/gorm"
)

func Transaction(db *gorm.DB) him.UpdateBuilder {
	return him.UpdateBuilder{DB: db}
}
