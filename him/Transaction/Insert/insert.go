package Insert

import (
	"github.com/dunpju/higo-orm/him"
	"gorm.io/gorm"
)

func Transaction(db *gorm.DB) him.InsertBuilder {
	return him.InsertBuilder{DB: db}
}
