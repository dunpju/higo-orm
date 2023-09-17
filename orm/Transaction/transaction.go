package Transaction

import (
	"github.com/dunpju/higo-orm/orm/Insert"
	"gorm.io/gorm"
)

type Transaction struct {
	db *gorm.DB
}

func Begin(db *gorm.DB) *Transaction {
	return &Transaction{db: db}
}

func (this *Transaction) Insert(into string) Insert.InsertBuilder {
	return Insert.Transaction(this.db).Into(into)
}

func (this *Transaction) Update(into string) {
}

func (this *Transaction) Delete(into string) {
}
