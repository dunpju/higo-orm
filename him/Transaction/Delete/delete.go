package Delete

import (
	"github.com/dunpju/higo-orm/him"
	"gorm.io/gorm"
)

func Transaction(db *gorm.DB) him.DeleteBuilder {
	return him.DeleteBuilder{DB: db}
}
