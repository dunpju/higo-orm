package Transaction

import (
	"github.com/dunpju/higo-orm/him"
	"github.com/dunpju/higo-orm/him/Transaction/Delete"
	"github.com/dunpju/higo-orm/him/Transaction/Insert"
	"github.com/dunpju/higo-orm/him/Transaction/Update"
	"gorm.io/gorm"
)

type Transaction struct {
	db    *gorm.DB
	Error error
}

func Begin(tx ...*gorm.DB) *Transaction {
	transaction := &Transaction{}
	if len(tx) > 0 {
		transaction.db = tx[0]
	} else {
		g, err := him.Gorm()
		if err != nil {
			transaction.Error = err
		}
		transaction.db = g.Begin()
	}
	return transaction
}

func (this *Transaction) Insert(into string) him.InsertBuilder {
	return Insert.Transaction(this.db).Insert(into)
}

func (this *Transaction) Update(table ...string) him.UpdateBuilder {
	return Update.Transaction(this.db).Update(table...)
}

func (this *Transaction) Delete(from ...string) him.DeleteBuilder {
	return Delete.Transaction(this.db).Delete(from...)
}
